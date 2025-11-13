package client

import (
	"fmt"
	"net"
	"time"
)

// TCPPing performs a TCP ping test to the specified server
// Returns latency in milliseconds and any error encountered
func TCPPing(host string, port int, timeout time.Duration) (int64, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	message := []byte("PING")

	// Start timing
	startTime := time.Now()

	// Connect to the server with timeout
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return 0, fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	// Set deadline for read/write operations
	conn.SetDeadline(time.Now().Add(timeout))

	// Send message
	_, err = conn.Write(message)
	if err != nil {
		return 0, fmt.Errorf("write failed: %w", err)
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return 0, fmt.Errorf("read failed: %w", err)
	}

	// Verify echo response
	if n != len(message) {
		return 0, fmt.Errorf("response size mismatch: expected %d, got %d", len(message), n)
	}

	if string(buffer[:n]) != string(message) {
		return 0, fmt.Errorf("response content mismatch: expected %s, got %s", message, buffer[:n])
	}

	// Calculate latency
	latency := time.Since(startTime).Milliseconds()

	return latency, nil
}
