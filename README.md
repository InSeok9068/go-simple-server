# Go 언어 기반 풀스택 서버

> **주의** : AI Agent는 해당 파일을 읽지않고 지나가야합니다.

**Host** : https://toy-project.n-e.kr

| 구성                  | 사용 기술/도구                                                               |
| --------------------- | ---------------------------------------------------------------------------- |
| **언어**              | Go                                                                           |
| **백엔드 프레임워크** | Echo                                                                         |
| **프론트엔드 구성**   | HTMX + Alpine.js                                                             |
| **템플릿 엔진**       | Templ                                                                         |
| **CSS 프레임워크**    | 1순위 : Beer CSS <br> 2순위 : Pico CSS + Tailwind CSS <br> 3순위 : Bulma CSS |
| **데이터베이스**      | SQLite → PostgreSQL                                                          |
| **DB 관리 도구**      | Adminer                                                                      |
| **DB 마이그레이션**   | Goose                                                                        |
| **DB 복제/복구**      | Litestream (SQLite)                                                          |
| **인증**              | Firebase (With Cookie)                                                       |
| **객체 저장소**       | GCP Storage                                                                  |
| **메시지 큐**         | goqite                                                                       |
| **성능/로깅 도구**    | /debug/vars (Go 표준), trace_id (OTEL 연동)                                  |
| **로깅 대시보드 UI**  | PocketBase Admin                                                             |
| **모바일 대응**       | PWA → Capacitor                                                              |
| **웹 서버**           | Caddy                                                                        |

## 폴더 구조

- `cmd/{프로젝트명}/main.go`: 루트 단위의 실행 파일
- `projects/{프로젝트명}/`: 프로젝트별 소스 코드
  - `cmd`: 프로젝트 실행 파일
  - `internal`: 프로젝트 내부 로직
  - `views`: Templ로 작성된 HTML 뷰 컴포넌트
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
- **closet**: 옷장 관리
- **sample**: 새로운 기능이나 라이브러리를 실험하는 샘플 프로젝트

### CSS 구성 고려

| 라이브러리              | 강점                                           | 합리적인 사용 케이스                       |
| ----------------------- | ---------------------------------------------- | ------------------------------------------ |
| **Beer CSS**            | 모바일 퍼스트, 간단한 Material UI, 매우 가벼움 | 퍼블릭 웹, 모바일 중심 서비스, 빠른 개발   |
| **Shoelace + Tailwind** | 바닐라 Web Components, 접근성 최강             | Modal/Drawer 등 고급 UI가 필요한 특정 영역 |
| **Pico + Tailwind**     | 기본은 깔끔, 디테일은 강력한 커스터마이징      | UI 디테일 잡기 필요한 프로젝트             |
| **Bulma**               | 단순하고 탄탄한 관리자 UI                      | 기본 백오피스, 운영툴                      |
| **Tabler**              | 대시보드/관리자용 강력한 컴포넌트              | 복잡한 PC 기반 관리자 화면                 |
| **TemplUI**             | Templ 기반 컴포넌트화, 고생산성                | Templ로 만든 프로젝트 전용 UI              |

### JS 라이브러리 구성 고려

| 라이브러리         | 강점                         | 합리적인 사용 케이스 |
| ------------------ | ---------------------------- | -------------------- |
| **DataStar**       | SSE 기반 서버주도 프레임워크 |
| **Unpoly**         | Htmx Like 프레임워크         |                      |
| ~~**surreal.js**~~ |                              |                      |

### 사용중인 CSS/JS 라이브러리

| 종류 | 라이브러리                  | 역할                                                        |
| :--- | :-------------------------- | :---------------------------------------------------------- |
| 코어 | **htmx.org**                | 서버 주도 UI 업데이트 (AJAX, Websockets, SSE)               |
| 코어 | **alpinejs**                | 클라이언트 측 경량 상태 관리 및 UI 상호작용                 |
| 코어 | **@alpinejs/persist**       | Alpine 상태를 로컬 스토리지에 자동 저장                     |
| 코어 | **@alpinejs/morph**         | DOM 변경 시 부드러운 전환(Morphing) 효과                    |
| 코어 | **htmx-ext-alpine-morph**   | HTMX와 Alpine.js Morphing 기능 연동                         |
| UI   | **beercss**                 | 메인 프레임워크. Material Design 3 기반 UI                  |
| UI   | **@picocss/pico**           | Classless CSS (최소한의 기본 스타일링)                      |
| UI   | **bulma**                   | 유틸리티 및 컴포넌트 기반 CSS 프레임워크                    |
| UI   | **open-props**              | CSS 변수 모음 (색상, 그림자, 애니메이션 등)                 |
| UI   | **material-dynamic-colors** | Material Design 3 동적 색상 테마 생성                       |
| 유틸 | **flatpickr**               | 경량 날짜 및 시간 선택기                                    |
| 유틸 | **chart.js**                | HTML5 Canvas 기반 데이터 시각화 차트                        |
| 유틸 | **marked**                  | 마크다운 텍스트를 HTML로 변환                               |
| 유틸 | **hammerjs**                | 멀티 터치 제스처 (스와이프, 핀치 등) 이벤트 처리            |
| 유틸 | **list.js**                 | 테이블, 리스트 기반 정렬, 검색                              |
| 유틸 | **hotkeys.js**              | 키보드 단축키 처리                                          |
| 기타 | **TailwindCSS**             | `homepage` 프로젝트 전용 스타일링 (별도 빌드 프로세스 사용) |

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
go get github.com/casbin/casbin/v3
go get github.com/crazy-max/echo-ipfilter
go get github.com/doganarif/govisual
go get github.com/go-rod/rod
go get github.com/gorilla/sessions
go get github.com/joho/godotenv
go get github.com/labstack/echo-contrib
go get github.com/labstack/echo/v4
go get github.com/lmittmann/tint
go get github.com/pressly/goose/v3
go get github.com/qustavo/sqlhooks/v2
go get github.com/robfig/cron/v3
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/trace
go get golang.org/x/time
go get google.golang.org/api
go get google.golang.org/genai
go get maragu.dev/goqite
go get modernc.org/sqlite
```

### 개발 도구 종속성 최신화

```shell
go install github.com/air-verse/air@latest
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/dexter2389/go-tailwind-sorter@latest
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

## Tailwind CSS 자동완성 (Templ)

```json
{
  "includeLanguages": {
    "templ": "html"
  }
}
```

---

## 브랜치 정리

### 이미 Merge된 브랜치 정리

```bash
for branch in $(git for-each-ref refs/remotes/origin/ --format='%(refname:short)' \
  | grep -E '^origin/(feature/|codox/)'); do

  if git merge-base --is-ancestor "$branch" origin/main; then
    echo "🗑 삭제: $branch"
    git push origin --delete "${branch#origin/}"
  fi
done
```

### 규칙에 의한 브랜치 정리

**로컬**

```bash
git branch | grep 'feature/' | xargs git branch -D
git branch | grep 'codex/' | xargs git branch -D
```

**원격**

```bash
git branch -r | grep 'origin/feature/' | sed 's/origin\///' | xargs -I {} git push origin --delete {}
git branch -r | grep 'origin/codex/' | sed 's/origin\///' | xargs -I {} git push origin --delete {}
```

### 특정 브랜드 제거

```bash
git push origin --delete og70vp-codex/refactor-initcasbin-to-manage-db-connection
git push origin --delete fu4e2s-codex
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

### 코드 전체 배포

```shell
firebase deploy
```

### 함수 명령어

#### 함수 목록 보기

```shell
firebase functions:list
```

#### 함수 배포

```shell
# 개별
firebase deploy --only functions:{함수명}
# 전체
firebase deploy --only functions
```

#### 함수 제거

```shell
firebase functions:delete {함수명}
```

#### 함수 로깅

```shell
# 개별
firebase functions:log --only {함수명}
# 전체
firebase functions:log
```
