package tls

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestEchoServerTLS(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	serverAddr := "localhost:34443"
	maxIdle := time.Second
	server := NewTLSServer(ctx, serverAddr, maxIdle, nil)
	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()
		err := server.ListenAndServeTLS("server.crt", "server.key")
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Error(err)
			return
		}
	}()

	server.Ready()

	cert, err := os.ReadFile("server.crt")
	if err != nil {
		t.Fatal(err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		t.Fatal("failed to append certificate to pool")
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		RootCAs:          certPool,
	}

	conn, err := tls.Dial("tcp", serverAddr, tlsConfig)

	if err != nil {
		t.Fatal(err)
	}

	hello := []byte("hello")
	_, err = conn.Write(hello)
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		t.Fatal(err)
	}

	if actual := b[:n]; !bytes.Equal(hello, actual) {
		t.Fatalf("expected %q; actual %q", hello, actual)
	}

	time.Sleep(2 * maxIdle)
	_, err = conn.Read(b)
	if err != io.EOF {
		t.Fatal(err)
	}
	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	<-done
}
