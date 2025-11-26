# Smoker

K8s 네트워크 진단용 echo 서버 + 클라이언트.

DaemonSet으로 각 노드에서 서버로 TCP/UDP/HTTP ping을 날려서 네트워크 상태 확인.

## 로컬 테스트

```bash
# 서버
go run cmd/server/main.go

# 클라이언트 (다른 터미널)
SERVER_HOST=localhost go run cmd/client/main.go
```

## K8s 배포

```bash
# 이미지 빌드
make docker-build

# 배포
make deploy

# 로그
kubectl logs -l app=smoker-client -f
```

## 서버 포트

- TCP: 8080 (echo)
- UDP: 8081 (echo)
- HTTP: 8082 (`/health` → healthy, `/` → OK)

## 환경변수

서버: `TCP_PORT`, `UDP_PORT`, `HTTP_PORT`

클라이언트: `SERVER_HOST`, `TEST_INTERVAL`(기본 30초), `TEST_TIMEOUT`(기본 5초)
