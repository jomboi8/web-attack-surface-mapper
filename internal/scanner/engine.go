package scanner

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"apguard/internal/scanner/detectors"
	"apguard/pkg/logger"
)

// Payload sets embedded directly for portability
var sqliPayloads = []string{
	"'", "\"", "' OR '1'='1", "' OR 1=1--", "\" OR \"1\"=\"1",
	"' OR 'x'='x", "1 OR 1=1", "' AND 1=0 UNION SELECT NULL--",
	"'; DROP TABLE users--", "' UNION SELECT 1,2,3--",
}

var xssPayloads = []string{
	`<script>alert('XSS')</script>`,
	`<img src=x onerror=alert('XSS')>`,
	`<svg onload=alert('XSS')>`,
	`javascript:alert('XSS')`,
	`"><script>alert(document.domain)</script>`,
	`'"><img src=# onerror=alert(1)>`,
}

var lfiPayloads = []string{
	"../../../etc/passwd",
	"..\\..\\..\\windows\\win.ini",
	"/etc/passwd",
	"....//....//....//etc/passwd",
	"php://filter/convert.base64-encode/resource=index.php",
}

var commonParams = []string{"id", "q", "search", "page", "file", "path", "url", "input", "data", "query", "name", "user"}

// ScanResult is a single vulnerability found during a scan.
type ScanResult struct {
	ScanID      int
	Type        string
	Severity    string
	Endpoint    string
	Parameter   string
	Payload     string
	Description string
}

// Engine orchestrates the scanning process.
type Engine struct {
	db         *sql.DB
	httpClient *HTTPClient
	mu         sync.Mutex
}

func NewEngine(db *sql.DB) *Engine {
	return &Engine{
		db:         db,
		httpClient: NewHTTPClient(15 * time.Second),
	}
}

// RunScan runs a full scan on the given target URL and scan record ID.
// It runs concurrently using goroutines per parameter and reports findings.
func (e *Engine) RunScan(scanID int, targetURL, profile string) ([]ScanResult, error) {
	logger.Info("Engine: starting scan #%d on %s (profile=%s)", scanID, targetURL, profile)

	// Update scan status to running
	e.updateScanStatus(scanID, "running", "")

	var (
		results []ScanResult
		wg      sync.WaitGroup
		resChan = make(chan ScanResult, 100)
	)

	params := commonParams
	if profile == "quick" {
		params = commonParams[:4]
	}

	// Launch concurrent goroutines per parameter
	for _, param := range params {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			found := e.scanParameter(scanID, targetURL, p, profile)
			for _, f := range found {
				resChan <- f
			}
		}(param)
	}

	// Also run misconfiguration check
	wg.Add(1)
	go func() {
		defer wg.Done()
		probe := e.httpClient.Inject(targetURL, "", "")
		if probe.Error == nil {
			headersMap := make(map[string][]string)
			for k, v := range probe.Headers {
				headersMap[k] = v
			}
			if f := detectors.DetectMisconfig(probe.StatusCode, probe.Body, headersMap); f != nil {
				resChan <- ScanResult{
					ScanID:      scanID,
					Type:        f.Type,
					Severity:    f.Severity,
					Endpoint:    targetURL,
					Parameter:   "-",
					Payload:     "-",
					Description: f.Description + " | Evidence: " + f.Evidence,
				}
			}
		}
	}()

	// Close channel when all goroutines finish
	go func() {
		wg.Wait()
		close(resChan)
	}()

	// Collect results
	for r := range resChan {
		results = append(results, r)
		e.storeVuln(r)
	}

	// Update scan record
	total := len(results)
	e.updateScanStatus(scanID, "completed", fmt.Sprintf("%d", total))
	logger.Info("Engine: scan #%d completed — %d vulnerabilities found", scanID, total)
	return results, nil
}

func (e *Engine) scanParameter(scanID int, targetURL, param, profile string) []ScanResult {
	var results []ScanResult

	// SQLi payloads
	for _, payload := range sqliPayloads {
		probe := e.httpClient.Inject(targetURL, param, payload)
		if probe.Error != nil {
			continue
		}
		if f := detectors.DetectSQLi(probe.Body); f != nil {
			results = append(results, ScanResult{
				ScanID:      scanID,
				Type:        f.Type,
				Severity:    f.Severity,
				Endpoint:    targetURL,
				Parameter:   param,
				Payload:     payload,
				Description: f.Description + " | Evidence: " + f.Evidence,
			})
			break // one finding per param is enough for this type
		}
	}

	// XSS payloads
	for _, payload := range xssPayloads {
		probe := e.httpClient.Inject(targetURL, param, payload)
		if probe.Error != nil {
			continue
		}
		if f := detectors.DetectXSS(probe.Body, payload); f != nil {
			results = append(results, ScanResult{
				ScanID:      scanID,
				Type:        f.Type,
				Severity:    f.Severity,
				Endpoint:    targetURL,
				Parameter:   param,
				Payload:     payload,
				Description: f.Description + " | Evidence: " + f.Evidence,
			})
			break
		}
	}

	// LFI payloads (only full profile)
	if profile != "quick" {
		for _, payload := range lfiPayloads {
			probe := e.httpClient.Inject(targetURL, param, payload)
			if probe.Error != nil {
				continue
			}
			if f := detectors.DetectLFI(probe.Body); f != nil {
				results = append(results, ScanResult{
					ScanID:      scanID,
					Type:        f.Type,
					Severity:    f.Severity,
					Endpoint:    targetURL,
					Parameter:   param,
					Payload:     payload,
					Description: f.Description + " | Evidence: " + f.Evidence,
				})
				break
			}
		}
	}

	return results
}

func (e *Engine) storeVuln(r ScanResult) {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, err := e.db.Exec(
		`INSERT INTO vulnerabilities (scan_id, type, severity, endpoint, parameter, payload, description)
		 VALUES (?,?,?,?,?,?,?)`,
		r.ScanID, r.Type, r.Severity, r.Endpoint, r.Parameter, r.Payload, r.Description,
	)
	if err != nil {
		logger.Error("failed to store vulnerability: %v", err)
	}
}

func (e *Engine) updateScanStatus(scanID int, status, total string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if strings.EqualFold(status, "running") {
		e.db.Exec(`UPDATE scans SET status=?, start_time=CURRENT_TIMESTAMP WHERE id=?`, status, scanID)
	} else {
		totalInt := 0
		fmt.Sscanf(total, "%d", &totalInt)
		e.db.Exec(
			`UPDATE scans SET status=?, end_time=CURRENT_TIMESTAMP, total_vulnerabilities=? WHERE id=?`,
			status, totalInt, scanID,
		)
	}
}
