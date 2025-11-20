# AGENTS.md

## 1. 공통 규칙

1. 모든 답변은 **한글**로 작성해.
2. 불필요한 복잡성보다 **직관적이고 명확한 코드**를 우선해.
3. 공통 요소가 충분히 반복되고 효과가 명확할 때만 **리팩토링**을 진행해.
4. 모든 화면과 UI는 **모바일 우선**(Mobile-first)으로 개발해.

---

## 2. 프로젝트 구성

- **언어**: Go
- **서버 프레임워크**: Echo
- **템플릿 엔진**: Templ (**Gomponents는 사용하지 않음**)
- **데이터베이스**: SQLite
- **DB 마이그레이션**: Goose
- **SQL 코드 생성기**: SQLC
- **프론트엔드 라이브러리**: HTMX, Alpine.js
- **CSS 프레임워크**: BeerCSS(기본), TailwindCSS(homepage만)
- **개발 서버**: GoVisual
- **소스 구조**: Monorepo

---

## 3. 프론트엔드 라이브러리 역할 정의

### HTMX

- **역할**: AJAX 요청으로 서버에서 HTML fragment를 받아 특정 영역만 동적으로 갱신.
- **목적**: 전체 페이지 리로드 없이 SPA와 유사한 사용자 경험 제공.

### Alpine.js

- **역할**: 클라이언트 측에서 경량 상태 관리 및 UI 인터랙션 처리.
- **목적**: 모달/드롭다운 등 즉각적 UI 동작 구현.

### HTMX × Alpine 적용 패턴

- 여러 프로젝트에서 재사용되는 HTML 조각은
  `shared/views` 또는 `shared/static/js`로 **공통화**.
- Alpine 상태는 반드시 **클라이언트에서만 결정 가능한 미세 상호작용**에 사용.
- 서버 데이터 기반 UI 상태는 **HTMX 응답**으로 갱신.

### Templ

- **역할**: HTML과 유사한 문법으로 UI를 작성하면, 정적 Go 코드로 컴파일하는 템플릿 엔진.
- **목적**: `.templ`에서 작성한 마크업을 `templ generate`로 미리 Go 코드로 변환하여
  런타임 파싱 비용 없이 빠르고 타입 안정적인 렌더링을 제공.
- **공통 컴포넌트 위치**: `simple-server/shared/views`

---

## 4. CSS 라이브러리 사용 정의

### BeerCSS (기본)

- HTML 안에 CSS 클래스를 직접 적용해 UI 구성.
- 모달/사이드바 등의 UI 동작은 `data-ui="#id"`만으로 제어 가능 → Alpine.js 불필요.
- 문서: https://beercss.com/
- 사용 가이드: `.doc/css/beercss/SUMMARY.md`
- ❗ **커스텀 CSS 금지**, BeerCSS에서 제공하지 않는 클래스도 사용 금지.

### TailwindCSS (homepage 전용)

- 설치: `task.sh install-tailwind`

---

## 5. 폴더 구조

```
cmd/{project}/main.go # 서버 엔트리포인트
projects/{project}/ # 프로젝트 소스
├─ handlers/ # Echo 핸들러
├─ views/ # Templ 템플릿
├─ static/ # 정적 파일
├─ migrations/ # Goose 마이그레이션
└─ query.sql # SQLC 쿼리
internal/ # 외부 노출 안 되는 공용 서버 코드
└─ validate/ # 검증 패키지
shared/ # 공통 프론트엔드 자원
├─ views/ # 공통 Templ
└─ static/ # 공통 JS, CSS
pkg/util/ # 순수 유틸 함수
```

---

## 6. 코딩 스타일 가이드

### 6.1 Handler 패턴

1. **바인딩 & 검증**

   - `c.Bind(&dto)`
   - `c.Validate(&dto)`
   - 실패 시 → `validate.HTTPError` 또는 `echo.NewHTTPError`

2. **DB 접근**

   - `queries, err := db.GetQueries()`
   - `queries.Method(ctx, params)`

3. **응답**
   - 정상:
     `views.Component(...).Render(ctx, c.Response().Writer)`
   - 오류:
     `echo.NewHTTPError(code, message)`

---

### 6.2 에러 핸들링 & 로깅

- 에러 처리는 Go 표준 방식: `if err != nil { ... }`
- 에러 응답: `echo.NewHTTPError`
- 로깅: `slog`
  - `DEBUG`, `INFO`, `WARN`, `ERROR`
  - 메시지는 **한글로 명확하게**

---

### 6.3 Templ 작성 가이드

- `.templ`에서 HTML과 동일한 구조로 작성.
- 서버 주도 렌더링이 기본값.

**동적 기능 규칙**

- **HTMX**로 데이터 기반 UI 동작 처리
  (`hx-get`, `hx-post`, `hx-target`, `hx-swap`)
- **Alpine.js**는 클라이언트 상태가 필요한 경우에만 사용
  (모달/탭/드롭다운 등)
- JS는 항상 **별도 .js 파일에 분리**.

**중요**

- `.templ` 파일을 수정했으면 반드시
  → `templ generate` 실행해 Go 코드 재생성.

---

## 7. HTTP / HTMX 응답 규칙

- **204 No Content**: 저장/수정/삭제 성공 (스왑 방지)
- **201 Created**: 생성된 리소스
- **202 Accepted**: 비동기 작업 접수
- **200 OK + HTML**: 부분 갱신
- **HX-Redirect** 헤더로 리다이렉션 처리
- 오류: 적절한 4xx/5xx + 한글 메시지

⚠ `200 OK` + **빈 본문** → HTMX가 대상 영역을 비워버림
→ 의도적으로 본문이 없을 때는 **204** 사용

---

## 8. 런타임 모드

### 개발 모드

- GoVisual 사용
- 정적 파일: 로컬 디스크
- 기능: `TransferEchoToGoVisualServerOnlyDev`

### 운영 모드

- Echo 단독 실행
- 정적 파일: embed.FS
- Gzip 활성화
- 타임아웃: 1분
- 요청 바디 제한: 5MB
- RateLimit: 초당 20회

---

## 9. 정적/임베드 자원

- 공통 정적: `/shared/static` → `shared/static/*`
- 프로젝트 정적: `/static` → `projects/{name}/static/*`
- PWA 파일 매핑:
  - `/manifest.json`
  - `/firebase-messaging-sw.js`
- Dev: `os.DirFS`
- Prod: `resources.EmbeddedFiles` (embed)

---

## 10. 인증 및 권한

- Firebase ID 토큰 → `/create-session` → `session_v2` 쿠키 저장
- 사용자 식별: `authutil.SessionUID(c)`
- 권한: Casbin(SQL adapter)
  - obj = `c.Path()`
  - act = HTTP METHOD
  - policy = DB 저장
  - model = `model.conf` (embed)

---

## 11. 라우팅 / 핸들러 가이드

- **공개 라우트**: 로그인/인덱스/프라이버시/리스트 조회
- **보호 라우트**: 저장/수정/삭제/검색/통계
- 바인딩 실패 → 400/422
- 에러 처리 → `echo.NewHTTPError`

---

## 12. 테스트 가이드

- CI 환경에는 DB가 없음 → **Mock 기반 유닛 테스트만**
- 테스트 파일명: `{파일}_test.go`
- 테스트 파일은 원본 파일과 같은 디렉토리 위치

---

## 13. 코드 생성 및 마이그레이션

- **스키마 수정**:
  `projects/{project}/migrations/*.sql`
- **SQL 수정**:
  `projects/{project}/query.sql`
- **SQLC 실행**:
  `./task.sh sqlc-generate {project}`
- **Templ 코드 생성**:
  `.templ` 수정 후 반드시 `templ generate`

---

## 14. 프로젝트 설명

- **homepage**

  - 여러 서비스를 소개하는 포털
  - CSS: TailwindCSS 사용

- **ai-study**

  - 특정 주제를 입력하면 관련 학습 주제 10개 추천

- **deario**

  - 일기를 분석해 AI 피드백 제공

- **closet**

  - 옷장 데이터 기반 AI 스타일 추천

- **sample**
  - 신기능 및 라이브러리 테스트용

---

## 15. ⭐ 최종 점검 체크리스트

- 위의 모든 규칙을 준수했는가?
- `./task.sh check` 실행해 빌드/테스트/린트 통과했는가?
- 한글 메시지/주석이 깨지지 않는가?
- 실패 시 → 원인을 분석 → 수정 → **다시 check 수행**
- check는 시간이 걸리므로 기다렸다가 결과를 확인해야 함.
