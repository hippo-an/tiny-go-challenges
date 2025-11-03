package tcpdata

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

const payload = "The bigger the interface, the weaker the abstraction."

func TestScanner(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}

		defer conn.Close()

		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// io.Reader 인터페이스를 받는 Scanner 구조체 생성
	// default seperator : \n
	sc := bufio.NewScanner(conn)
	sc.Split(bufio.ScanWords) // Scanwords space sepereated word

	var words []string
	for sc.Scan() {
		te := sc.Text()
		t.Logf("Scanner: %s", te)
		words = append(words, te)
	}

	// sc.Scan() 이 false 를 반환하며 종료된 후
	// io.EOF 가 아닌 다른 에러 발생한 경우
	// sc.Err() 함수에서 Scan 동안 발생한 에러를 리턴
	err = sc.Err()

	if err != nil {
		t.Error(err)
	}

	expected := []string{"The", "bigger", "the", "interface,", "the", "weaker", "the", "abstraction."}

	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}

	t.Logf("Scanned words: %#v", words)
}
