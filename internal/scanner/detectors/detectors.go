package detectors

import (
	"regexp"
	"strings"
)

// Finding represents a detected vulnerability hit.
type Finding struct {
	Type        string
	Severity    string
	Description string
	Evidence    string
}

// ---- SQL Injection Detector ----

var sqliErrorPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)sql syntax.*mysql`),
	regexp.MustCompile(`(?i)warning.*mysql_`),
	regexp.MustCompile(`(?i)valid mysql result`),
	regexp.MustCompile(`(?i)mysqlclient\.`),
	regexp.MustCompile(`(?i)postgresql.*error`),
	regexp.MustCompile(`(?i)warning.*pg_`),
	regexp.MustCompile(`(?i)valid postgresql result`),
	regexp.MustCompile(`(?i)quoted string not properly terminated`),
	regexp.MustCompile(`(?i)ORA-[0-9]{4,5}`),
	regexp.MustCompile(`(?i)microsoft.*odbc.*sql server`),
	regexp.MustCompile(`(?i)unclosed quotation mark after the character string`),
	regexp.MustCompile(`(?i)sqlexception`),
	regexp.MustCompile(`(?i)syntax error.*sql`),
	regexp.MustCompile(`(?i)jdbc.*exception`),
	regexp.MustCompile(`(?i)com\.mysql\.jdbc`),
}

func DetectSQLi(body string) *Finding {
	for _, p := range sqliErrorPatterns {
		if m := p.FindString(body); m != "" {
			return &Finding{
				Type:        "SQL Injection",
				Severity:    "CRITICAL",
				Description: "SQL error pattern detected in response body indicating possible SQLi vulnerability",
				Evidence:    truncate(m, 200),
			}
		}
	}
	return nil
}

// ---- XSS Detector ----

var xssReflectPatterns = []*regexp.Regexp{
	regexp.MustCompile(`<script[^>]*>.*?</script>`),
	regexp.MustCompile(`javascript:.*?alert`),
	regexp.MustCompile(`on\w+\s*=\s*["']?.*?alert`),
	regexp.MustCompile(`<img[^>]*onerror`),
	regexp.MustCompile(`<svg[^>]*onload`),
}

// DetectXSS checks if the XSS payload was reflected in the response.
func DetectXSS(body, payload string) *Finding {
	// Check direct reflection
	if payload != "" && strings.Contains(body, payload) {
		return &Finding{
			Type:        "Cross-Site Scripting (XSS)",
			Severity:    "HIGH",
			Description: "Injected payload was reflected in response without encoding",
			Evidence:    truncate(payload, 200),
		}
	}
	// Check for script tags
	for _, p := range xssReflectPatterns {
		if m := p.FindString(body); m != "" {
			return &Finding{
				Type:        "Cross-Site Scripting (XSS)",
				Severity:    "HIGH",
				Description: "XSS pattern detected in response body",
				Evidence:    truncate(m, 200),
			}
		}
	}
	return nil
}

// ---- LFI Detector ----

var lfiPatterns = []*regexp.Regexp{
	regexp.MustCompile(`root:.*:0:0:`),
	regexp.MustCompile(`\[extensions\]`),
	regexp.MustCompile(`windows\s*\[system32\]`),
	regexp.MustCompile(`\[boot loader\]`),
	regexp.MustCompile(`<\?php`),
	regexp.MustCompile(`failed to open stream: no such file`),
	regexp.MustCompile(`open_basedir restriction`),
}

func DetectLFI(body string) *Finding {
	for _, p := range lfiPatterns {
		if m := p.FindString(body); m != "" {
			return &Finding{
				Type:        "Local File Inclusion (LFI)",
				Severity:    "HIGH",
				Description: "File inclusion pattern detected in response body",
				Evidence:    truncate(m, 200),
			}
		}
	}
	return nil
}

// ---- Misconfiguration Detector ----

func DetectMisconfig(statusCode int, body string, headers map[string][]string) *Finding {
	// Directory listing
	if strings.Contains(body, "Index of /") || strings.Contains(body, "Directory listing for") {
		return &Finding{
			Type:        "Misconfiguration",
			Severity:    "MEDIUM",
			Description: "Directory listing is enabled on the server",
			Evidence:    "Directory listing detected in response body",
		}
	}
	// Server version disclosure
	if sv, ok := headers["Server"]; ok && len(sv) > 0 {
		if regexp.MustCompile(`(?i)(apache|nginx|iis)\/[\d.]+`).MatchString(sv[0]) {
			return &Finding{
				Type:        "Misconfiguration",
				Severity:    "LOW",
				Description: "Server version disclosed in response headers",
				Evidence:    "Server: " + sv[0],
			}
		}
	}
	return nil
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
