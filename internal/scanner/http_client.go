package scanner

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient wraps http.Client with convenience helpers.
type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // intentional for scanning
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("stopped after 5 redirects")
				}
				return nil
			},
		},
	}
}

// ProbeResult holds the result of a single HTTP probe.
type ProbeResult struct {
	URL        string
	StatusCode int
	Body       string
	Headers    http.Header
	Duration   time.Duration
	Error      error
}

// GET sends a GET request with the given query parameter injected.
func (h *HTTPClient) Inject(targetURL, param, payload string) *ProbeResult {
	start := time.Now()

	u, err := url.Parse(targetURL)
	if err != nil {
		return &ProbeResult{URL: targetURL, Error: err}
	}
	q := u.Query()
	q.Set(param, payload)
	u.RawQuery = q.Encode()
	finalURL := u.String()

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return &ProbeResult{URL: finalURL, Error: err}
	}
	req.Header.Set("User-Agent", "APGUARD-Scanner/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return &ProbeResult{URL: finalURL, Error: err, Duration: time.Since(start)}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // max 1MB
	return &ProbeResult{
		URL:        finalURL,
		StatusCode: resp.StatusCode,
		Body:       string(bodyBytes),
		Headers:    resp.Header,
		Duration:   time.Since(start),
	}
}

// POST sends a POST form request with the payload injected into a form field.
func (h *HTTPClient) InjectPost(targetURL, param, payload string) *ProbeResult {
	start := time.Now()
	form := url.Values{}
	form.Set(param, payload)

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(form.Encode()))
	if err != nil {
		return &ProbeResult{URL: targetURL, Error: err}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "APGUARD-Scanner/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return &ProbeResult{URL: targetURL, Error: err, Duration: time.Since(start)}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	return &ProbeResult{
		URL:        targetURL,
		StatusCode: resp.StatusCode,
		Body:       string(bodyBytes),
		Headers:    resp.Header,
		Duration:   time.Since(start),
	}
}
