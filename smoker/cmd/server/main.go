package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/hippo-an/tiny-go-challenges/smoker/internal/server"
)

func main() {
	// Get port configurations from environment variables
	tcpPort := getEnvAsInt("TCP_PORT", 8080)
	udpPort := getEnvAsInt("UDP_PORT", 8081)
	httpPort := getEnvAsInt("HTTP_PORT", 8082)

	log.Printf("Starting Smoker Echo Server...")
	log.Printf("TCP Port: %d, UDP Port: %d, HTTP Port: %d", tcpPort, udpPort, httpPort)

	// Create servers
	tcpServer := server.NewTCPServer(tcpPort)
	udpServer := server.NewUDPServer(udpPort)
	httpServer := server.NewHTTPServer(httpPort)

	// WaitGroup to track server goroutines
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// Start TCP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := tcpServer.Start(); err != nil {
			log.Printf("TCP server error: %v", err)
			errChan <- err
		}
	}()

	// Start UDP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := udpServer.Start(); err != nil {
			log.Printf("UDP server error: %v", err)
			errChan <- err
		}
	}()

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(); err != nil {
			log.Printf("HTTP server error: %v", err)
			errChan <- err
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v, initiating graceful shutdown...", sig)
	case err := <-errChan:
		log.Printf("Server error occurred: %v, initiating shutdown...", err)
	}

	// Stop all servers
	log.Println("Stopping servers...")

	if err := tcpServer.Stop(); err != nil {
		log.Printf("Error stopping TCP server: %v", err)
	}

	if err := udpServer.Stop(); err != nil {
		log.Printf("Error stopping UDP server: %v", err)
	}

	if err := httpServer.Stop(); err != nil {
		log.Printf("Error stopping HTTP server: %v", err)
	}

	log.Println("All servers stopped. Goodbye!")
}

// getEnvAsInt reads an environment variable and returns it as an integer
// If the variable is not set or invalid, it returns the default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s: %s, using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}
