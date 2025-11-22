package tls

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"
)

func NewTLSServer(ctx context.Context, address string, maxIdle time.Duration, tlsConfig *tls.Config) *Server {
	return &Server{
		ctx:       ctx,
		ready:     make(chan struct{}),
		addr:      address,
		maxIdle:   maxIdle,
		tlsConfig: tlsConfig,
	}
}

type Server struct {
	ctx       context.Context
	ready     chan struct{}
	addr      string
	maxIdle   time.Duration
	tlsConfig *tls.Config
}

func (s *Server) Ready() {
	if s.ready != nil {
		<-s.ready
	}
}

func (s *Server) ListenAndServeTLS(certFn, keyFn string) error {
	if s.addr == "" {
		s.addr = "localhost:443"
	}

	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("[TLS] Failed to bind to %s: %v", s.addr, err)
		return fmt.Errorf("binding to tcp %s: %w", s.addr, err)
	}

	if s.ctx != nil {
		go func() {
			<-s.ctx.Done()
			_ = l.Close()
		}()
	}

	return s.ServeTLS(l, certFn, keyFn)
}

func (s Server) ServeTLS(l net.Listener, certFn, keyFn string) error {
	if s.tlsConfig == nil {
		s.tlsConfig = &tls.Config{
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			MinVersion:       tls.VersionTLS12,
		}
	}

	if len(s.tlsConfig.Certificates) == 0 && s.tlsConfig.GetCertificate == nil {
		cert, err := tls.LoadX509KeyPair(certFn, keyFn)
		if err != nil {
			log.Printf("[TLS] Failed to load key pair: %v", err)
			return fmt.Errorf("loading key pair: %v", err)
		}
		s.tlsConfig.Certificates = []tls.Certificate{cert}
	}

	tlsListener := tls.NewListener(l, s.tlsConfig)
	if s.ready != nil {
		close(s.ready)
	}

	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Printf("[TLS] Accept error: %v", err)
			return fmt.Errorf("accept: %v", err)
		}

		remoteAddr := conn.RemoteAddr().String()

		go func() {
			defer func() {
				_ = conn.Close()
			}()

			for {
				if s.maxIdle > 0 {
					err := conn.SetDeadline(time.Now().Add(s.maxIdle))
					if err != nil {
						log.Printf("[TLS] Failed to set deadline for %s: %v", remoteAddr, err)
						return
					}
				}

				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					log.Printf("[TLS] Read error from %s: %v", remoteAddr, err)
					return
				}

				_, err = conn.Write(buf[:n])
				if err != nil {
					log.Printf("[TLS] Write error to %s: %v", remoteAddr, err)
					return
				}
			}
		}()
	}
}
