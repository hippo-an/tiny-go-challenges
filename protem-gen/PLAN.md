# protem-gen: Go Web Application Template Generator

## Implementation Plan

> **Document Version**: 1.2.0
> **Last Updated**: 2025-11-27
> **Status**: Draft - Pending Review

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
> Database type: [postgres/mysql/sqlite]: postgres

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

**결정**: `Gin` 프레임워크 기본값 (옵션으로 `Echo`, `Chi`, `Fiber` 제공)

#### 프레임워크 비교 분석 (2025년 기준)

**GitHub Stars & 커뮤니티**:
| Framework | GitHub Stars | 첫 릴리즈 | 최근 활동 |
|-----------|--------------|-----------|-----------|
| **Gin** | 81k+ | 2014 | 활발 |
| **Echo** | 30k+ | 2015 | 활발 |
| **Fiber** | 35k+ | 2020 | 매우 활발 |
| **Chi** | 19k+ | 2016 | 활발 |

**성능 벤치마크** (Go 1.23.5, MacBook Pro M2 기준):
| Framework | Requests/sec | Latency (p99) | 메모리 사용량 |
|-----------|--------------|---------------|---------------|
| **Fiber** | 최고 | 최저 | 낮음 |
| **Echo** | 우수 | 낮음 | 낮음 |
| **Gin** | 우수 | 낮음 | 낮음 |
| **Chi** | 우수 | 낮음 | 최저 |

> **참고**: 2025년 벤치마크에서 실제 워크로드(DB 연동, JSON 처리, 미들웨어 체인) 테스트 시 프레임워크 간 성능 차이는 미미함

**기능 비교**:
| Feature | Gin | Echo | Fiber | Chi |
|---------|-----|------|-------|-----|
| HTTP/2 지원 | ✅ | ✅ | ✅ | ✅ |
| WebSocket | ✅ | ✅ | ✅ (내장) | 외부 |
| 미들웨어 생태계 | 풍부 | 풍부 | 풍부 | 표준 호환 |
| JSON 검증 | 내장 | 내장 | 내장 | 외부 |
| 라우트 그룹화 | ✅ | ✅ | ✅ | ✅ |
| net/http 호환 | ✅ | ✅ | ❌ (fasthttp) | ✅ |
| 자동 TLS | ✅ | ✅ | ✅ | 외부 |

**상세 분석**:

| Framework | 장점 | 단점 | 적합한 상황 |
|-----------|------|------|-------------|
| **Gin** | • 가장 큰 커뮤니티 (81k stars)<br>• 풍부한 문서/튜토리얼<br>• Martini 대비 40x 빠름<br>• 안정적인 API | • 기능이 미니멀<br>• 추가 라이브러리 필요할 수 있음 | API, 마이크로서비스, 입문자 |
| **Echo** | • 타입 안전성 강조<br>• 깔끔한 API 설계<br>• 중앙 집중 에러 처리<br>• 다양한 템플릿 엔진 지원 | • 학습 곡선 약간 높음<br>• Gin보다 작은 커뮤니티 | 엔터프라이즈, API 게이트웨이 |
| **Fiber** | • 최고 성능 (fasthttp 기반)<br>• Express.js 스타일 API<br>• Node.js 개발자 친화적<br>• 내장 Rate Limiting | • net/http 비호환<br>• 일부 라이브러리 호환 문제<br>• 상대적으로 신생 | 고성능 요구, Node.js 전환자 |
| **Chi** | • 가장 가벼움<br>• net/http 완전 호환<br>• 표준 라이브러리 철학<br>• 컴포저블 설계 | • 내장 기능 최소<br>• 많은 외부 의존성 필요 | 미니멀리스트, stdlib 선호자 |

**권장 선택 가이드**:

```
사용 케이스별 권장 프레임워크:

┌─────────────────────────────────┬──────────────┐
│ 사용 케이스                      │ 권장         │
├─────────────────────────────────┼──────────────┤
│ 일반 웹 API / 마이크로서비스      │ Gin (기본값)  │
│ 엔터프라이즈 / 대규모 팀          │ Echo         │
│ 최고 성능 요구 / Node.js 전환    │ Fiber        │
│ stdlib 호환성 / 미니멀 설계       │ Chi          │
└─────────────────────────────────┴──────────────┘
```

**최종 결정: `Gin` 기본값 선택 근거**:
1. **커뮤니티 규모**: 81k+ stars로 압도적 1위, 문제 해결 리소스 풍부
2. **생태계**: 가장 많은 미들웨어, 플러그인, 예제 코드 존재
3. **안정성**: 10년+ 역사, 프로덕션 검증 완료
4. **학습 용이성**: 초보자도 빠르게 습득 가능
5. **성능**: 실제 워크로드에서 충분한 성능 (프레임워크 간 차이 미미)

> **참고 자료**:
> - [LogRocket: Top Go Frameworks 2025](https://blog.logrocket.com/top-go-frameworks-2025/)
> - [Tech Tonic: Go Framework Performance 2025](https://medium.com/deno-the-complete-reference/go-the-fastest-web-framework-in-2025-dfa2ddfd09e9)
> - [Fiber Benchmarks](https://docs.gofiber.io/extra/benchmarks/)

### ADR-004: Project Structure Pattern

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
- [ ] PLAN.md 승인
- [ ] PRD.md 승인
- [ ] go.mod 의존성 정의

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
    Framework   string   // echo, chi, fiber, stdlib
    Database    string   // postgres, mysql, sqlite, none
    IncludeGRPC bool     // gRPC 포함 여부
    IncludeAuth bool     // 인증 보일러플레이트
    IncludeAI   bool     // AI 통합 준비 코드
}
```

**완료 기준**:
- [ ] `protem-gen version` 동작
- [ ] `protem-gen create` 대화형 프롬프트 동작
- [ ] 빈 디렉토리 생성 확인

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
│   ├── echo/
│   │   └── server.go.tmpl
│   ├── chi/
│   │   └── server.go.tmpl
│   └── fiber/
│       └── server.go.tmpl
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
│   │   ├── schema.sql.tmpl
│   │   └── queries/
│   ├── mysql/
│   └── sqlite/
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
- [ ] 모든 템플릿 파일 작성 완료
- [ ] 템플릿 렌더링 테스트 통과
- [ ] 기본 프로젝트 생성 가능

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

# Tailwind CSS 감시
.PHONY: watch-tailwind
watch-tailwind:
	npx tailwindcss -i ./web/tailwind/input.css -o ./web/static/css/output.css --watch

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
- [ ] `make dev` 단일 명령어로 모든 워처 실행
- [ ] Go 코드 변경 시 자동 재빌드
- [ ] templ 파일 변경 시 자동 재생성
- [ ] Tailwind 변경 시 CSS 자동 컴파일

---

### Phase 4: Database Integration
**목표**: sqlc 기반 데이터베이스 레이어 구현

#### 4.1 sqlc 설정

```yaml
# sqlc.yaml 템플릿
version: "2"
sql:
  - engine: "postgresql"  # 또는 mysql, sqlite
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
// 생성될 go.mod - HTTP 프레임워크 선택에 따라 변경
require (
    // HTTP 프레임워크 (택1)
    github.com/gin-gonic/gin v1.11.0          // 기본값
    // github.com/labstack/echo/v4 v4.13.4    // 옵션
    // github.com/go-chi/chi/v5 v5.2.0        // 옵션
    // github.com/gofiber/fiber/v2 v2.52.6    // 옵션 (net/http 비호환)

    // 템플릿 & 프론트엔드
    github.com/a-h/templ v0.3.865             // HTML 템플릿 (Go)

    // 데이터베이스
    github.com/jackc/pgx/v5 v5.7.5            // PostgreSQL
    // github.com/go-sql-driver/mysql v1.8.1  // MySQL (옵션)
    // modernc.org/sqlite v1.34.5             // SQLite (옵션)
)

// 개발 도구 (go install)
// github.com/air-verse/air@latest            // Go 핫 리로드
// github.com/a-h/templ/cmd/templ@v0.3.865    // templ CLI
// github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0  // SQL 코드 생성
```

| 패키지 | 버전 | 용도 | 비고 |
|--------|------|------|------|
| **HTTP 프레임워크** |
| gin-gonic/gin | v1.11.0 | HTTP 프레임워크 | 기본값, Go 1.23+ |
| labstack/echo/v4 | v4.13.4 | HTTP 프레임워크 | 옵션 |
| go-chi/chi/v5 | v5.2.0 | HTTP 라우터 | 옵션, stdlib 호환 |
| gofiber/fiber/v2 | v2.52.6 | HTTP 프레임워크 | 옵션, v3는 베타 |
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
    "tailwindcss": "^3.4.17",
    "@tailwindcss/forms": "^0.5.10",
    "@tailwindcss/typography": "^0.5.16"
  }
}
```

| 패키지 | 버전 | 용도 | 비고 |
|--------|------|------|------|
| tailwindcss | ^3.4.17 | CSS 프레임워크 | v4는 설정 방식 대폭 변경 |
| @tailwindcss/forms | ^0.5.10 | 폼 스타일 리셋 | |
| @tailwindcss/typography | ^0.5.16 | 타이포그래피 | |

### CDN 의존성 (생성될 프로젝트 HTML)

| 라이브러리 | 버전 | CDN URL |
|------------|------|---------|
| htmx | 2.0.4 | `https://unpkg.com/htmx.org@2.0.4` |
| Alpine.js | 3.14.8 | `https://cdn.jsdelivr.net/npm/alpinejs@3.14.8/dist/cdn.min.js` |

### 버전 선택 근거

1. **Tailwind CSS v3.4.x 선택 (v4 아님)**
   - v4는 2025년 1월 출시, CSS-first 설정으로 대폭 변경
   - v3.4.x는 안정적이고 문서/예제 풍부
   - v4 마이그레이션은 향후 버전에서 지원 예정

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
8. **Tailwind CSS v4**: CSS-first 설정 방식 지원

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.2.0 | 2025-11-27 | Technical Specifications 전면 업데이트: 최신 안정화 버전 반영, 버전 선택 근거 추가 | AI Assistant |
| 1.1.0 | 2025-11-27 | ADR-003 확장: Gin 추가 및 4개 프레임워크 상세 비교 분석 | AI Assistant |
| 1.0.0 | 2025-11-27 | Initial draft | AI Assistant |
