# protem-gen: Product Requirements Document

> **Version**: 2.7.0
> **Status**: Implemented - Phase 7 Complete
> **Last Updated**: 2025-12-02

---

## 1. Product Overview

| Item | Description |
|------|-------------|
| **Product Name** | protem-gen |
| **Type** | CLI Binary Tool |
| **Purpose** | Go 웹 애플리케이션 프로젝트 생성기 |
| **Target Users** | Go 개발자 |

**핵심 가치**: 반복적인 프로젝트 설정 시간을 5분 이내로 단축

---

## 2. Problem Statement

| 현재 상태 | 목표 상태 |
|-----------|-----------|
| Go + Tailwind + templ + htmx 프로젝트 설정에 2-4시간 소요 | 5분 이내 완료 |
| 매번 동일한 보일러플레이트 작성 | 자동 생성 |
| 핫 리로드 설정 복잡 | 단일 명령어로 실행 |
| 아키텍처 패턴 일관성 부재 | Clean Architecture 강제 |

---

## 3. Functional Requirements

### 3.1 CLI 명령어

| ID | 명령어 | 설명 | 우선순위 |
|----|--------|------|----------|
| FR-001 | `protem-gen create` | 대화형 프로젝트 생성 | P0 |
| FR-002 | `protem-gen create --name <name>` | 비대화형 생성 | P1 |
| FR-003 | `protem-gen version` | 버전 출력 | P0 |
| FR-004 | `protem-gen list-templates` | 사용 가능 템플릿 목록 | P2 |

### 3.2 프로젝트 설정 옵션

| ID | 옵션 | 값 | 기본값 | 우선순위 |
|----|------|-----|--------|----------|
| FR-010 | Project Name | string | 필수 입력 | P0 |
| FR-011 | Module Path | string | github.com/user/name | P0 |
| FR-012 | HTTP Framework | gin (고정) | gin | P0 |
| FR-013 | Database | postgres, sqlite, none | postgres | P0 |
| FR-014 | Include gRPC | boolean | false | P1 |
| FR-015 | Include Auth | boolean | false | P2 |
| FR-016 | Include AI Ready | boolean | false | P2 |

> **Note**: HTTP Framework는 Gin으로 고정됨. 향후 버전에서 추가 프레임워크 지원 예정

### 3.3 생성되는 프로젝트 기능

| ID | 기능 | 설명 | 우선순위 |
|----|------|------|----------|
| FR-020 | Hot Reload | air + templ + tailwind 통합 | P0 |
| FR-021 | Clean Architecture | 도메인 기반 디렉토리 구조 | P0 |
| FR-022 | sqlc Integration | 타입 안전 SQL 쿼리 | P0 |
| FR-023 | templ Templates | 타입 안전 HTML 템플릿 | P0 |
| FR-024 | Tailwind CSS | 유틸리티 우선 CSS | P0 |
| FR-025 | htmx Integration | 서버 렌더링 상호작용 | P0 |
| FR-026 | Alpine.js | 클라이언트 상태 관리 | P0 |
| FR-027 | Makefile | 개발 명령어 모음 | P0 |

---

## 4. Non-Functional Requirements

### 4.1 성능

| ID | 요구사항 | 기준 |
|----|----------|------|
| NFR-001 | 프로젝트 생성 시간 | < 10초 |
| NFR-002 | CLI 시작 시간 | < 500ms |
| NFR-003 | 바이너리 크기 | < 50MB |

### 4.2 호환성

| ID | 요구사항 | 지원 범위 |
|----|----------|-----------|
| NFR-010 | Go 버전 | 1.22+ |
| NFR-011 | OS | macOS, Linux, Windows |
| NFR-012 | 터미널 | ANSI 색상 지원 터미널 |

### 4.3 사용성

| ID | 요구사항 | 기준 |
|----|----------|------|
| NFR-020 | 설치 단계 | 1단계 (go install) |
| NFR-021 | 첫 실행까지 시간 | < 5분 |
| NFR-022 | 문서 참조 없이 사용 | 대화형 가이드 제공 |

---

## 5. Generated Project Structure

```
<project-name>/
├── cmd/
│   └── server/
│       └── main.go              # 애플리케이션 진입점
├── internal/
│   ├── domain/                  # 비즈니스 엔티티
│   │   └── user/
│   │       └── user.go
│   ├── application/             # 유스케이스/서비스
│   │   └── user/
│   │       └── service.go
│   ├── infrastructure/          # 외부 시스템
│   │   ├── database/
│   │   │   ├── db.go
│   │   │   └── sqlc/            # sqlc 생성 코드
│   │   └── http/
│   │       └── server.go
│   └── interfaces/              # 핸들러/컨트롤러
│       └── http/
│           ├── handler.go
│           └── routes.go
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
├── .air.toml                    # air 설정
├── .gitignore
├── go.mod
├── Makefile
├── package.json                 # Tailwind 의존성
└── README.md
```

---

## 6. User Workflow

### 6.1 설치

```bash
go install github.com/hippo-an/tiny-go-challenges/protem-gen@latest
```

### 6.2 프로젝트 생성

```bash
$ protem-gen create

? Project name: my-app
? Module path: github.com/myuser/my-app
? Database: [postgres]
? Include gRPC support? [N]

Creating project: my-app

Checking required tools...
  ✓ go (go1.23.0)
  ✓ npm (10.2.0)
  ✓ air (v1.61.0)
  ✓ templ (v0.3.865)

Creating directory structure...
Initializing project...
  go mod init github.com/myuser/my-app
  npm init
  npm install [tailwindcss @tailwindcss/cli]
  air init
Configuring project files...
Generating source files...
Installing dependencies...
  go get [github.com/gin-gonic/gin github.com/a-h/templ github.com/jackc/pgx/v5]
  go mod tidy

✓ Project 'my-app' created successfully!

Next steps:
  cd my-app
  make setup      # Install dependencies
  make dev        # Start development server
```

### 6.3 개발 시작

```bash
cd my-app
make setup    # npm install, go mod download
make dev      # 핫 리로드 개발 서버 시작
```

---

## 7. Makefile Commands (Generated)

| Command | Description |
|---------|-------------|
| `make dev` | 모든 워처 병렬 실행 (air + templ + tailwind) |
| `make build` | 프로덕션 빌드 |
| `make test` | 테스트 실행 |
| `make sqlc-generate` | sqlc 코드 생성 |
| `make templ-generate` | templ 코드 생성 |
| `make setup` | 의존성 설치 |
| `make clean` | 빌드 결과물 삭제 |

---

## 8. Technology Stack

### 8.1 Generator (protem-gen)

| Component | Technology | Version |
|-----------|------------|---------|
| Language | Go | 1.22+ |
| CLI Framework | Cobra | v1.8.x |
| TUI | Bubbletea | v1.2.x |
| Template Engine | text/template | stdlib |

### 8.2 Generated Project

| Component | Technology | Version |
|-----------|------------|---------|
| Language | Go | 1.23+ |
| HTTP Framework | Gin | v1.11.0 |
| Templating | templ | v0.3.x |
| CSS | Tailwind CSS | v4.x |
| Interactivity | htmx | v2.0.x |
| Client State | Alpine.js | v3.x |
| Database | PostgreSQL/SQLite | - |
| SQL Codegen | sqlc | v1.30.x |
| Hot Reload | air | latest |

---

## 9. Success Criteria

| ID | Criteria | Metric |
|----|----------|--------|
| SC-001 | 프로젝트 생성 성공률 | 100% (지원 환경) |
| SC-002 | `make dev` 성공률 | 100% |
| SC-003 | 핫 리로드 동작 | Go, templ, Tailwind 모두 |
| SC-004 | 생성 프로젝트 빌드 | 에러 없음 |
| SC-005 | 생성 프로젝트 테스트 | 예제 테스트 통과 |

---

## 10. Out of Scope (v1.0)

| Item | Reason | Future Version |
|------|--------|----------------|
| 웹 기반 UI | CLI 우선 | v2.0 |
| 플러그인 시스템 | 복잡성 | v2.0 |
| 프로젝트 업데이트 | 생성만 지원 | v1.5 |
| Docker Compose | 선택적 기능 | v1.2 |
| CI/CD 템플릿 | 선택적 기능 | v1.3 |
| Kubernetes manifests | 선택적 기능 | v1.3 |

---

## 11. Dependencies Matrix

### 11.1 External Dependencies (Generator)

| Dependency | Purpose | Risk Level |
|------------|---------|------------|
| spf13/cobra | CLI framework | Low |
| charmbracelet/bubbletea | TUI | Low |
| charmbracelet/lipgloss | Styling | Low |

### 11.2 External Dependencies (Generated)

| Dependency | Purpose | Risk Level |
|------------|---------|------------|
| gin-gonic/gin | HTTP framework | Low |
| a-h/templ | Templating | Medium |
| sqlc-dev/sqlc | SQL codegen | Low |
| air-verse/air | Hot reload | Low |

---

## 12. Risks & Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| templ API 변경 | 템플릿 재작성 | Medium | 버전 고정, 모니터링 |
| OS별 호환성 | 지원 범위 제한 | Low | CI 다중 OS 테스트 |
| 의존성 충돌 | 빌드 실패 | Medium | 버전 고정, 통합 테스트 |

---

## 13. Glossary

| Term | Definition |
|------|------------|
| Hot Reload | 코드 변경 시 자동 재빌드/새로고침 |
| Clean Architecture | 의존성 역전 기반 레이어 분리 아키텍처 |
| templ | Go용 타입 안전 HTML 템플릿 라이브러리 |
| htmx | 서버 렌더링 기반 SPA 유사 기능 라이브러리 |
| sqlc | SQL 쿼리에서 타입 안전 Go 코드 생성 도구 |
| air | Go 애플리케이션 핫 리로드 도구 |

---

## Appendix A: Command Reference

```bash
# 설치
go install github.com/hippo-an/tiny-go-challenges/protem-gen@latest

# 버전 확인
protem-gen version

# 대화형 프로젝트 생성
protem-gen create

# 비대화형 프로젝트 생성
protem-gen create \
  --name my-app \
  --module github.com/user/my-app \
  --database postgres \
  --grpc=false
```

---

## Appendix B: Generated Project Quick Start

```bash
# 1. 프로젝트 생성
protem-gen create
# Project name: my-app

# 2. 디렉토리 이동
cd my-app

# 3. 의존성 설치
make setup

# 4. 개발 서버 시작
make dev

# 5. 브라우저 열기
# http://localhost:8080
```
