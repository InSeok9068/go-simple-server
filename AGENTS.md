# AGENTS_V2.md

> 이 문서는 **현재 go-simple-server 코드베이스의 실제 패턴**을 기준으로 정리한 작업 규칙집입니다.
> 목적은 "짧게 읽고, 실수 없이 일하기"입니다.

---

## 0) 문서 우선순위

- 이 문서는 `AGENTS_BACKUP.md`의 절대 규칙을 보완하는 실행 가이드입니다.
- 충돌 시 `AGENTS_BACKUP.md`의 강제 규칙을 우선합니다.

---

## 1) 절대 규칙

1. 모든 응답, 로그, 에러 메시지, 주석은 한국어로 작성합니다.
2. 작업 완료 전에 반드시 `./task.sh check`를 실행하고 통과해야 합니다.
3. `./task.sh check` 실패 시 원인 수정 후 재실행합니다.
4. 생성 파일은 직접 수정하지 않습니다.
5. `./task.sh`로 실행 가능한 검증/생성/포맷 작업은 반드시 Git Bash를 통해 실행합니다.

생성 파일 예시:

- `*_templ.go`
- `projects/*/db/*.go` (sqlc 생성물)

---

## 2) 프로젝트 빠른 지도

- 런타임: Go + Echo
- 템플릿: Templ
- 데이터: SQLite + Goose + SQLC
- 프론트: HTMX + Alpine.js
- 스타일: 서비스별 선택(세부 규칙은 `10) CSS 프레임워크 선택 규칙` 참조)
- 메시지 큐: goqite (경량 메시지 큐, 필요 시 선택 사용)

서비스별 성격:

- `homepage`: 소개 포털
- `deario`: 일기 + AI 피드백
- `closet`: 옷장 + 추천
- `ai-study`: 학습 주제 추천 서비스
- `sample`: 실험용, 레거시 코드 포함

레이어/패키지 경계 규칙:

1. 루트 `internal/`은 모노레포 전체 서비스가 함께 쓰는 공통 서버 레이어입니다.
2. 루트 `internal/`에는 미들웨어, 공통 에러 처리, 검증, 공통 런타임 보조 코드만 둡니다.
3. 루트 `internal/`에서 특정 서비스(`projects/{project}`)의 도메인/DB 패키지를 import하지 않습니다.
4. 루트 `pkg/`는 프로젝트 의존성이 없는 매우 작은 유틸/래퍼만 둡니다.
5. `pkg/`는 범용 함수, 경량 어댑터, 순수 헬퍼 중심으로 유지하고 서비스 정책/비즈니스 규칙은 넣지 않습니다.
6. `projects/{project}/internal/`은 해당 프로젝트 전용 서버 로직(핸들러, 서비스, 도메인 흐름)을 둡니다.
7. `projects/{project}/internal/`은 `internal/`, `pkg/`를 사용할 수 있지만, 다른 프로젝트의 `projects/{other}/internal/`은 import하지 않습니다.
8. 코드 배치 기준:
   - 여러 서비스에서 재사용 + 서버 공통 관심사: 루트 `internal/`
   - 서비스 맥락 없는 초소형 범용 유틸: 루트 `pkg/`
   - 특정 서비스 요구사항/정책/도메인 규칙: `projects/{project}/internal/`

---

## 3) 명령 실행 규칙 (Git Bash)

1. `./task.sh` 명령은 항상 Git Bash로 실행합니다.
2. PowerShell에서 실행할 때는 아래 형태를 사용합니다.

```bash
"C:/Program Files/Git/bin/bash.exe" -lc "./task.sh <command>"
```

3. 대표 명령:

- 전체 검증: `./task.sh check`
- Templ 생성: `./task.sh templ-generate`
- SQLC 생성: `./task.sh sqlc-generate {project}`
- 포맷: `./task.sh fmt`

---

## 4) UI/뷰 작성 규칙

1. 신규 화면은 기본적으로 `.templ`로 작성합니다.
2. 기존 화면도 `.templ` 기준으로 유지/수정합니다.
3. `.templ` 수정 후 반드시 Git Bash에서 `./task.sh templ-generate`를 실행합니다.
4. 인라인 `<script>`는 신규 코드에서 금지하고, JS는 `shared/static` 또는 `projects/{project}/static`로 분리합니다.
5. HTMX는 서버 상태 반영용, Alpine.js는 클라이언트 로컬 상태용으로만 사용합니다.

---

## 5) 백엔드 핸들러/응답/오류 규칙

핸들러 기본 흐름:

1. 입력 파싱(`Bind`, `FormValue`, `QueryParam`)
2. 검증(`c.Validate`, 필요 시 `validate.HTTPError`)
3. 인증 필요 시 `authutil.SessionUID(c)`
4. DB 접근(`db.GetQueries()`)
5. 응답(`Render`, `c.HTML`, `c.JSON`, `c.NoContent`, `echo.NewHTTPError`)

검증 태그 규칙:

1. DTO 검증은 반드시 `go-playground/validator` 태그(`validate:"..."`)를 사용합니다.
2. 사용자에게 노출되는 검증 문구는 DTO 필드에 `message:"..."` 태그로 정의합니다.
3. 검증 실패 응답은 `validate.HTTPError(err, &dto)`를 우선 사용해 `message` 태그 문구가 반영되도록 합니다.
4. `message`는 해당 검증 조건과 일치하는 한글 문장으로 작성합니다.

권장 패턴:

```go
var dto SomeDTO
if err := c.Bind(&dto); err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "요청 본문이 올바르지 않습니다.")
}
if err := c.Validate(&dto); err != nil {
    return validate.HTTPError(err, &dto)
}

uid, err := authutil.SessionUID(c)
if err != nil {
    return err
}

queries, err := db.GetQueries()
if err != nil {
    return err
}

return c.NoContent(http.StatusNoContent)
```

HTTP/HTMX 응답 규칙:

1. HTMX 부분 갱신: `200 + HTML`
2. 본문 없이 성공: `204 No Content`
3. 생성 성공: `201 Created`
4. 비동기 작업 접수: `202 Accepted`
5. HTMX 전체 이동: `HX-Redirect` 헤더 + `204`

주의:

- `200`에 빈 본문을 보내면 HTMX 타겟이 비워질 수 있으므로 피합니다.

오류/로깅 규칙:

1. 비즈니스 오류는 `echo.NewHTTPError(코드, "한글 메시지")`를 기본으로 사용합니다.
2. 공통 처리 가능한 오류는 `return err`로 상위 미들웨어/에러핸들러에 위임합니다.
3. 로깅은 `slog`를 사용하고, 메시지/필드도 한국어 중심으로 작성합니다.

---

## 6) DB/SQLC 규칙

1. 쿼리는 `projects/{project}/query.sql`에 작성합니다.
2. 스키마는 `projects/{project}/migrations/*.sql`에서 관리합니다.
3. 변경 후 반드시 Git Bash에서 SQLC 생성을 실행합니다.

```bash
./task.sh sqlc-generate {project}
```

4. 핸들러에서 DB 직접 SQL 작성 대신 sqlc 쿼리 메서드를 사용합니다.
5. 트랜잭션이 필요한 경우 `db.GetDB()` + `BeginTx` 패턴을 사용합니다.

---

## 7) 인증/권한 규칙

1. 로그인 세션은 `session_v2` 쿠키를 사용합니다.
2. 사용자 식별은 `authutil.SessionUID(c)`를 사용합니다.
3. 권한 라우트는 Casbin 미들웨어 그룹에서 관리합니다.
4. 인증/권한 오류는 `echo.NewHTTPError`로 명확히 반환합니다.

---

## 8) 런타임 규칙

1. 각 서비스 `cmd/main.go`에서 공통 미들웨어를 먼저 등록합니다.
2. 개발 환경은 GoVisual 래핑(`TransferEchoToGoVisualServerOnlyDev`)을 사용합니다.
3. 운영 환경은 Embed FS + gzip + rate limit + timeout(공통 미들웨어) 구성으로 실행합니다.

---

## 9) 정적 자원 규칙

1. 공통 정적 경로는 `/shared/static` -> `shared/static`입니다.
2. 프로젝트 정적 경로는 `/static` -> `projects/{service}/static`입니다.

---

## 10) CSS 프레임워크 선택 규칙

서비스별 UI 요구사항에 따라 아래 3가지 후보 중 하나를 선택해 사용합니다.
특별한 사유가 없으면 한 서비스 안에서 다중 스타일 시스템을 혼용하지 않습니다.

### 10.1 후보 1: BeerCSS (기본 선택)

1. 적용 대상: 대다수 서비스
2. 목적: 빠른 개발, 일관된 UI, Material Design 3 기반 화면 구성
3. 규칙:
   - BeerCSS 제공 클래스/컴포넌트를 우선 사용합니다.
   - 커스텀 CSS는 최소화하고 필요한 경우에만 제한적으로 사용합니다.
   - 모달/사이드바는 `data-ui` 기반 동작을 우선 검토합니다.

### 10.2 후보 2: Shoelace + TailwindCSS 조합

1. 적용 대상: Web Components 중심 설계가 필요한 서비스
2. 구성:
   - Shoelace: UI 컴포넌트(Web Components)
   - TailwindCSS: 레이아웃/간격/반응형 유틸리티
3. 규칙:
   - 컴포넌트 자체 표현은 Shoelace를 우선 사용합니다.
   - TailwindCSS는 레이아웃 보조 용도로 사용합니다.
   - 컴포넌트 역할과 레이아웃 역할을 분리해 유지보수성을 확보합니다.

### 10.3 후보 3: TailwindCSS 단독

1. 적용 대상: 고자유도 커스텀 디자인이 필요한 서비스
2. 목적: 브랜드/화면 스타일을 세밀하게 제어
3. 규칙:
   - 색상/간격/타이포 토큰을 일관되게 관리합니다.
   - 유틸리티 클래스 반복 패턴은 재사용 가능한 템플릿 단위로 정리합니다.
   - 모바일 우선 반응형 설계를 기본으로 작성합니다.

### 10.4 선택 방식

1. 서비스별 CSS 프레임워크 선택은 사용자(요청자) 지정에 따릅니다.
2. 에이전트는 지정된 후보를 임의 변경하지 않습니다.
3. 지정이 없을 때만 BeerCSS를 기본값으로 제안할 수 있습니다.

### 10.5 기록 위치(필수)

1. 서비스별 선택 결과와 선택 이유를 `projects/{project}/README.md`에 기록합니다.
2. 최소 기록 항목:
   - 선택한 CSS 프레임워크 후보(1/2/3)
   - 선택 이유(요구사항 기준)
   - 예외 규칙(혼용 여부, 제한 조건)

---

## 11) 작업 절차 체크리스트

1. 수정 범위를 먼저 고정합니다(요청 범위 외 리팩토링 금지).
2. 코드 수정 후 생성 작업을 수행합니다.

- `.templ` 변경 시: `./task.sh templ-generate`
- `query.sql`/`migrations` 변경 시: `./task.sh sqlc-generate {project}`

3. 마지막에 반드시 Git Bash에서 아래를 실행합니다.

```bash
./task.sh check
```

4. 실패하면 수정 후 3번을 반복합니다.
5. 결과 보고 시 아래를 함께 남깁니다.

- 변경 파일
- 생성 명령 실행 여부
- `./task.sh check` 통과 여부

참고:

- 명령 예시는 `13) 빠른 명령 모음`을 사용합니다.

---

## 12) Agent 금지 사항

1. 생성 파일 직접 수정 금지
2. 요청 없는 대규모 구조 개편 금지
3. HTMX 대상에 의도치 않은 빈 응답(`200 + empty body`) 금지
4. 검증 명령 생략 후 완료 선언 금지

---

## 13) 빠른 명령 모음

```bash
# 공통 실행 방식 (Git Bash 필수)
"C:/Program Files/Git/bin/bash.exe" -lc "./task.sh <command>"

# 전체 검증 (필수)
./task.sh check

# Templ 생성
./task.sh templ-generate

# SQLC 생성
./task.sh sqlc-generate deario
./task.sh sqlc-generate closet

# 포맷
./task.sh fmt
```
