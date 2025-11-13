package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	port   int
	server *http.Server
}

func NewHTTPServer(port int) *HTTPServer {
	return &HTTPServer{
		port: port,
	}
}

// Start begins listening for HTTP requests and echoes received data
func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleEcho)
	mux.HandleFunc("/health", s.handleHealth)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("HTTP echo server listening on port %d", s.port)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

// handleEcho handles echo requests
func (s *HTTPServer) handleEcho(w http.ResponseWriter, r *http.Request) {
	remoteAddr := r.RemoteAddr
	log.Printf("HTTP request from %s: %s %s", remoteAddr, r.Method, r.URL.Path)

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("HTTP read error from %s: %v", remoteAddr, err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Echo back the body if present, otherwise send a default message
	if len(body) > 0 {
		log.Printf("HTTP received %d bytes from %s", len(body), remoteAddr)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Printf("HTTP write error to %s: %v", remoteAddr, err)
			return
		}
		log.Printf("HTTP echoed %d bytes to %s", len(body), remoteAddr)
	} else {
		// For GET requests with no body, return a simple message
		message := "OK"
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(message))
		if err != nil {
			log.Printf("HTTP write error to %s: %v", remoteAddr, err)
			return
		}
		log.Printf("HTTP sent OK response to %s", remoteAddr)
	}
}

// handleHealth handles health check requests
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}

// Stop stops the HTTP server gracefully
func (s *HTTPServer) Stop() error {
	if s.server != nil {
		log.Printf("Shutting down HTTP server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}
