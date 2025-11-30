# protem-gen: Go Web Application Template Generator

## Implementation Plan

> **Document Version**: 2.4.0
> **Last Updated**: 2025-12-01
> **Status**: Implemented - Phase 3 Complete

---

## Table of Contents

1. [Project Vision](#project-vision)
2. [Architecture Decisions](#architecture-decisions)
3. [Implementation Phases](#implementation-phases)
4. [Technical Specifications](#technical-specifications)
5. [Risk Assessment](#risk-assessment)
6. [Version History](#version-history)

---

## Project Vision

### What We're Building

**protem-gen**은 Go 웹 애플리케이션 프로젝트를 생성하는 CLI 바이너리 도구입니다.

**핵심 구분점**: 이것은 복사-붙여넣기 템플릿이 아닌, 사용자 선택에 따라 커스터마이즈된 프로젝트를 생성하는 **실행 가능한 바이너리**입니다.

### Target Usage

```bash
# 설치
go install github.com/hippo-an/tiny-go-challenges/protem-gen@latest

# 실행
protem-gen create

# 대화형 프롬프트
> Project name: my-app
> Include gRPC? [y/N]: y
> Database type: [postgres/sqlite]: postgres

# 결과: ./my-app/ 에 완전한 프로젝트 생성
```

---

## Architecture Decisions

### ADR-001: CLI Framework Selection

**결정**: `cobra` + `bubbletea` 조합 사용

**근거**:
- `cobra`: Go CLI 표준, 서브커맨드 지원, 자동 완성 지원
- `bubbletea`: 풍부한 TUI 경험, 대화형 프롬프트에 적합

**대안 검토**:
| Framework | 장점 | 단점 |
|-----------|------|------|
| urfave/cli | 간단함 | TUI 지원 부족 |
| cobra only | 안정적 | 대화형 UI 제한적 |
| survey | 프롬프트 특화 | 유지보수 중단 상태 |

### ADR-002: Template Engine Selection

**결정**: Go `text/template` + `embed` 패키지 사용

**근거**:
- 외부 의존성 없음
- Go 표준 라이브러리로 안정성 보장
- `embed`로 바이너리에 템플릿 포함

### ADR-003: Generated Project HTTP Framework

**결정**: `Gin` 프레임워크 고정 (단일 프레임워크 지원)

**근거**:
1. **커뮤니티 규모**: 81k+ stars로 압도적 1위, 문제 해결 리소스 풍부
2. **생태계**: 가장 많은 미들웨어, 플러그인, 예제 코드 존재
3. **안정성**: 10년+ 역사, 프로덕션 검증 완료
4. **학습 용이성**: 초보자도 빠르게 습득 가능
5. **유지보수 단순화**: 단일 프레임워크 지원으로 코드 복잡도 감소

> **v2 계획**: 향후 버전에서 Echo, Chi 등 추가 프레임워크 지원 검토

### ADR-004: CLI 실행 기반 프로젝트 초기화

**결정**: 정적 템플릿 대신 실제 CLI 명령어 실행으로 프로젝트 초기화

**근거**:
- 라이브러리 업데이트 시 자동으로 최신 기본값 반영
- 템플릿 유지보수 부담 감소
- 실제 도구의 기본 설정 활용

**CLI 실행 대상**:
| 도구 | 명령어 | 결과물 |
|------|--------|--------|
| Go | `go mod init <module>` | go.mod |
| npm | `npm init -y` | package.json |
| npm | `npm install tailwindcss @tailwindcss/cli -D` | node_modules |
| air | `air init` | .air.toml |

> **Note**: Tailwind CSS v4는 `tailwindcss init` 명령어가 제거되었습니다. 대신 CSS-first 설정 방식을 사용합니다.

**템플릿 유지 대상**:
- Makefile, README.md, .gitignore (CLI 대안 없음)
- Go 소스 코드 (프레임워크별 커스텀 필요)
- templ 파일 (프로젝트별 커스텀 필요)

### ADR-005: Project Structure Pattern

**결정**: Clean Architecture + Domain-Driven Directory Layout

```
generated-project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/           # 엔티티, 비즈니스 규칙
│   │   └── user/
│   ├── application/      # 유스케이스, 서비스
│   │   └── user/
│   ├── infrastructure/   # 외부 시스템 연결
│   │   ├── database/
│   │   ├── http/
│   │   └── external/
│   └── interfaces/       # 핸들러, 프레젠테이션
│       ├── http/
│       └── grpc/
├── pkg/                  # 재사용 가능한 패키지
├── web/
│   ├── templates/        # templ 파일
│   ├── static/           # 정적 파일
│   └── tailwind/         # Tailwind 설정
├── migrations/           # 데이터베이스 마이그레이션
├── sqlc/                 # sqlc 설정 및 쿼리
└── Makefile
```

---

## Implementation Phases

### Phase 0: Foundation Setup (현재)
**기간**: Day 1-2
**목표**: 프로젝트 기반 구축

| Task | 상세 | 우선순위 |
|------|------|----------|
| 0.1 | PLAN.md 작성 | P0 |
| 0.2 | PRD.md 작성 | P0 |
| 0.3 | 문서 리뷰 및 승인 | P0 |
| 0.4 | 개발 환경 설정 | P0 |

**완료 기준**:
- [x] PLAN.md 승인
- [x] PRD.md 승인
- [x] go.mod 의존성 정의

---

### Phase 1: CLI Core Development
**목표**: 기본 CLI 인터페이스 구현

#### 1.1 프로젝트 구조 설정

```
protem-gen/
├── cmd/
│   └── protem-gen/
│       └── main.go
├── internal/
│   ├── cli/              # CLI 명령어
│   │   ├── create.go
│   │   └── version.go
│   ├── config/           # 설정 처리
│   │   └── config.go
│   ├── generator/        # 프로젝트 생성 로직
│   │   ├── generator.go
│   │   └── scaffold.go
│   ├── prompt/           # 대화형 프롬프트
│   │   └── prompt.go
│   └── template/         # 템플릿 처리
│       └── template.go
├── templates/            # 임베드될 템플릿 파일들
│   ├── base/
│   ├── http/
│   ├── grpc/
│   └── database/
├── PLAN.md
├── PRD.md
├── README.md
├── Makefile
└── go.mod
```

#### 1.2 구현 순서

| Step | Task | 의존성 | 산출물 |
|------|------|--------|--------|
| 1.2.1 | Cobra CLI 스캐폴딩 | 없음 | `cmd/protem-gen/main.go` |
| 1.2.2 | `create` 서브커맨드 구현 | 1.2.1 | `internal/cli/create.go` |
| 1.2.3 | 설정 구조체 정의 | 없음 | `internal/config/config.go` |
| 1.2.4 | Bubbletea 프롬프트 구현 | 1.2.2 | `internal/prompt/prompt.go` |
| 1.2.5 | 기본 생성기 로직 | 1.2.3 | `internal/generator/generator.go` |

#### 1.3 프롬프트 설계

```go
type ProjectConfig struct {
    Name        string   // 프로젝트 이름
    ModulePath  string   // Go 모듈 경로
    Database    string   // postgres, sqlite, none
    IncludeGRPC bool     // gRPC 포함 여부
    IncludeAuth bool     // 인증 보일러플레이트
    IncludeAI   bool     // AI 통합 준비 코드
}
// Note: Framework is fixed to Gin, not configurable
```

**완료 기준**:
- [x] `protem-gen version` 동작
- [x] `protem-gen create` 대화형 프롬프트 동작
- [x] 빈 디렉토리 생성 확인

---

### Phase 2: Template System Implementation
**목표**: 프로젝트 템플릿 시스템 구현

#### 2.1 템플릿 구조

```
templates/
├── base/                    # 모든 프로젝트 공통
│   ├── go.mod.tmpl
│   ├── Makefile.tmpl
│   ├── README.md.tmpl
│   ├── .gitignore.tmpl
│   ├── .air.toml.tmpl
│   └── cmd/
│       └── server/
│           └── main.go.tmpl
├── http/
│   └── gin/
│       ├── server.go.tmpl
│       ├── routes.go.tmpl
│       └── handler.go.tmpl
├── architecture/
│   ├── domain/
│   │   └── user.go.tmpl
│   ├── application/
│   │   └── user_service.go.tmpl
│   ├── infrastructure/
│   │   └── user_repository.go.tmpl
│   └── interfaces/
│       └── user_handler.go.tmpl
├── database/
│   ├── postgres/
│   │   ├── sqlc.yaml.tmpl
│   │   ├── db.go.tmpl
│   │   ├── schema.sql.tmpl
│   │   └── queries.sql.tmpl
│   └── sqlite/
│       ├── sqlc.yaml.tmpl
│       ├── db.go.tmpl
│       ├── schema.sql.tmpl
│       └── queries.sql.tmpl
├── frontend/
│   ├── tailwind.config.js.tmpl
│   ├── package.json.tmpl
│   └── templates/
│       ├── base.templ.tmpl
│       ├── index.templ.tmpl
│       └── components/
├── grpc/                    # 선택적
│   ├── proto/
│   └── server/
└── ai/                      # 선택적
    ├── llm/
    │   └── client.go.tmpl
    └── prompt/
        └── manager.go.tmpl
```

#### 2.2 구현 순서

| Step | Task | 의존성 | 산출물 |
|------|------|--------|--------|
| 2.2.1 | 템플릿 임베딩 설정 | Phase 1 | `internal/template/embed.go` |
| 2.2.2 | 템플릿 렌더링 엔진 | 2.2.1 | `internal/template/render.go` |
| 2.2.3 | 기본 템플릿 작성 | 2.2.1 | `templates/base/*` |
| 2.2.4 | HTTP 프레임워크 템플릿 | 2.2.3 | `templates/http/*` |
| 2.2.5 | 아키텍처 템플릿 | 2.2.3 | `templates/architecture/*` |
| 2.2.6 | 데이터베이스 템플릿 | 2.2.3 | `templates/database/*` |
| 2.2.7 | 프론트엔드 템플릿 | 2.2.3 | `templates/frontend/*` |
| 2.2.8 | 생성기와 템플릿 통합 | 2.2.2-7 | `internal/generator/scaffold.go` |

**완료 기준**:
- [x] 모든 템플릿 파일 작성 완료
- [x] 템플릿 렌더링 테스트 통과
- [x] 기본 프로젝트 생성 가능 (PostgreSQL, SQLite, None 모두 빌드 성공)

---

### Phase 3: Hot Reload Pipeline
**목표**: 생성된 프로젝트의 핫 리로드 환경 구현

#### 3.1 핫 리로드 아키텍처

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   air       │    │   templ     │    │  tailwind   │
│  (Go 코드)   │    │  (템플릿)    │    │   (CSS)     │
└──────┬──────┘    └──────┬──────┘    └──────┬──────┘
       │                  │                   │
       │   .go 변경       │   .templ 변경     │   .css 변경
       ▼                  ▼                   ▼
┌─────────────────────────────────────────────────────┐
│                    make dev                          │
│  (모든 워처를 병렬로 실행)                              │
└─────────────────────────────────────────────────────┘
```

#### 3.2 Makefile 명령어 (생성될 프로젝트용)

```makefile
# 개발 환경 실행 (모든 워처 병렬 실행)
.PHONY: dev
dev:
	@echo "Starting development environment..."
	@make -j3 watch-air watch-templ watch-tailwind

# Go 코드 핫 리로드
.PHONY: watch-air
watch-air:
	air -c .air.toml

# templ 파일 감시 및 재생성
.PHONY: watch-templ
watch-templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false

# Tailwind CSS 감시 (v4 uses @tailwindcss/cli)
.PHONY: watch-tailwind
watch-tailwind:
	npx @tailwindcss/cli -i ./web/tailwind/input.css -o ./web/static/css/output.css --watch

# sqlc 코드 생성
.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate

# 빌드
.PHONY: build
build:
	go build -o bin/server ./cmd/server

# 테스트
.PHONY: test
test:
	go test -v ./...
```

#### 3.3 구현 순서

| Step | Task | 의존성 | 산출물 |
|------|------|--------|--------|
| 3.3.1 | air 설정 템플릿 작성 | Phase 2 | `.air.toml.tmpl` |
| 3.3.2 | templ 워처 설정 | Phase 2 | Makefile 업데이트 |
| 3.3.3 | Tailwind 설정 템플릿 | Phase 2 | `tailwind.config.js.tmpl` |
| 3.3.4 | package.json 템플릿 | 3.3.3 | `package.json.tmpl` |
| 3.3.5 | Makefile 템플릿 완성 | 3.3.1-4 | `Makefile.tmpl` |
| 3.3.6 | 통합 테스트 | 3.3.5 | 테스트 프로젝트 생성 및 검증 |

**완료 기준**:
- [x] `make dev` 단일 명령어로 모든 워처 실행
- [x] Go 코드 변경 시 자동 재빌드
- [x] templ 파일 변경 시 자동 재생성
- [x] Tailwind 변경 시 CSS 자동 컴파일

---

### Phase 4: Database Integration
**목표**: sqlc 기반 데이터베이스 레이어 구현

#### 4.1 sqlc 설정

```yaml
# sqlc.yaml 템플릿
version: "2"
sql:
  - engine: "postgresql"  # 또는 sqlite
    queries: "sqlc/queries/"
    schema: "migrations/"
    gen:
      go:
        package: "db"
        out: "internal/infrastructure/database/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
```

#### 4.2 구현 순서

| Step | Task | 의존성 | 산출물 |
|------|------|--------|--------|
| 4.2.1 | sqlc.yaml 템플릿 (DB별) | Phase 2 | `templates/database/*/sqlc.yaml.tmpl` |
| 4.2.2 | 기본 스키마 템플릿 | 4.2.1 | `schema.sql.tmpl` |
| 4.2.3 | 예제 쿼리 템플릿 | 4.2.2 | `queries/*.sql.tmpl` |
| 4.2.4 | DB 연결 설정 | 4.2.1 | `db.go.tmpl` |
| 4.2.5 | 리포지토리 패턴 구현 | 4.2.3 | `repository.go.tmpl` |

**완료 기준**:
- [ ] `make sqlc-generate` 동작
- [ ] 생성된 코드 컴파일 성공
- [ ] 예제 CRUD 작업 동작

---

### Phase 5: Frontend Integration
**목표**: templ + htmx + Alpine.js + Tailwind 통합

#### 5.1 프론트엔드 스택 구성

```
web/
├── templates/
│   ├── layouts/
│   │   └── base.templ      # 기본 레이아웃 (htmx, Alpine 로드)
│   ├── pages/
│   │   ├── index.templ     # 홈페이지
│   │   └── about.templ
│   └── components/
│       ├── header.templ
│       ├── footer.templ
│       └── button.templ    # 재사용 컴포넌트
├── static/
│   ├── css/
│   │   └── output.css      # Tailwind 출력
│   └── js/
│       └── app.js          # 커스텀 JS (필요시)
└── tailwind/
    ├── input.css           # Tailwind 입력
    └── tailwind.config.js
```

#### 5.2 구현 순서

| Step | Task | 의존성 | 산출물 |
|------|------|--------|--------|
| 5.2.1 | 기본 레이아웃 템플릿 | Phase 2 | `base.templ.tmpl` |
| 5.2.2 | htmx 통합 예제 | 5.2.1 | 동적 컴포넌트 예제 |
| 5.2.3 | Alpine.js 통합 예제 | 5.2.1 | 상태 관리 예제 |
| 5.2.4 | Tailwind 컴포넌트 | 5.2.1 | UI 컴포넌트 템플릿 |
| 5.2.5 | 핸들러와 템플릿 연결 | 5.2.1-4 | 라우트 설정 |

**완료 기준**:
- [ ] templ 파일 렌더링 성공
- [ ] htmx 부분 업데이트 동작
- [ ] Alpine.js 상호작용 동작
- [ ] Tailwind 스타일 적용

---

### Phase 6: Optional Features
**목표**: 선택적 기능 구현

#### 6.1 gRPC 지원

| Step | Task | 산출물 |
|------|------|--------|
| 6.1.1 | Proto 파일 템플릿 | `proto/*.proto.tmpl` |
| 6.1.2 | gRPC 서버 설정 | `grpc/server.go.tmpl` |
| 6.1.3 | buf 설정 | `buf.yaml.tmpl` |
| 6.1.4 | Makefile 명령어 추가 | `make proto-generate` |

#### 6.2 AI 통합 준비

| Step | Task | 산출물 |
|------|------|--------|
| 6.2.1 | LLM 클라이언트 인터페이스 | `llm/client.go.tmpl` |
| 6.2.2 | 프롬프트 매니저 | `prompt/manager.go.tmpl` |
| 6.2.3 | 스트리밍 응답 핸들러 | `stream/handler.go.tmpl` |
| 6.2.4 | API 키 관리 패턴 | `config/secrets.go.tmpl` |

#### 6.3 인증 보일러플레이트

| Step | Task | 산출물 |
|------|------|--------|
| 6.3.1 | 세션 관리 | `auth/session.go.tmpl` |
| 6.3.2 | JWT 유틸리티 | `auth/jwt.go.tmpl` |
| 6.3.3 | 미들웨어 | `middleware/auth.go.tmpl` |

---

### Phase 7: Testing & Documentation
**목표**: 테스트 및 문서화

#### 7.1 테스트 전략

| 테스트 유형 | 대상 | 도구 |
|------------|------|------|
| 단위 테스트 | 생성기 로직 | Go testing |
| 통합 테스트 | 전체 생성 파이프라인 | Go testing |
| E2E 테스트 | 생성된 프로젝트 | 스크립트 |
| 골든 테스트 | 템플릿 출력 | testscript |

#### 7.2 구현 순서

| Step | Task | 산출물 |
|------|------|--------|
| 7.2.1 | 단위 테스트 작성 | `*_test.go` |
| 7.2.2 | 통합 테스트 작성 | `integration_test.go` |
| 7.2.3 | E2E 테스트 스크립트 | `scripts/e2e_test.sh` |
| 7.2.4 | README.md 작성 | 생성기용 README |
| 7.2.5 | 생성 프로젝트 README 템플릿 | `README.md.tmpl` |

**완료 기준**:
- [ ] 테스트 커버리지 80% 이상
- [ ] 모든 테스트 통과
- [ ] 문서화 완료

---

### Phase 8: Release & Distribution
**목표**: 배포 준비

#### 8.1 배포 전략

| Step | Task | 산출물 |
|------|------|--------|
| 8.1.1 | GoReleaser 설정 | `.goreleaser.yml` |
| 8.1.2 | GitHub Actions CI/CD | `.github/workflows/` |
| 8.1.3 | 버전 관리 | 시맨틱 버저닝 |
| 8.1.4 | Homebrew Formula | homebrew-tap (선택적) |

---

## Technical Specifications

> **버전 검증일**: 2025-11-27
> **Go 최소 요구 버전**: 1.23+

### 의존성 버전 (생성기: protem-gen)

```go
// go.mod (protem-gen)
require (
    github.com/spf13/cobra v1.9.1              // CLI 프레임워크
    github.com/charmbracelet/bubbletea v1.2.4  // TUI 프레임워크
    github.com/charmbracelet/lipgloss v1.0.0   // 터미널 스타일링
    github.com/charmbracelet/bubbles v0.20.0   // TUI 컴포넌트
)
```

| 패키지 | 버전 | 용도 | 비고 |
|--------|------|------|------|
| spf13/cobra | v1.9.1 | CLI 프레임워크 | Go 1.22+ 권장 |
| charmbracelet/bubbletea | v1.2.4 | TUI 프레임워크 | v2는 베타 |
| charmbracelet/lipgloss | v1.0.0 | 터미널 스타일링 | v2는 베타 |
| charmbracelet/bubbles | v0.20.0 | TUI 컴포넌트 | v2는 베타 |

### 의존성 버전 (생성될 프로젝트)

```go
// 생성될 go.mod - go mod init 및 go get으로 자동 생성
require (
    // HTTP 프레임워크 (Gin 고정)
    github.com/gin-gonic/gin v1.11.0

    // 템플릿 & 프론트엔드
    github.com/a-h/templ v0.3.865             // HTML 템플릿 (Go)

    // 데이터베이스 (선택)
    github.com/jackc/pgx/v5 v5.7.5            // PostgreSQL
    // modernc.org/sqlite v1.34.5             // SQLite (옵션)
)

// 개발 도구 (go install) - 필수 도구, 없으면 에러
// github.com/air-verse/air@latest            // Go 핫 리로드
// github.com/a-h/templ/cmd/templ@latest      // templ CLI
// github.com/sqlc-dev/sqlc/cmd/sqlc@latest   // SQL 코드 생성
```

| 패키지 | 버전 | 용도 | 비고 |
|--------|------|------|------|
| **HTTP 프레임워크** |
| gin-gonic/gin | v1.11.0 | HTTP 프레임워크 | 고정, Go 1.23+ |
| **템플릿** |
| a-h/templ | v0.3.865 | HTML 템플릿 | |
| **데이터베이스** |
| jackc/pgx/v5 | v5.7.5 | PostgreSQL 드라이버 | |
| sqlc-dev/sqlc | v1.30.0 | SQL 코드 생성 | CLI 도구 |
| **개발 도구** |
| air-verse/air | latest | Go 핫 리로드 | cosmtrek/air에서 이전됨 |

### 생성될 프로젝트 npm 의존성

```json
{
  "devDependencies": {
    "tailwindcss": "^4.1.17",
    "@tailwindcss/cli": "^4.1.17"
  }
}
```

| 패키지 | 버전 | 용도 | 비고 |
|--------|------|------|------|
| tailwindcss | ^4.1.17 | CSS 프레임워크 | v4 CSS-first 설정 방식 |
| @tailwindcss/cli | ^4.1.17 | Tailwind CLI | 빌드 명령어용 |

> **Note**: Tailwind CSS v4는 CSS-first 설정 방식을 사용합니다. `tailwind.config.js` 대신 `input.css`에서 `@import "tailwindcss"` 및 `@source` 디렉티브로 설정합니다. forms, typography 플러그인은 v4에 기본 포함됩니다.

### CDN 의존성 (생성될 프로젝트 HTML)

| 라이브러리 | 버전 | CDN URL |
|------------|------|---------|
| htmx | 2.0.4 | `https://unpkg.com/htmx.org@2.0.4` |
| Alpine.js | 3.14.8 | `https://cdn.jsdelivr.net/npm/alpinejs@3.14.8/dist/cdn.min.js` |

### 버전 선택 근거

1. **Tailwind CSS v4.x 선택**
   - v4는 2025년 1월 출시, CSS-first 설정으로 대폭 변경
   - JavaScript 설정 파일 불필요, CSS 내 `@import "tailwindcss"` 사용
   - `@source` 디렉티브로 스캔 경로 지정
   - forms, typography 플러그인 기본 포함

2. **Charmbracelet v1.x 선택 (v2 베타 아님)**
   - v2는 현재 베타 상태
   - v1은 프로덕션 검증된 안정 버전
   - v2 정식 출시 후 업그레이드 검토

3. **Fiber v2 선택 (v3 베타 아님)**
   - v3는 현재 beta.4 상태
   - v2.52.6은 안정적인 최신 버전

4. **air 저장소 변경**
   - `cosmtrek/air` → `air-verse/air`로 이전됨
   - 새 저장소에서 유지보수 진행 중

---

## Risk Assessment

### 높은 위험

| 위험 | 영향 | 완화 전략 |
|------|------|-----------|
| 의존성 버전 충돌 | 생성된 프로젝트 빌드 실패 | 모든 의존성 버전 고정, 통합 테스트 |
| 템플릿 복잡성 증가 | 유지보수 어려움 | 모듈화, 템플릿 조합 패턴 사용 |

### 중간 위험

| 위험 | 영향 | 완화 전략 |
|------|------|-----------|
| OS 호환성 | Windows에서 Makefile 제한 | Cross-platform 스크립트 대안 제공 |
| 핫 리로드 동기화 | 워처 간 충돌 | 순차적 재시작 로직 |

### 낮은 위험

| 위험 | 영향 | 완화 전략 |
|------|------|-----------|
| TUI 호환성 | 일부 터미널에서 UI 깨짐 | 폴백 텍스트 모드 제공 |

---

## Extension Points

### 향후 확장 가능 영역

1. **추가 프레임워크**: Gorilla Mux, Hertz
2. **추가 데이터베이스**: MongoDB, Redis, CockroachDB
3. **인증 프로바이더**: OAuth2, OIDC, Passkey
4. **결제 통합**: Stripe, PayPal
5. **관측성**: OpenTelemetry, Prometheus, Grafana
6. **컨테이너화**: Docker, Kubernetes manifests
7. **CI/CD 템플릿**: GitHub Actions, GitLab CI

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 2.4.0 | 2025-12-01 | Phase 3 완료: Hot Reload Pipeline 검증 완료, .air.toml pre_cmd 버그 수정 | AI Assistant |
| 2.3.0 | 2025-11-30 | Phase 2 완료 검증: templ generate 자동화 추가, 모든 DB 옵션 빌드 테스트 통과 | AI Assistant |
| 2.2.0 | 2025-11-29 | Phase 2 완료: 아키텍처/DB 템플릿 추가, MySQL 제거 (PostgreSQL/SQLite만 지원) | AI Assistant |
| 2.1.0 | 2025-11-28 | Tailwind CSS v4 지원: CSS-first 설정 방식, @tailwindcss/cli 사용, Phase 0/1 완료 기준 마킹 | AI Assistant |
| 2.0.0 | 2025-11-27 | **Breaking**: Gin 단일 프레임워크 지원, CLI 실행 기반 프로젝트 초기화 (ADR-004), Chi/Fiber/Echo 제거 | AI Assistant |
| 1.2.0 | 2025-11-27 | Technical Specifications 전면 업데이트: 최신 안정화 버전 반영, 버전 선택 근거 추가 | AI Assistant |
| 1.1.0 | 2025-11-27 | ADR-003 확장: Gin 추가 및 4개 프레임워크 상세 비교 분석 | AI Assistant |
| 1.0.0 | 2025-11-27 | Initial draft | AI Assistant |
