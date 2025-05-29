# Go 언어 기반 풀스택 서버

**Host** : https://toy-project.n-e.kr

| 구성              | 사용 기술                   |
|-----------------|-------------------------|
| 언어              | Go                      |
| 백엔드 프레임워크       | Echo                    |
| 프론트엔드 프레임워크     | Htmx + Alpine.js        |
| CSS 프레임워크 선택지 1 | Beer CSS                |
| CSS 프레임워크 선택지 2 | Pico CSS + Tailwind CSS |
| CSS 프레임워크 선택지 3 | Bulma CSS               |
| Template 엔진     | Gomponents              |
| 데이터베이스          | SQLite3                 |

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

### 어드민 서버 실행

#### 윈도우

```shell
./pocketbase.exe serve --dir ./projects/homepage/pb_data
```

#### 리눅스

```shell
chmod +x pocketbase
./pocketbase serve --dir ./projects/homepage/pb_data
```

---

## 종속성 최신화

### 프로젝트 종속성 최신화

```shell
go get firebase.google.com/go/v4@latest
go get github.com/doganarif/govisual@latest
go get github.com/eduardolat/gomponents-lucide@latest
go get github.com/glsubri/gomponents-alpine@latest
go get github.com/go-rod/rod@latest
go get github.com/gorilla/sessions@latest
go get github.com/joho/godotenv@latest
go get github.com/labstack/echo-contrib@latest
go get github.com/labstack/echo/v4@latest
go get github.com/mattn/go-sqlite3@latest
go get github.com/openai/openai-go@latest
go get github.com/robfig/cron/v3@latest
go get github.com/willoma/bulma-gomponents@latest
go get github.com/willoma/gomplements@latest
go get google.golang.org/api@latest
go get google.golang.org/genai@latest
go get maragu.dev/gomponents@latest
go get maragu.dev/gomponents-htmx@latest
```

### 개발 도구 종속성 최신화

```shell
#go install github.com/a-h/templ/cmd/templ@latest
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
CGO_ENABLED=0 go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
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
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.63.4
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
