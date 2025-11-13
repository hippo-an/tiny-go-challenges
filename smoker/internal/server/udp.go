package server

import (
	"fmt"
	"log"
	"net"
)

type UDPServer struct {
	port int
	conn *net.UDPConn
}

func NewUDPServer(port int) *UDPServer {
	return &UDPServer{
		port: port,
	}
}

// Start begins listening for UDP packets and echoes received data
func (s *UDPServer) Start() error {
	addr := net.UDPAddr{
		Port: s.port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return fmt.Errorf("failed to start UDP server: %w", err)
	}
	s.conn = conn

	log.Printf("UDP echo server listening on port %d", s.port)

	// Handle incoming packets
	buf := make([]byte, DefaultBufferSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("UDP read error: %v", err)
			continue
		}

		if n > 0 {
			log.Printf("UDP received %d bytes from %s", n, remoteAddr.String())

			// Echo back the data
			_, err = conn.WriteToUDP(buf[:n], remoteAddr)
			if err != nil {
				log.Printf("UDP write error to %s: %v", remoteAddr.String(), err)
				continue
			}

			log.Printf("UDP echoed %d bytes to %s", n, remoteAddr.String())
		}
	}
}

// Stop stops the UDP server
func (s *UDPServer) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
