package tcpbasic

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func TestPingerDeadline(t *testing.T) {
	done := make(chan struct{})

	// "127.0.0.1:" 주소로 TCP 리스너를 생성합니다. 포트는 자동으로 선택됩니다.
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	begin := time.Now() // 테스트 시작 시간 기록

	// server 역할을 하는 함수를 고루틴에서 실행
	go func() {
		defer close(done) // 함수 종료 시 done 채널을 닫아 완료 신호를 보냅니다.

		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel() // Pinger 고루틴에 종료 신호를 보냅니다.
			conn.Close()
		}()

		// Pinger의 타이머를 리셋하기 위한 채널을 생성합니다.
		resetTimer := make(chan time.Duration, 1)
		resetTimer <- time.Second // 초기 ping 간격을 1초로 설정합니다.
		// Pinger를 별도의 고루틴으로 시작합니다.
		go Pinger(ctx, conn, resetTimer)

		// 연결에 대한 첫 데드라인을 5초로 설정합니다.
		// 5초 동안 아무런 Read/Write 활동이 없으면 연결은 타임아웃됩니다.
		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1024) // 데이터를 읽기 위한 버퍼
		for {
			n, err := conn.Read(buf)
			if err != nil {
				return
			}

			t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])

			// Pinger의 타이머를 리셋합니다.
			resetTimer <- 0
			// 데이터를 성공적으로 읽었으므로, 데드라인을 다시 5초 연장합니다.
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				t.Error(err)
				return
			}
		}
	}()

	// 클라이언트 역할을 하는 코드
	// 서버 리스너 주소로 TCP 연결을 시도합니다.
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024) // 서버로부터 데이터를 읽기 위한 버퍼

	// 서버로부터 4개의 "ping" 메시지를 읽습니다.
	for range 4 {
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	// 서버로 "PONG!!!" 메시지를 보냅니다.
	// 이 메시지를 받은 서버는 데드라인을 5초 연장합니다.
	_, err = conn.Write([]byte("PONG!!!"))
	if err != nil {
		t.Fatal(err)
	}

	// 서버가 데드라인 타임아웃으로 연결을 닫을 때까지 계속 읽기를 시도합니다.
	for range 4 {
		n, err := conn.Read(buf)
		if err != nil {
			// 서버가 연결을 닫았으므로 io.EOF 에러가 예상됩니다.
			if err != io.EOF {
				t.Fatal(err)
			}
			break // EOF를 받으면 루프를 빠져나옵니다.
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	<-done // 서버 고루틴이 종료될 때까지 대기합니다.
	end := time.Since(begin).Truncate(time.Second)
	t.Logf("[%s] done", end)

	if end != 9*time.Second {
		t.Fatalf("expected EOF at 9 seconds; actual %s", end)
	}
}

const defaultPingInterval = 30 * time.Second

// Pinger는 주어진 간격(interval)마다 "ping" 메시지를 io.Writer에 씁니다.
// reset 채널을 통해 간격을 동적으로 변경할 수 있습니다.
// 컨텍스트(ctx)가 취소되면 Pinger는 중지됩니다.
func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration
	select {
	case <-ctx.Done():
		return
	case interval = <-reset: // reset 채널에서 초기 간격 값을 받아옵니다.
	default:
		// 채널에 값이 없으면 non-blocking으로 기본값을 설정하기 위해 default를 사용합니다.
	}

	if interval <= 0 {
		interval = defaultPingInterval // 간격이 0 이하면 기본값(30초)을 사용합니다.
	}

	timer := time.NewTimer(interval)
	defer func() {
		// Pinger가 종료될 때 타이머를 확실히 중지시켜 고루틴 릭을 방지합니다.
		if !timer.Stop() {
			// 타이머가 이미 만료된 경우, 채널에 남아있는 값을 비워줍니다.
			<-timer.C
		}
	}()

	// for-select 루프는 여러 채널 이벤트를 동시에 기다립니다.
	for {
		select {
		case <-ctx.Done():
			// 컨텍스트가 취소되면(예: 부모 고루틴이 종료될 때) Pinger를 종료합니다.
			return
		case newInterval := <-reset:
			// reset 채널에 새로운 간격 값이 들어오면 타이머를 재설정합니다.
			// 먼저 현재 타이머를 멈춥니다.
			if !timer.Stop() {
				// 만약 타이머가 이미 만료되어 Stop()이 false를 반환하면,
				// 만료된 타이머의 채널을 비워줍니다.
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval // 0보다 큰 경우에만 간격을 업데이트합니다.
			}
		case <-timer.C:
			// 타이머가 만료되면 "ping" 메시지를 씁니다.
			if _, err := w.Write([]byte("ping")); err != nil {
				// 쓰기 작업에 실패하면(예: 연결이 끊긴 경우) Pinger를 종료합니다.
				return
			}
		}

		// 다음 이벤트를 위해 현재 간격(interval)으로 타이머를 리셋합니다.
		_ = timer.Reset(interval)
	}
}
