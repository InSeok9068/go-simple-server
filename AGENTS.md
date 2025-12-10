# AGENTS.md

> 이 문서는 **Go + Echo + Templ 기반 Monorepo**에서
> 사람과 LLM(코드 어시스턴트) 모두가 일관된 코드를 작성하기 위한 **단일 규칙집**이다.

---

> [!IMPORTANT] > **CRITICAL RULE**: 모든 작업을 마치기 전, 반드시 `./task.sh check` 명령어를 실행하고 **모든 검사를 통과(Pass)** 해야 한다.
> LLM(Codex)은 이 명령어가 실패하면 절대 작업을 완료했다고 판단해서는 안 되며, 반드시 원인을 수정하고 재시도해야 한다.
> 이 규칙은 절대 무시할 수 없다.

---

## 1. 공통 규칙

1. 모든 답변·로그·주석·UI 텍스트는 **한글**로 작성한다.
2. **불필요한 추상화나 복잡성보다 직관적이고 명확한 코드**를 우선한다.
3. 공통 요소가 충분히 반복되고, 효과가 검증되었을 때만 **리팩토링**한다.
4. 모든 화면과 UI는 **모바일 우선(Mobile-first)**으로 설계하고, 이후 데스크톱을 보완한다.

> LLM이 코드를 생성할 때도 위 네 가지 원칙을 항상 우선 적용한다.

---

## 2. 프로젝트 구성 개요

- **언어**: Go
- **서버 프레임워크**: Echo
- **템플릿 엔진**: Templ (**Gomponents는 사용하지 않음**)
- **데이터베이스**: SQLite
- **DB 마이그레이션**: Goose
- **SQL 코드 생성기**: SQLC
- **프론트엔드 라이브러리**: HTMX, Alpine.js
- **CSS 프레임워크**
  - 기본: BeerCSS
  - 예외: `homepage` 프로젝트만 TailwindCSS 사용
- **개발 서버**: GoVisual
- **소스 구조**: Monorepo

---

## 3. 프론트엔드 라이브러리 역할 정의

### 3.1 HTMX

- **역할**
  - AJAX 요청으로 서버에서 **HTML fragment**를 받아 특정 영역만 동적으로 갱신.
- **목적**
  - 전체 페이지 리로드 없이도 **SPA와 유사한 사용자 경험** 제공.

### 3.2 Alpine.js

- **역할**
  - 클라이언트 측에서 경량 상태 관리 및 UI 인터랙션 처리.
- **목적**
  - 모달, 드롭다운, 탭 등 **즉각적인 UI 동작** 구현.

### 3.3 HTMX × Alpine.js 적용 패턴

- 여러 프로젝트에서 재사용되는 HTML 조각은
  `shared/views` 또는 `shared/static/js`로 **공통화**한다.
- **Alpine.js 사용 범위**
  - 반드시 **클라이언트에서만 결정 가능한 미세 상호작용**에 사용한다.
  - 예: 토글, 모달 열기/닫기, 단순 상태 토글 등.
- **서버 데이터 기반 UI 상태**
  - 항상 **HTMX 응답(서버 렌더링 HTML)** 으로 갱신한다.
  - 클라이언트에서 임의로 서버 데이터 상태를 가정하지 않는다.

### 3.4 Templ

- **역할**
  - HTML과 유사한 문법으로 UI를 작성하면, 정적 Go 코드로 컴파일하는 템플릿 엔진.
- **목적**
  - `.templ`에서 작성한 마크업을 `templ generate`로 미리 Go 코드로 변환하여
    런타임 파싱 비용 없이 **빠르고 타입 안정적인 렌더링**을 제공한다.
- **공통 컴포넌트 위치**
  - `simple-server/shared/views`

> LLM은 HTML 템플릿이 필요할 때 **반드시 Templ 문법**을 사용하고,
> 기존 규칙에 따라 공통 컴포넌트는 `shared/views`에 배치한다.

---

## 4. CSS 라이브러리 사용 정의

### 4.1 BeerCSS (기본)

- **역할**
  - HTML 안에 BeerCSS 클래스만을 사용해 UI를 구성한다.
  - **Material Design 3**를 계승한 프레임워크이다.
- **동작**
  - 모달·사이드바 등의 UI 동작은 `data-ui="#id"`만으로 제어 가능 → 이 경우 Alpine.js 불필요.
- **참고**
  - 공식 문서: <https://beercss.com/>
  - Material Design 3 가이드라인: <https://m3.material.io/>
  - 사용 가능 태그 및 클래스 목록: `.doc/css/beercss/SUMMARY.md`
  - 올바른 BeerCSS 사용 가이드: `.doc/css/beercss/GUIDE.md`
- ❗ **규칙**
  - **커스텀 CSS 금지**
  - BeerCSS에서 제공하지 않는 클래스 사용 금지
  - **Material Design 3 가이드라인**을 준수하여 작업한다.

### 4.2 TailwindCSS (homepage 전용)

- 사용 대상: **`homepage` 프로젝트에만** TailwindCSS 사용 가능.
- 설치/세팅:
  ```bash
  ./task.sh install-tailwind
  ```

---

## 5. 폴더 구조 규칙

```text
cmd/{project}/main.go        # 루트 CLI 프로젝트 실행파일
projects/{project}/          # 프로젝트 소스
├─ cmd/                      # 프로젝트 실행 파일
├─ internal/                 # 프로젝트 내부 로직
├─ views/                    # Templ 템플릿
├─ static/                   # 정적 파일 (JS, 이미지 등)
├─ migrations/               # Goose 마이그레이션
└─ query.sql                 # SQLC 쿼리 정의
internal/                    # 외부 노출되지 않는 공용 서버 코드
└─ validate/                 # 검증 패키지
shared/                      # 공통 프론트엔드 자원
├─ views/                    # 공통 Templ 컴포넌트
└─ static/                   # 공통 JS, CSS
pkg/util/                    # 순수 유틸 함수 (비즈니스 로직 없음)
```

> 새로운 프로젝트를 추가할 때는 반드시 위 구조를 그대로 따른다.

---

## 6. 코딩 스타일 가이드

### 6.1 Handler 패턴 (Echo)

Handler의 기본 흐름은 **항상 같은 패턴**을 따른다.

1. **바인딩 & 검증**

   ```go
   var dto SomeDTO
   if err := c.Bind(&dto); err != nil {
       return echo.NewHTTPError(http.StatusBadRequest, "요청 형식이 올바르지 않습니다.")
   }

   if err := c.Validate(&dto); err != nil {
       // validate.HTTPError 또는 echo.NewHTTPError 사용
       return validate.HTTPError(err)
   }
   ```

   - 바인딩 실패 → 400 또는 422
   - 검증 실패 → `validate.HTTPError` 또는 `echo.NewHTTPError` 사용
   - 에러 메시지는 항상 **한글로 명확하게** 작성

2. **DB 접근**

   ```go
   queries, err := db.GetQueries()
   if err != nil {
       return echo.NewHTTPError(http.StatusInternalServerError, "데이터베이스 연결에 실패했습니다.")
   }

   ctx := c.Request().Context()
   result, err := queries.Method(ctx, params)
   if err != nil {
       // 상황에 맞게 4xx 또는 5xx 반환
       // 예: 데이터 없음 → 404, 권한 문제 → 403, 서버 문제 → 500
   }
   ```

3. **응답**
   - 정상 응답 (Templ 렌더링):

     ```go
     return views.Component(...).Render(ctx, c.Response().Writer)
     ```

   - 오류 응답:

     ```go
     return echo.NewHTTPError(code, "에러 메시지")
     ```

### 6.2 에러 핸들링 & 로깅

- 에러 처리는 Go 표준 방식 사용:

  ```go
  if err != nil {
      // 로깅 후 적절한 HTTPError 반환
  }
  ```

- 에러 응답: 항상 `echo.NewHTTPError` 사용
- 로깅: `slog` 사용
  - 레벨: `DEBUG`, `INFO`, `WARN`, `ERROR`
  - 로그 메시지와 필드는 **한글로 명확하게** 작성

### 6.3 Templ 작성 가이드

- `.templ` 파일은 **HTML과 거의 동일한 구조**로 작성한다.
- 기본 전략은 **서버 주도 렌더링(SSR)** 이다.

**동적 기능 규칙**

- 데이터 기반 UI 동작 → **HTMX** 사용
  (주요 속성: `hx-get`, `hx-post`, `hx-target`, `hx-swap` 등)
- 클라이언트 로컬 상태 위주의 UI 동작 → **Alpine.js** 사용
  (모달, 탭, 드롭다운, 단순 토글 등)
- JS 코드는 항상 **별도 `.js` 파일**로 분리한다.
  (Templ 안에 인라인 `<script>`를 두지 않는다.)

**중요**

- `.templ` 파일을 수정하면 **반드시** 다음을 실행한다.

  ```bash
  templ generate
  ```

  → Go 코드 재생성 후 커밋한다.

---

## 7. HTTP / HTMX 응답 규칙

HTMX와 함께 사용할 때의 응답 규칙을 명확히 한다.

### 7.1 상태 코드 사용 원칙

- `204 No Content`
  - 저장/수정/삭제 성공 등, **본문 없이** UI를 유지하고 싶은 경우
  - HTMX 타겟 영역이 비워지지 않아야 할 때 사용
- `201 Created`
  - 새 리소스가 생성된 경우
- `202 Accepted`
  - 비동기 작업이 접수되었으나 아직 완료되지 않은 경우
- `200 OK + HTML`
  - HTMX로 **부분 갱신**이 필요한 경우
- 오류 응답
  - 적절한 4xx/5xx 코드 + 한글 오류 메시지

> ⚠ 주의: `200 OK` + **빈 본문**을 보내면 HTMX가 대상 영역을 비워버린다.
> 의도적으로 본문이 없을 때는 **반드시 `204`를 사용**한다.

### 7.2 리다이렉션 (HTMX)

- HTMX 요청 후 전체 페이지 이동이 필요할 때는
  `HX-Redirect` 헤더를 사용한다.

---

## 8. 런타임 모드

### 8.1 개발 모드

- 개발 서버: **GoVisual** 사용
- 정적 파일: 로컬 디스크에서 직접 제공
- 대표 기능:
  - `TransferEchoToGoVisualServerOnlyDev`

### 8.2 운영 모드

- Echo 단독 실행
- 정적 파일: `embed.FS` 사용
- Gzip 압축 활성화
- 요청 타임아웃: **1분**
- 요청 바디 최대 크기: **5MB**
- Rate Limit: **초당 20회**

---

## 9. 정적 / 임베드 자원 규칙

- 공통 정적 자원
  - URL: `/shared/static`
  - 실제 경로: `shared/static/*`
- 프로젝트별 정적 자원
  - URL: `/static`
  - 실제 경로: `projects/{name}/static/*`

- PWA 관련 파일 매핑
  - `/manifest.json`
  - `/firebase-messaging-sw.js`

- 개발 환경
  - `os.DirFS`로 파일 시스템에서 직접 읽기
- 운영 환경
  - `resources.EmbeddedFiles` (embed) 사용

---

## 10. 인증 및 권한

### 10.1 인증

- Firebase ID 토큰을 받아 `/create-session`에서 세션 생성
- 세션은 `session_v2` 쿠키에 저장
- 요청에서 사용자 식별 시:

  ```go
  uid, err := authutil.SessionUID(c)
  ```

### 10.2 권한 (Authorization)

- 권한 관리는 Casbin(SQL adapter) 사용
  - `obj` = `c.Path()` (요청 경로)
  - `act` = HTTP Method (`GET`, `POST` 등)
  - policy = DB에 저장
  - model = `model.conf` (embed)

---

## 11. 라우팅 / 핸들러 가이드

- **공개 라우트**
  - 로그인, 인덱스, 프라이버시, 리스트 조회 등
- **보호 라우트**
  - 저장, 수정, 삭제, 검색, 통계 등

- 바인딩 실패
  - 400 또는 422 반환
- 모든 에러 처리
  - `echo.NewHTTPError`로 통일

---

## 12. 테스트 가이드

- CI 환경에는 **실제 DB가 없음**
  - → **Mock 기반 유닛 테스트만** 작성한다.
- 테스트 파일 네이밍
  - `{파일명}_test.go`
- 테스트 파일 위치
  - 원본 파일과 **같은 디렉토리**에 둔다.

---

## 13. 코드 생성 및 마이그레이션

- **스키마 수정**
  - 경로: `projects/{project}/migrations/*.sql`
- **SQL 수정**
  - 경로: `projects/{project}/query.sql`
- **SQLC 코드 생성**
  ```bash
  ./task.sh sqlc-generate {project}
  ```
- **Templ 코드 생성**
  - `.templ` 수정 후 반드시 아래 명령 실행:

    ```bash
    templ generate
    ```

---

## 14. 프로젝트 설명

- **homepage**
  - 여러 서비스를 소개하는 포털 역할
  - CSS: **TailwindCSS** 사용

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

아래 항목을 모두 만족해야 머지/배포할 수 있다.

1. 이 문서의 모든 규칙을 준수했는가?
2. `./task.sh check` 다음 명령이 **성공적으로 통과**했는가?
   - 빌드
   - 테스트
   - 린트
3. 한글 메시지/주석/로그가 깨지지 않는가?
4. `./task.sh check` 실패 시
   - 원인을 분석한다.
   - 수정한다.
   - 다시 `./task.sh check`를 수행한다.

> ⚠ 주의: `./task.sh check`는 시간이 걸리므로 기다렸다가 결과를 확인한다.
> **다시 한 번 강조한다. `./task.sh check`를 통과하지 못하면 작업을 완료할 수 없다.**

---

> 이 문서가 변경될 경우, 관련 레포의 코드/설정도 함께 검토하고,
> 변경 이유를 커밋 메시지에 **명확한 한글로 기록**한다.
