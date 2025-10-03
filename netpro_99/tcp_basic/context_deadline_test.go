package tcpbasic

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestContextDeadline(t *testing.T) {
	dl := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel() // 테스트 종료 시 컨텍스트를 정리합니다.

	// 커스텀 다이얼러를 생성합니다.
	var d net.Dialer
	// Control 함수는 다이얼링 프로세스 중에 호출됩니다.
	// 여기서는 데드라인(5초)보다 약간 긴 시간 동안 대기하여 타임아웃을 강제로 발생시킵니다.
	d.Control = func(network, address string, c syscall.RawConn) error {
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	}

	// 컨텍스트를 사용하여 연결을 시도합니다.
	// 데드라인이 지나면 이 함수는 에러를 반환합니다.
	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil { // 에러가 없다면, 타임아웃이 발생하지 않은 것이므로 테스트 실패입니다.
		conn.Close()
		t.Fatal("connection did not timeout")
	}

	// 반환된 에러가 net.Error 타입인지 확인합니다. ner.Error 는 네트워크 연결에 대한 에러가 자세하게 포함되어 있습니다.
	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}

	// 컨텍스트의 에러가 DeadlineExceeded인지 확인합니다.
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
}
