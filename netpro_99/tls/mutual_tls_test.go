package tls

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"os"
	"strings"
	"testing"
)

func caCertPool(caCertFn string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caCertFn)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to add certificate to pool")
	}

	return certPool, nil
}

func TestMutualTLSAuthentication(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Server TLS
	serverPool, err := caCertPool("client.crt")
	if err != nil {
		t.Fatal(err)
	}

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		t.Fatal(err)
	}

	serverConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		GetConfigForClient: func(chi *tls.ClientHelloInfo) (*tls.Config, error) {
			return &tls.Config{
				Certificates:     []tls.Certificate{cert},
				ClientAuth:       tls.RequireAndVerifyClientCert,
				ClientCAs:        serverPool,
				CurvePreferences: []tls.CurveID{tls.CurveP256},
				MinVersion:       tls.VersionTLS13,
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					opts := x509.VerifyOptions{
						KeyUsages: []x509.ExtKeyUsage{
							x509.ExtKeyUsageClientAuth,
						},
						Roots: serverPool,
					}

					ip := strings.Split(chi.Conn.RemoteAddr().String(), ":")[0]
					hostNames, err := net.LookupAddr(ip)
					if err != nil {
						t.Errorf("PTR lookup: %v", err)
					}
					hostNames = append(hostNames, ip)

					for _, chain := range verifiedChains {
						opts.Intermediates = x509.NewCertPool()
						for _, cert := range chain[1:] {
							opts.Intermediates.AddCert(cert)
						}

						for _, hostName := range hostNames {
							opts.DNSName = hostName
							_, err = chain[0].Verify(opts)
							if err == nil {
								return nil
							}
						}
					}

					return errors.New("client authentication failed")
				},
			}, nil
		},
	}

	serverAddr := "localhost:44443"
	server := NewTLSServer(ctx, serverAddr, 0, serverConfig)
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

	// Client TLS
	clientPool, err := caCertPool("server.crt")
	if err != nil {
		t.Fatal(err)
	}
	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := tls.Dial("tcp", serverAddr, &tls.Config{
		Certificates:     []tls.Certificate{clientCert},
		CurvePreferences: []tls.CurveID{tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
		RootCAs:          clientPool,
	})

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

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}

	cancel()
	<-done
}
