# protem-gen

Go 웹 애플리케이션 프로젝트 생성기 CLI 도구

## Overview

protem-gen은 Go + Tailwind CSS + templ + htmx 기반의 웹 애플리케이션 프로젝트를 빠르게 생성하는 CLI 도구입니다. Clean Architecture 패턴을 따르며, 핫 리로드 개발 환경을 자동으로 구성합니다.

### Features

- **빠른 프로젝트 설정**: 2-4시간 걸리던 설정을 5분 이내로 단축
- **Clean Architecture**: 도메인 기반 디렉토리 구조 자동 생성
- **핫 리로드**: air + templ + tailwind 통합 개발 환경
- **타입 안전 SQL**: sqlc를 통한 타입 안전 데이터베이스 쿼리
- **선택적 기능**: gRPC, Auth, AI Ready 옵션

## Installation

```bash
go install github.com/hippo-an/tiny-go-challenges/protem-gen@latest
```

### Prerequisites

다음 도구들이 설치되어 있어야 합니다:

| Tool | Installation |
|------|--------------|
| Go 1.22+ | https://go.dev/dl/ |
| Node.js/npm | https://nodejs.org/ |
| air | `go install github.com/air-verse/air@latest` |
| templ | `go install github.com/a-h/templ/cmd/templ@latest` |

## Quick Start

```bash
# 1. 대화형 프로젝트 생성
protem-gen create

# 2. 디렉토리 이동
cd my-app

# 3. 의존성 설치
make setup

# 4. 개발 서버 시작
make dev

# 5. 브라우저 열기
# http://localhost:8080
```

## CLI Usage

### 대화형 모드

```bash
protem-gen create
```

대화형 TUI를 통해 프로젝트 설정을 진행합니다.

### 비대화형 모드

```bash
protem-gen create \
  --name my-app \
  --module github.com/user/my-app \
  --database postgres \
  --grpc \
  --auth \
  --ai \
  --no-interactive
```

### 옵션

| Flag | Description | Default |
|------|-------------|---------|
| `--name` | 프로젝트 이름 | (필수) |
| `--module` | Go 모듈 경로 | github.com/user/{name} |
| `--database` | 데이터베이스 (postgres/sqlite/none) | postgres |
| `--grpc` | gRPC 지원 포함 | false |
| `--auth` | 인증 기능 포함 | false |
| `--ai` | AI Ready 구조 포함 | false |
| `--no-interactive` | 비대화형 모드 | false |

### 버전 확인

```bash
protem-gen version
```

## Generated Project Structure

```
<project-name>/
├── cmd/
│   └── server/
│       └── main.go              # 애플리케이션 진입점
├── internal/
│   ├── domain/                  # 비즈니스 엔티티
│   ├── application/             # 유스케이스/서비스
│   ├── infrastructure/          # 외부 시스템
│   │   ├── database/
│   │   ├── http/
│   │   ├── auth/               # (--auth 옵션)
│   │   ├── llm/                # (--ai 옵션)
│   │   └── ...
│   └── interfaces/              # 핸들러/컨트롤러
│       ├── http/
│       └── grpc/               # (--grpc 옵션)
├── pkg/                         # 재사용 패키지
├── web/
│   ├── templates/               # templ 파일
│   │   ├── layouts/
│   │   ├── pages/
│   │   └── components/
│   ├── static/
│   │   └── css/
│   └── tailwind/
│       └── input.css            # Tailwind v4 CSS-first 설정
├── migrations/                  # DB 마이그레이션
├── sqlc/
│   ├── sqlc.yaml
│   └── queries/
├── proto/                       # (--grpc 옵션)
├── .air.toml                    # air 설정
├── .gitignore
├── go.mod
├── Makefile
├── package.json                 # Tailwind 의존성
└── README.md
```

## Configuration Options

### Database

| Option | Description | Driver |
|--------|-------------|--------|
| `postgres` | PostgreSQL (권장) | pgx/v5 |
| `sqlite` | SQLite | modernc.org/sqlite |
| `none` | 데이터베이스 없음 | - |

### Optional Features

| Feature | Description | Adds |
|---------|-------------|------|
| `--grpc` | gRPC 서버 지원 | proto/, grpc handlers |
| `--auth` | JWT/세션 인증 | auth middleware, JWT manager |
| `--ai` | LLM 통합 준비 | llm client, prompt manager |

## Development

### Makefile Commands

생성된 프로젝트에서 사용 가능한 명령어:

| Command | Description |
|---------|-------------|
| `make dev` | 모든 워처 병렬 실행 (air + templ + tailwind) |
| `make build` | 프로덕션 빌드 |
| `make test` | 테스트 실행 |
| `make sqlc-generate` | sqlc 코드 생성 |
| `make templ-generate` | templ 코드 생성 |
| `make setup` | 의존성 설치 |
| `make clean` | 빌드 결과물 삭제 |

### Running Tests

```bash
# 단위 테스트
go test ./...

# 커버리지
go test -cover ./...

# E2E 테스트
./scripts/e2e_test.sh
```

## Technology Stack

### Generator (protem-gen)

| Component | Technology |
|-----------|------------|
| Language | Go 1.22+ |
| CLI Framework | Cobra |
| TUI | Bubbletea |
| Template Engine | text/template |

### Generated Project

| Component | Technology |
|-----------|------------|
| Language | Go 1.23+ |
| HTTP Framework | Gin |
| Templating | templ |
| CSS | Tailwind CSS v4 |
| Interactivity | htmx |
| Client State | Alpine.js |
| SQL Codegen | sqlc |
| Hot Reload | air |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`go test ./...`)
4. Commit your changes
5. Push to the branch
6. Open a Pull Request

## License

MIT License
