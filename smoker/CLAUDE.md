# CLAUDE.md

이 파일은 Claude Code (claude.ai/code)가 이 저장소의 코드 작업을 수행할 때 지침을 제공합니다.

## 프로젝트 개요

Smoker는 Echo 서버와 진단 클라이언트로 구성된 Kubernetes 클러스터 네트워크 진단 도구입니다. Kubernetes 클러스터 내에서 TCP, UDP, HTTP 프로토콜을 통해 네트워크 상태를 지속적으로 모니터링합니다.

**목적**: Node 간, Pod 간 통신을 실시간으로 테스트하여 네트워크 문제를 조기에 발견

## 개발 명령어

### 빌드
```bash
# 서버 바이너리 빌드
go build -o bin/server ./cmd/server

# 클라이언트 바이너리 빌드 (구현 완료 시)
go build -o bin/client ./cmd/client
```

### 로컬 실행
```bash
# 서버 시작 (포트 8080, 8081, 8082에서 3개의 echo 서버 실행)
go run cmd/server/main.go

# 기본 포트 오버라이드 (선택사항)
TCP_PORT=9080 UDP_PORT=9081 HTTP_PORT=9082 go run cmd/server/main.go

# 클라이언트 실행 (구현 완료 시)
SERVER_HOST=localhost go run cmd/client/main.go
```

### 테스트
```bash
# 모든 테스트 실행
go test ./...

# 상세 출력과 함께 테스트 실행
go test -v ./...

# 특정 패키지 테스트
go test ./internal/server
go test ./internal/client
```

## 아키텍처

### 멀티 프로토콜 Echo 서버 설계

서버는 **세 개의 동시 실행 echo 서버**를 별도 포트에서 운영합니다:
- **TCP 서버** (기본 포트 8080): 30초 타임아웃과 함께 연결 지향 echo 요청 처리
- **UDP 서버** (기본 포트 8081): 비연결형 데이터그램 echo 처리
- **HTTP 서버** (기본 포트 8082): Health check 및 echo 기능을 위한 REST 엔드포인트 제공

세 서버 모두 `cmd/server/main.go`에서 고루틴으로 시작되며, 시그널 채널(SIGTERM/SIGINT)을 통한 graceful shutdown을 지원합니다.

### 서버 구현 패턴

각 서버 타입(`internal/server/tcp.go`, `udp.go`, `http.go`)은 일관된 인터페이스를 따릅니다:
- `NewXXXServer(port int)` - 생성자
- `Start() error` - 서버를 시작하는 블로킹 호출
- `Stop() error` - Graceful shutdown

메인 진입점은 다음을 사용합니다:
1. WaitGroup으로 서버 고루틴 추적
2. 에러 채널로 시작 실패 감지
3. 시그널 채널로 graceful shutdown
4. 종료 시 모든 서버에 순차적으로 Stop() 호출

### 클라이언트 아키텍처 (구현 예정)

클라이언트는 **DaemonSet** (노드당 하나)으로 배포되며:
- 30초마다 주기적 테스트 실행 (`TEST_INTERVAL`로 설정 가능)
- 세 프로토콜을 순차적으로 테스트: TCP → UDP → HTTP
- Latency 측정 및 구조화된 포맷으로 결과 로깅
- Kubernetes Service DNS를 통해 서버에 연결: `smoker-server`

예상 로그 형식:
```
[TIMESTAMP] [NODE:node-name] [PROTOCOL] STATUS latency=XXms
```

### 배포 아키텍처

**서버**: Kubernetes Deployment (1개 레플리카) + 3개 포트를 노출하는 ClusterIP Service
**클라이언트**: 노드당 하나의 클라이언트 Pod를 보장하는 Kubernetes DaemonSet

## 주요 구현 세부사항

### TCP 서버 (internal/server/tcp.go)
- 연결당 고루틴을 사용하는 `net.Listen` 사용
- `conn.SetDeadline()`을 통한 30초 읽기/쓰기 데드라인 설정
- 버퍼 크기: 4096 바이트 (`DefaultBufferSize`에 정의)
- 연결 생명주기 로깅: 연결 수립 → 데이터 수신 → echo → 종료

### UDP 서버 (internal/server/udp.go)
- 비연결형 통신을 위한 `net.ListenUDP` 사용
- 단일 고루틴이 `ReadFromUDP`/`WriteToUDP`를 통해 모든 패킷 처리
- 연결 추적 없음 - 상태 비저장 패킷 echo
- 동일한 4096바이트 버퍼 크기 공유

### HTTP 서버 (internal/server/http.go)
- `http.ServeMux`를 사용하는 표준 `net/http` 사용
- 두 개의 엔드포인트:
  - `GET /` 또는 `POST /` → 요청 body를 echo (body가 비어있으면 "OK" 반환)
  - `GET /health` → health check를 위해 "healthy" 반환
- `http.Server`에 30초 읽기/쓰기 타임아웃 설정
- 5초 타임아웃 컨텍스트를 사용한 graceful shutdown

### 환경 변수

**서버**:
- `TCP_PORT` (기본값: 8080)
- `UDP_PORT` (기본값: 8081)
- `HTTP_PORT` (기본값: 8082)

**클라이언트** (예정):
- `SERVER_HOST` (기본값: "smoker-server")
- `TCP_PORT`, `UDP_PORT`, `HTTP_PORT` (서버 포트와 일치)
- `TEST_INTERVAL` (기본값: 30초)
- `TEST_TIMEOUT` (기본값: 5초)

## 현재 구현 상태

### 완료
- 동시 연결 처리를 지원하는 TCP echo 서버
- 패킷 기반 통신을 지원하는 UDP echo 서버
- Health check 엔드포인트가 있는 HTTP echo 서버
- Graceful shutdown을 지원하는 서버 메인 진입점
- TCP 클라이언트 ping 테스트 함수 (`internal/client/tcp.go`)

### 대기 중 (상세 계획은 PRD.md 참조)
- UDP 클라이언트 구현 (`internal/client/udp.go`)
- HTTP 클라이언트 구현 (`internal/client/http.go`)
- 클라이언트 메인 진입점 (`cmd/client/main.go`)
- Dockerfile (서버 및 클라이언트 이미지를 위한 멀티 스테이지 빌드)
- Kubernetes 매니페스트 (`deployments/server.yaml`, `deployments/client.yaml`)
- 빌드/배포 자동화를 위한 Makefile

## 중요한 패턴 및 규칙

### 에러 처리
- 모든 서버는 `%w`와 함께 `fmt.Errorf`를 사용하여 래핑된 에러 반환
- 에러는 로깅되지만 다른 서버를 중단시키지 않음 (고루틴을 통한 격리)
- 네트워크 에러(EOF, connection reset)는 우아하게 처리

### 로깅
- 구조화된 메시지를 위해 `log.Printf` 사용
- 컨텍스트 포함: 프로토콜, 원격 주소, 바이트 수, 작업
- 주요 생명주기 지점에서 로그: 시작, 요청 수신, 응답 전송, 종료

### Graceful Shutdown
- 메인 함수는 시그널 채널과 에러 채널을 사용하는 select 사용
- Stop() 메서드는 멱등성을 가져야 함 (여러 번 호출해도 안전)
- HTTP 서버는 graceful shutdown을 위해 타임아웃이 있는 컨텍스트 사용
- TCP/UDP 서버는 새 연결 수락을 중지하기 위해 리스너를 닫음

### 포트 설정
- Kubernetes 유연성을 위해 모든 포트는 환경 변수로 설정 가능
- 헬퍼 함수 `getEnvAsInt()`가 기본값과 함께 타입 안전 환경 변수 파싱 제공
- 잘못된 환경 변수 값은 경고 로그와 함께 기본값으로 대체

## 모듈 및 의존성

- 모듈 경로: `github.com/hippo-an/tiny-go-challenges/smoker`
- Go 버전: 1.25.3
- **외부 의존성 없음** - Go 표준 라이브러리만 사용
