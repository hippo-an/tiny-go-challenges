package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hippo-an/tiny-go-challenges/smoker/internal/client"
)

func main() {
	// Get configuration from environment variables
	serverHost := getEnv("SERVER_HOST", "smoker-server")
	tcpPort := getEnvAsInt("TCP_PORT", 8080)
	udpPort := getEnvAsInt("UDP_PORT", 8081)
	httpPort := getEnvAsInt("HTTP_PORT", 8082)
	testInterval := getEnvAsInt("TEST_INTERVAL", 30)
	testTimeout := getEnvAsInt("TEST_TIMEOUT", 5)
	nodeName := getEnv("NODE_NAME", "unknown")

	log.Printf("Starting Smoker Network Diagnostic Client...")
	log.Printf("Node: %s", nodeName)
	log.Printf("Server: %s", serverHost)
	log.Printf("Ports - TCP:%d UDP:%d HTTP:%d", tcpPort, udpPort, httpPort)
	log.Printf("Test Interval: %ds, Timeout: %ds", testInterval, testTimeout)

	// Create ticker for periodic tests
	ticker := time.NewTicker(time.Duration(testInterval) * time.Second)
	defer ticker.Stop()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run initial test immediately
	runTests(nodeName, serverHost, tcpPort, udpPort, httpPort, time.Duration(testTimeout)*time.Second)

	// Main loop
	for {
		select {
		case <-ticker.C:
			runTests(nodeName, serverHost, tcpPort, udpPort, httpPort, time.Duration(testTimeout)*time.Second)
		case sig := <-sigChan:
			log.Printf("Received signal: %v, shutting down...", sig)
			return
		}
	}
}

// runTests executes all three protocol tests and logs the results
func runTests(nodeName, host string, tcpPort, udpPort, httpPort int, timeout time.Duration) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// TCP Test
	tcpLatency, tcpErr := client.TCPPing(host, tcpPort, timeout)
	if tcpErr != nil {
		log.Printf("[%s] [NODE:%s] [TCP] FAILED error=%v", timestamp, nodeName, tcpErr)
	} else {
		log.Printf("[%s] [NODE:%s] [TCP] SUCCESS latency=%dms", timestamp, nodeName, tcpLatency)
	}

	// UDP Test
	udpLatency, udpErr := client.UDPPing(host, udpPort, timeout)
	if udpErr != nil {
		log.Printf("[%s] [NODE:%s] [UDP] FAILED error=%v", timestamp, nodeName, udpErr)
	} else {
		log.Printf("[%s] [NODE:%s] [UDP] SUCCESS latency=%dms", timestamp, nodeName, udpLatency)
	}

	// HTTP Test
	httpLatency, httpErr := client.HTTPPing(host, httpPort, timeout)
	if httpErr != nil {
		log.Printf("[%s] [NODE:%s] [HTTP] FAILED error=%v", timestamp, nodeName, httpErr)
	} else {
		log.Printf("[%s] [NODE:%s] [HTTP] SUCCESS latency=%dms", timestamp, nodeName, httpLatency)
	}
}

// getEnv reads an environment variable and returns it as a string
// If the variable is not set, it returns the default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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
