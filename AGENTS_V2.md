# AGENTS_V2.md

> 이 문서는 **현재 go-simple-server 코드베이스의 실제 패턴**을 기준으로 정리한 작업 규칙집입니다.
> 목적은 "짧게 읽고, 실수 없이 일하기"입니다.

---

## 0) 문서 우선순위

- `AGENTS.md`의 절대 규칙(특히 `./task.sh check` 통과 전 완료 금지)은 그대로 유지합니다.
- 이 문서는 기존 규칙을 대체하기보다, **현행 코드 패턴 중심의 실행 가이드**를 제공합니다.

---

## 1) 절대 규칙

1. 모든 응답, 로그, 에러 메시지, 주석은 한국어로 작성합니다.
2. 작업 완료 전에 반드시 `./task.sh check`를 실행하고 통과해야 합니다.
3. `./task.sh check` 실패 시 원인 수정 후 재실행합니다.
4. 생성 파일은 직접 수정하지 않습니다.

- `*_templ.go`
- `projects/*/db/*.go` (sqlc 생성물)

---

## 2) 프로젝트 빠른 지도

- 런타임: Go + Echo
- 템플릿: Templ
- 데이터: SQLite + Goose + SQLC
- 프론트: HTMX + Alpine.js
- 스타일:
- 기본: BeerCSS
- 예외: `homepage`는 Tailwind + Shoelace

서비스별 성격:

- `homepage`: 소개 포털, Tailwind 중심
- `deario`: 일기 + AI 피드백 (BeerCSS + HTMX)
- `closet`: 옷장 + 추천 (BeerCSS + HTMX)
- `ai-study`: 학습 주제 추천 서비스
- `sample`: 실험용, 레거시 코드 포함

---

## 3) UI/뷰 작성 규칙

1. 신규 화면은 기본적으로 `.templ`로 작성합니다.
2. 기존 화면도 동일하게 `.templ` 기준으로 유지/수정합니다.
3. `.templ` 수정 후 반드시 `templ generate`를 실행합니다.
4. 인라인 `<script>`는 신규 코드에서 금지하고, JS는 `shared/static` 또는 `projects/{project}/static`로 분리합니다.
5. HTMX는 서버 상태 반영용, Alpine.js는 클라이언트 로컬 상태용으로만 사용합니다.

---

## 4) 핸들러 표준 패턴 (Echo)

핵심 흐름:

1. 입력 파싱(`Bind`, `FormValue`, `QueryParam`)
2. 검증(`c.Validate`, 필요 시 `validate.HTTPError`)
3. 인증 필요 시 `authutil.SessionUID(c)`
4. DB 접근(`db.GetQueries()`)
5. 응답(`Render`, `c.HTML`, `c.JSON`, `c.NoContent`, `echo.NewHTTPError`)

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

---

## 5) HTTP/HTMX 응답 규칙

1. HTMX 부분 갱신: `200 + HTML`
2. 본문 없이 성공: `204 No Content`
3. 생성 성공: `201 Created`
4. 비동기 작업 접수: `202 Accepted`
5. HTMX 전체 이동: `HX-Redirect` 헤더 + `204`

주의:

- `200`에 빈 본문을 보내면 HTMX 타겟이 비워질 수 있으므로 피합니다.

---

## 6) DB/SQLC 규칙

1. 쿼리는 `projects/{project}/query.sql`에 작성합니다.
2. 스키마는 `projects/{project}/migrations/*.sql`에서 관리합니다.
3. 변경 후 반드시 SQLC 생성:

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

## 8) 런타임/정적 자원 규칙

1. 각 서비스 `cmd/main.go`에서 공통 미들웨어를 먼저 등록합니다.
2. 정적 자원 경로:

- 공통: `/shared/static` -> `shared/static`
- 프로젝트: `/static` -> `projects/{service}/static`

3. 개발: GoVisual 래핑(`TransferEchoToGoVisualServerOnlyDev`)
4. 운영: Embed FS + gzip + rate limit + timeout(공통 미들웨어)

---

## 9) 로깅/에러 처리 규칙

1. 비즈니스 오류는 `echo.NewHTTPError(코드, "한글 메시지")`를 기본으로 사용합니다.
2. 공통 처리 가능한 오류는 `return err`로 상위 미들웨어/에러핸들러에 위임합니다.
3. 로깅은 `slog`를 사용하고, 메시지/필드도 한국어 중심으로 작성합니다.

---

## 10) CSS 규칙 (현 코드 기준)

1. `homepage`는 Tailwind + Shoelace 조합을 유지합니다.
2. `deario`, `closet`은 BeerCSS 기반 + 프로젝트별 `static/style.css` 보강 패턴을 유지합니다.
3. 스타일 변경은 "기능 전달" 목적을 우선하고, 단순 미관 리팩토링은 요청 시에만 진행합니다.
4. 모바일 우선 레이아웃을 기본으로 작성합니다.

---

## 11) 테스트 규칙

1. 테스트 파일은 `{name}_test.go` 형식, 원본 파일과 같은 디렉터리에 둡니다.
2. 외부 인프라 의존 테스트보다 유닛 테스트를 우선합니다.
3. 최소한 아래는 검증합니다.

- 입력 검증/에러 분기
- 핵심 유틸 함수
- 미들웨어/핸들러의 실패 경로

---

## 12) 작업 절차 체크리스트

1. 수정 범위를 먼저 고정합니다(요청 범위 외 리팩토링 금지).
2. 코드 수정 후 생성 작업을 수행합니다.

- `.templ` 변경 시: `templ generate`
- `query.sql`/`migrations` 변경 시: `./task.sh sqlc-generate {project}`

3. 마지막에 반드시 실행합니다.

```bash
./task.sh check
```

4. 실패하면 수정 후 3번을 반복합니다.
5. 결과 보고 시 아래를 함께 남깁니다.

- 변경 파일
- 생성 명령 실행 여부
- `./task.sh check` 통과 여부

---

## 13) Agent 금지 사항

1. 생성 파일 직접 수정 금지
2. 요청 없는 대규모 구조 개편 금지
3. HTMX 대상에 의도치 않은 빈 응답(`200 + empty body`) 금지
4. 검증 명령 생략 후 완료 선언 금지

---

## 14) 빠른 명령 모음

```bash
# 전체 검증 (필수)
./task.sh check

# Templ 생성
./task.sh templ-generate
# 또는
templ generate

# SQLC 생성
./task.sh sqlc-generate deario
./task.sh sqlc-generate closet

# 포맷
./task.sh fmt
```
