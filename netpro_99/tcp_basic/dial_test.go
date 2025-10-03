package tcpbasic

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {

	// net.Listen 을 사용하여 루프백 아이피의 랜덤 포트를 통해 TCP 연결 포트 바인딩을 수행합니다.
	l, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// 종료 처리를 위한 채널을 생성합니다.
	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()

		for {
			// 리스너를 통해 들어온 tcp 연결을 수락하여 새로운 tcp 세션을 만듭니다.
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

				// 한번에 1024 바이트만큼 데이터를 수신하며, EOF 가 발생하기 전까지는 계속해서 데이터를 수신합니다.
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

	// Go routine 으로 실행되는 listener 의 동작과 별개로 메인 고루틴에서 Dial 요청을 보냅니다.
	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// 연결된 tcp connection 을 통해 hello 라는 문자열을 발신합니다.
	_, err = conn.Write([]byte("hello"))
	if err != nil {
		t.Error(err)
	}

	// 안정적인 연결 종료
	conn.Close() // connection close 로 tcp session 에 EOF 가 발생합니다.
	<-done       // 이때 고루틴으로 실행되는 listener 의 핸들러가 종료되며 defer 에서 done 채널로 struct 를 발송할 때 까지 블락됩니다.
	l.Close()    // 리스너를 close 하게 되면 line28 에 error 가 발생합니다.
	<-done       // 고루틴으로 실행되는 함수의 defer 에서 done 채널로 struct 발송할 때 까지 블락됩니다.
}
