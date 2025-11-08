//go:build darwin || linux

package uxnidomainsock

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestEchoServerUnixDatagram(t *testing.T) {
	dir, err := os.MkdirTemp("", "echo_unixgram")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(err)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	serverSocketAddr := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))

	serverAddr, err := datagramEchoServer(ctx, "unixgram", serverSocketAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	// /tmp/echo_unixgram/123.sock 에
	err = os.Chmod(serverSocketAddr, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	clientSocketAddr := filepath.Join(dir, fmt.Sprintf("c%d.sock", os.Getpid()))

	client, err := net.ListenPacket("unixgram", clientSocketAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	err = os.Chmod(clientSocketAddr, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("ping")
	for range 3 {
		_, err = client.WriteTo(msg, serverAddr)
		if err != nil {
			t.Fatal(err)
		}
	}
	buf := make([]byte, 1024)
	for range 3 {
		n, addr, err := client.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}
		if addr.String() != serverAddr.String() {
			t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
		}
		if !bytes.Equal(msg, buf[:n]) {
			t.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}
}
