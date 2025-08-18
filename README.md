# Go 언어 기반 풀스택 서버

**Host** : https://toy-project.n-e.kr

| 구성                  | 사용 기술/도구                                                |
| --------------------- | ------------------------------------------------------------- |
| **언어**              | Go                                                            |
| **백엔드 프레임워크** | Echo                                                          |
| **프론트엔드 구성**   | HTMX + Alpine.js                                              |
| **템플릿 엔진**       | Gomponents                                                    |
| **CSS 프레임워크**    | 1. Beer CSS <br> 2. Pico CSS + Tailwind CSS <br> 3. Bulma CSS |
| **데이터베이스**      | SQLite → PostgreSQL                                           |
| **DB 관리 도구**      | Adminer                                                       |
| **DB 마이그레이션**   | Goose                                                         |
| **DB 복제/복구**      | Litestream (SQLite)                                           |
| **인증**              | Firebase (With Cookie)                                        |
| **객체 저장소**       | GCP Storage                                                   |
| **메시지 큐**         | goqite                                                        |
| **성능/로깅 도구**    | /debug/vars (Go 표준), trace_id (OTEL 연동)                   |
| **로깅 대시보드 UI**  | PocketBase Admin                                              |
| **모바일 대응**       | PWA → Capacitor                                               |
| **웹 서버**           | Caddy                                                         |

## 폴더 구조

- `cmd/{프로젝트명}/main.go`: 각 프로젝트의 서버 실행 파일
- `projects/{프로젝트명}/`: 프로젝트별 소스 코드
  - `handlers`: HTTP 요청을 처리하는 핸들러
  - `views`: Gomponents로 작성된 HTML 뷰 컴포넌트
  - `static`: CSS, JavaScript 등 정적 파일
  - `migrations`: Goose 기반 데이터베이스 마이그레이션
  - `query.sql`: SQLC가 사용하는 쿼리 정의
- `internal`: 여러 프로젝트에서 공유하는 서버 공통 패키지
- `shared`: 공통 뷰 컴포넌트와 정적 자산
- `pkg`: 외부 의존성이 없는 유틸리티 함수

## 프로젝트 설명

- **homepage**: 여러 서비스의 소개와 진입점을 제공하는 포털 (TailwindCSS 사용)
- **ai-study**: 입력한 주제와 관련된 학습 주제를 AI가 제안
- **deario**: 일기를 작성하면 AI가 피드백을 제공
- **sample**: 새로운 기능이나 라이브러리를 실험하는 샘플 프로젝트

### 추가 구성 고려

- Beer CSS : CSS 프레임워크 (Material 모바일 우선)
- Pico CSS + Tailwind CSS : CSS 프레임워크 (커스터마이징 용이)
- Bulma CSS : CSS 프레임워크 (심플한 관리자 UI)
- Tabler : CSS 프레임워크 (복잡한 관리자 UI) - PC 환경
- Shoelace : CSS 프레임워크 (바닐라 웹 컴포넌트)
- ~~surreal.js : [surreal.js](https://cdn.jsdelivr.net/gh/gnat/surreal@main/surreal.js)~~
- DataStar : SSE 기반 서버주도 프레임워크
- Unpoly : Htmx Like 프레임워크

---

## 실행

### 서비스 실행

#### 윈도우 (개발 환경)

```shell
air
```

#### 리눅스

```shell
chmod +x main
./main
```

```shell
sudo systemctl start main.service
```

#### 윈도우 (로그 어드민)

```shell
./pocketbase.exe serve --dir ./shared/log
```

#### 리눅스 (로그 어드민)

```shell
./pocketbase serve --dir ./srv/log
```

---

## 종속성 최신화

### 프로젝트 종속성 최신화

```shell
go get cloud.google.com/go/storage
go get firebase.google.com/go/v4
go get github.com/AlecAivazis/survey/v2
go get github.com/Blank-Xu/sql-adapter
go get github.com/casbin/casbin/v2
go get github.com/crazy-max/echo-ipfilter
go get github.com/doganarif/govisual
go get github.com/eduardolat/gomponents-lucide
go get github.com/glsubri/gomponents-alpine
go get github.com/go-rod/rod
go get github.com/gorilla/sessions
go get github.com/joho/godotenv
go get github.com/labstack/echo-contrib
go get github.com/labstack/echo/v4
go get github.com/lmittmann/tint
go get github.com/pressly/goose/v3
go get github.com/qustavo/sqlhooks/v2
go get github.com/robfig/cron/v3
go get github.com/willoma/bulma-gomponents
go get github.com/willoma/gomplements
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/trace
go get golang.org/x/time
go get google.golang.org/api
go get google.golang.org/genai
go get maragu.dev/gomponents
go get maragu.dev/gomponents-htmx
go get maragu.dev/goqite
go get modernc.org/sqlite
```

### 개발 도구 종속성 최신화

```shell
#go install github.com/a-h/templ/cmd/templ@latest
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

### 불필요한 종속성 제거

```shell
go mod tidy
```

### JS, CSS 파일 벤더링

#### CSS

[bulma.min.css](https://cdn.jsdelivr.net/npm/bulma/css/bulma.min.css) </br>
[beer.min.css](https://cdn.jsdelivr.net/npm/beercss/dist/cdn/beer.min.css) </br>
[open-props.min.css](https://cdn.jsdelivr.net/npm/open-props/open-props.min.css) </br>
[pico.classless.min.css](https://cdn.jsdelivr.net/npm/@picocss/pico/css/pico.classless.min.css)
[flatpickr.min.css](https://cdn.jsdelivr.net/npm/flatpickr@4.6.13/dist/flatpickr.min.css)
[dark.css](https://cdn.jsdelivr.net/npm/flatpickr@4.6.13/dist/themes/dark.css)

#### JS

[cdn.min.js](https://cdn.jsdelivr.net/npm/alpinejs/dist/cdn.min.js) </br>
[htmx.min.js](https://cdn.jsdelivr.net/npm/htmx.org/dist/htmx.min.js) </br>
[beer.min.js](https://cdn.jsdelivr.net/npm/beercss/dist/cdn/beer.min.js) </br>
[material-dynamic-colors.min.js](https://cdn.jsdelivr.net/npm/material-dynamic-colors/dist/cdn/material-dynamic-colors.min.js)
[flatpickr.min.js](https://cdn.jsdelivr.net/npm/flatpickr@4.6.13/dist/flatpickr.min.js)
[ko.js](https://cdn.jsdelivr.net/npm/flatpickr@4.6.13/dist/l10n/ko.js)
[marked.umd.js](https://cdn.jsdelivr.net/npm/marked/lib/marked.umd.js)

---

## GCC 활성화

Window : [tdm-gcc](https://jmeubank.github.io/tdm-gcc/)

```shell
go env -w CGO_ENABLED=1
```

---

## 오류 검증

```shell
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6
```

```shell
golangci-lint run ./...
```

---

## 보안 이슈 검증

```shell
govulncheck ./...
```

---

## 참고 글

https://ntorga.com/full-stack-go-app-with-htmx-and-alpinejs/

## Tailwind CSS 자동완성 (Gomponents)

```json
{
  "includeLanguages": {
    ...
    "go": "html"
  },
  "experimental": {
    ...
    "classRegex": [
      [
        "Class(?:es)?[({]([^)}]*)[)}]",
        "[\"`]([^\"`]*)[\"`]"
      ]
    ]
  }
}
```

---

## 이미 Merge된 브랜치 제거

```bash
for branch in $(git for-each-ref refs/remotes/origin/ --format='%(refname:short)' \
  | grep -E '^origin/(feature/|codox/)'); do

  if git merge-base --is-ancestor "$branch" origin/main; then
    echo "🗑 삭제: $branch"
    git push origin --delete "${branch#origin/}"
  fi
done
```

---

## Gemini CLI 설치 및 자동화 커밋메시지 도구 설치

```shell
npm install -g @google/gemini-cli

gemini

npm install -g gemini-commit-assistant

aic
```

---

## 모바일, 데스크톱 성능 분석

[MoCheck](https://mocheck.netlify.app/ko)

---

## PWA 앱 출시

1. npm install -g @bubblewrap/cli
2. bubblewrap init --manifest https://deario.toy-project.n-e.kr/manifest.json
3. assetlinks.json 추가
4. bubblewrap build

---

## Firebase Cloud Function

### 콘솔 설치

```shell
npm install -g firebase-tools
```

### 코드 배포

```shell
firebase deploy
```
