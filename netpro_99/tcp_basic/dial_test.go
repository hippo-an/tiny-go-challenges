package tcpbasic

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()

		for {
			conn, err := l.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			t.Log("tcp connection connected")
			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		t.Error(err)
	}

	conn.Close()
	<-done
	l.Close()
	<-done
}
