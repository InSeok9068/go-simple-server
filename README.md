# Go 언어 기반 풀스택 서버

**Host** : https://toy-project.n-e.kr

| 구성                    | 사용 기술               |
| ----------------------- | ----------------------- |
| 언어                    | Go                      |
| 백엔드 프레임워크       | Echo                    |
| 프론트엔드 프레임워크   | Htmx + Alpine.js        |
| CSS 프레임워크 선택지 1 | Beer CSS                |
| CSS 프레임워크 선택지 2 | Pico CSS + Tailwind CSS |
| CSS 프레임워크 선택지 3 | Bulma CSS               |
| Template 엔진           | Gomponents              |
| 데이터베이스            | SQLite3                 |

### 추가 구성 고려

- Beer CSS : CSS 프레임워크 (Material 모바일 우선)
- Pico CSS + Tailwind CSS : CSS 프레임워크 (커스터마이징 용이)
- Bulma CSS : CSS 프레임워크 (심플한 관리자 UI)
- Tabler : CSS 프레임워크 (복잡한 관리자 UI) - PC 환경
- Shoelace : CSS 프레임워크 (바닐라 웹 컴포넌트)
- ~~surreal.js : [surreal.js](https://cdn.jsdelivr.net/gh/gnat/surreal@main/surreal.js)~~
- DataStar : SSE 기반 서버주도 프레임워크
- Unpoly : Htmx+ Like 프레임워크

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
go get firebase.google.com/go/v4
go get github.com/Blank-Xu/sql-adapter
go get github.com/casbin/casbin/v2
go get github.com/doganarif/govisual
go get github.com/eduardolat/gomponents-lucide
go get github.com/glsubri/gomponents-alpine
go get github.com/go-rod/rod
go get github.com/gorilla/sessions
go get github.com/joho/godotenv
go get github.com/labstack/echo-contrib
go get github.com/labstack/echo/v4
go get github.com/lmittmann/tint
go get github.com/openai/openai-go
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
go get modernc.org/sqlite
```

### 개발 도구 종속성 최신화

```shell
#go install github.com/a-h/templ/cmd/templ@latest
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
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

#### JS

[cdn.min.js](https://cdn.jsdelivr.net/npm/alpinejs/dist/cdn.min.js) </br>
[htmx.min.js](https://cdn.jsdelivr.net/npm/htmx.org/dist/htmx.min.js) </br>
[beer.min.js](https://cdn.jsdelivr.net/npm/beercss/dist/cdn/beer.min.js) </br>
[material-dynamic-colors.min.js](https://cdn.jsdelivr.net/npm/material-dynamic-colors/dist/cdn/material-dynamic-colors.min.js)

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
