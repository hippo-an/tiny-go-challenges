package tftp

import (
	"bytes"
	"errors"
	"log"
	"net"
	"time"
)

type Server struct {
	Payload []byte
	Retries uint8
	Timeout time.Duration
}

func (s Server) ListenAndServe(addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	log.Printf("Listening on %s...\n", conn.LocalAddr())

	return s.Serve(conn)
}

func (s *Server) Serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("nil connection")
	}

	if s.Payload == nil {
		return errors.New("payload is required")
	}

	if s.Retries == 0 {
		s.Retries = 10
	}

	if s.Timeout == 0 {
		s.Timeout = 6 * time.Second
	}

	var rrq ReadReq
	for {
		buf := make([]byte, DataGramSize)
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}

		err = rrq.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("[%s] bad request: %v", addr, err)
			continue
		}

		go s.handle(addr.String(), rrq)
	}
}

// Server 의 필드값에 접근 필요 -> not function but method
func (s Server) handle(clientAddr string, rrq ReadReq) {
	log.Printf("[%s] requested file: %s", clientAddr, rrq.FileName)
	conn, err := net.Dial("udp", clientAddr)
	if err != nil {
		log.Printf("[%s] dial: %v", clientAddr, err)
		return
	}

	defer conn.Close()

	var (
		ackPkt  Ack
		errPkt  Err
		dataPkt = Data{Payload: bytes.NewReader(s.Payload)}
		buf     = make([]byte, DataGramSize)
	)

NEXTPACKET:
	for n := DataGramSize; n == DataGramSize; {
		data, err := dataPkt.MarshalBinary()
		if err != nil {
			log.Printf("[%s] preparing data packet: %v", clientAddr, err)
			return
		}
	RETRY:
		for i := s.Retries; i > 0; i-- {
			// 루프 순회시 패킷의 크기를 확인하는데 사용하는 변수인 n 을 계속해서 업데이트 한다.
			n, err = conn.Write(data)
			if err != nil {
				log.Printf("[%s] write: %v", clientAddr, err)
				return
			}

			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))

			// client 로 부터 수신 확인 패킷을 읽어온다.
			_, err = conn.Read(buf)
			if err != nil {
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					continue RETRY
				}
				log.Printf("[%s] waiting for ACK: %v", clientAddr, err)
				return
			}

			// 클라이언트로부터 읽은 바이트를 언마샬링하여 응답 패킷 코드를 확인한다.
			switch {
			case ackPkt.UnmarshalBinary(buf) == nil:
				if uint16(ackPkt) == dataPkt.Block {
					continue NEXTPACKET
				}
			case errPkt.UnmarshalBinary(buf) == nil:
				log.Printf("[%s] received error: %v", clientAddr, errPkt.Message)
				return
			default:
				log.Printf("[%s] bad packet", clientAddr)

			}

		}
		log.Printf("[%s] exgausted retries", clientAddr)
		return
	}
	log.Printf("[%s] went %d blocks", clientAddr, dataPkt.Block)

}
