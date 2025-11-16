# Smoker: Kubernetes 네트워크 진단 도구

## 프로젝트 개요

Echo 서버와 진단 클라이언트로 구성된 Kubernetes 클러스터 네트워크 모니터링 도구. TCP/UDP/HTTP 프로토콜의 연결성과 latency를 주기적으로 테스트하여 네트워크 문제를 조기 발견.

## 현재 상태

### ✅ 완료 (Phase 1 & 2)
**서버 컴포넌트**
- TCP/UDP/HTTP Echo 서버 구현 (`internal/server/`)
- 서버 메인 진입점 (`cmd/server/main.go`)
- Graceful shutdown, 환경변수 설정 지원

**클라이언트 컴포넌트**
- TCP/UDP/HTTP Ping 클라이언트 구현 (`internal/client/`)
- 클라이언트 메인 진입점 (`cmd/client/main.go`)
- 주기적 테스트 실행, 구조화된 로그 출력
- 로컬 통합 테스트 성공

**문서**
- CLAUDE.md (개발 가이드)

### ⬜ 다음 작업 (Phase 3, 4, 5)
1. **Dockerfile** - Multi-stage build로 서버/클라이언트 이미지 생성
2. **Kubernetes 매니페스트**
   - `deployments/server.yaml` - Deployment + Service
   - `deployments/client.yaml` - DaemonSet
3. **README.md** - 사용자 문서 (빌드/배포/사용법)
4. **Makefile** - 빌드/배포 자동화

## 아키텍처

**서버 (Deployment)**
- TCP(8080)/UDP(8081)/HTTP(8082) echo 서버
- ClusterIP Service로 노출: `smoker-server`

**클라이언트 (DaemonSet)**
- 모든 노드에 1개씩 배포
- 30초마다 서버에 TCP/UDP/HTTP 테스트 실행
- 결과를 stdout으로 로깅

## 핵심 기능

**서버**
- TCP(8080): "PING" 메시지를 그대로 반환
- UDP(8081): "PING" 패킷을 그대로 반환
- HTTP(8082): GET `/` → "OK", `/health` → "healthy"

**클라이언트**
- 30초마다 TCP/UDP/HTTP 순차 테스트
- Latency 측정 (ms)
- 로그 형식: `[TIMESTAMP] [NODE:name] [PROTOCOL] STATUS latency=Xms`

## 기술 스택

- **언어:** Go 1.25 (표준 라이브러리만 사용, 외부 의존성 없음)
- **컨테이너:** Docker (Multi-stage build)
- **배포:** Kubernetes (Deployment + DaemonSet)

## 다음 작업 (Phase 3, 4, 5)

### 1. Dockerfile 작성
Multi-stage build로 경량 이미지 생성:
- Stage 1: Go builder
- Stage 2: Server 이미지 (Alpine 베이스)
- Stage 3: Client 이미지 (Alpine 베이스)

### 2. Kubernetes 매니페스트
**`deployments/server.yaml`**
- Deployment (replicas: 1)
- Service (ClusterIP, 3개 포트 노출)
- 환경변수: TCP_PORT, UDP_PORT, HTTP_PORT

**`deployments/client.yaml`**
- DaemonSet (모든 노드에 배포)
- 환경변수: SERVER_HOST, NODE_NAME (Downward API)

### 3. 문서 및 자동화
**`README.md`**
- 프로젝트 소개, 빌드/배포/사용법

**`Makefile`**
- `build`: 로컬 빌드
- `docker-build`: 이미지 빌드
- `deploy`: K8s 배포
- `logs-server/client`: 로그 확인
- `clean`: 리소스 정리

## 환경 변수

**서버**
- `TCP_PORT` (기본값: 8080)
- `UDP_PORT` (기본값: 8081)
- `HTTP_PORT` (기본값: 8082)

**클라이언트**
- `SERVER_HOST` (기본값: smoker-server)
- `TCP_PORT`, `UDP_PORT`, `HTTP_PORT` (기본값: 8080, 8081, 8082)
- `TEST_INTERVAL` (기본값: 30초)
- `TEST_TIMEOUT` (기본값: 5초)
- `NODE_NAME` (K8s Downward API로 주입)

## 테스트 방법

**로컬**
```bash
# 터미널 1: 서버 실행
go run cmd/server/main.go

# 터미널 2: 클라이언트 실행
SERVER_HOST=localhost go run cmd/client/main.go
```

**Kubernetes**
```bash
# 배포
kubectl apply -f deployments/

# 로그 확인
kubectl logs -l app=smoker-server -f
kubectl logs -l app=smoker-client -f
```

## MVP 완료 기준

- ✅ 서버/클라이언트 로컬 실행 및 테스트 성공
- ⬜ Docker 이미지 빌드 성공
- ⬜ Kubernetes 배포 성공
- ⬜ DaemonSet이 모든 노드에 배포
- ⬜ kubectl logs로 진단 정보 확인 가능

## 진행률: 60% (Phase 1, 2 완료)
