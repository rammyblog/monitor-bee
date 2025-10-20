package checker

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	storage "github.com/rammyblog/monitor-bee/internal/storage/sql"
)

type CheckResult struct {
	Status         string // "success" or "failed"
	ResponseTimeMs int    // How long it took
	StatusCode     int    // HTTP status code
	ErrorMessage   string // If failed, what went wrong
}

func HTTPMonitor(mon storage.Monitor) CheckResult {
	timeout := time.Duration(mon.TimeoutSeconds) * time.Second
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Validate SSL
			},
		},
	}
	parsedUrl, err := url.Parse(mon.Url)
	if err != nil {
		return CheckResult{
			Status:       "failed",
			ErrorMessage: "Failed to parse URL: " + err.Error(),
		}
	}

	var bodyReader io.Reader
	if mon.Body.Valid {
		bodyReader = strings.NewReader(mon.Body.String)
	}

	req, err := http.NewRequest(mon.Method, parsedUrl.String(), bodyReader)
	if err != nil {
		return CheckResult{
			Status:       "failed",
			ErrorMessage: "Failed to create request: " + err.Error(),
		}
	}

	// Add headers
	if mon.Headers != nil {
		var headers map[string][]string
		err := json.Unmarshal(mon.Headers, &headers)
		if err == nil {
			for key, values := range headers {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		}
	}

	// Measure start time
	startTime := time.Now()

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return CheckResult{
			Status:       "failed",
			ErrorMessage: "Request failed: " + err.Error(),
		}
	}
	defer resp.Body.Close()

	// Calculate response time
	responseTime := time.Since(startTime).Milliseconds()

	// Check if status code matches expected (if specified)
	status := "success"
	if mon.ExpectedStatusCode.Valid && resp.StatusCode != int(mon.ExpectedStatusCode.Int32) {
		status = "failed"
	}

	return CheckResult{
		Status:         status,
		ResponseTimeMs: int(responseTime),
		StatusCode:     resp.StatusCode,
	}
}
