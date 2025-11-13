package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	DefaultBufferSize int = 4096
)

type TCPServer struct {
	port     int
	listener net.Listener
}

func NewTCPServer(port int) *TCPServer {
	return &TCPServer{
		port: port,
	}
}

// Start begins listening for TCP connections and echoes received data
func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}
	s.listener = listener

	log.Printf("TCP echo server listening on port %d", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCP accept error: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a single TCP connection
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	log.Printf("TCP connection established from %s", remoteAddr)

	// Set read/write deadlines to prevent hanging connections
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Echo back everything received
	buf := make([]byte, DefaultBufferSize)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("TCP read error from %s: %v", remoteAddr, err)
			}
			break
		}

		if n > 0 {
			log.Printf("TCP received %d bytes from %s", n, remoteAddr)

			// Echo back the data
			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Printf("TCP write error to %s: %v", remoteAddr, err)
				break
			}

			log.Printf("TCP echoed %d bytes to %s", n, remoteAddr)
		}
	}

	log.Printf("TCP connection closed from %s", remoteAddr)
}

// Stop stops the TCP server
func (s *TCPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
