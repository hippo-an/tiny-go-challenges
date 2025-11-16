package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPPing performs an HTTP ping test to the specified server
// Returns latency in milliseconds and any error encountered
func HTTPPing(host string, port int, timeout time.Duration) (int64, error) {
	url := fmt.Sprintf("http://%s:%d/", host, port)

	// Start timing
	startTime := time.Now()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Send GET request
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: expected 200, got %d", resp.StatusCode)
	}

	// Read and verify response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read body failed: %w", err)
	}

	// Verify response (should be "OK" for empty GET requests)
	expectedResponse := "OK"
	if string(body) != expectedResponse {
		return 0, fmt.Errorf("response content mismatch: expected %s, got %s", expectedResponse, string(body))
	}

	// Calculate latency
	latency := time.Since(startTime).Milliseconds()

	return latency, nil
}
