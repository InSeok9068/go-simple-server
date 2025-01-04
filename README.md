# Go 언어 풀스택 심플 서버

Go + Htmx + AlpineJS + Bulma + SQLite3 기반 심플 서버 프로젝트입니다.

## 실행

### 서비스 실행

#### 리눅스

```shell
chmod +x main
./main
```

### 어드민 서버 실행

#### 윈도우

```shell
.\pocketbase.exe serve --dir ./internal/main/pb_data
```

#### 리눅스

```shell
chmod +x pocketbase
./pocketbase serve --dir ./internal/main/pb_data
```

---

## 종속성 최신화

```shell
go get firebase.google.com/go/v4@latest
go get github.com/a-h/templ@latest
go get github.com/joho/godotenv@latest
go get github.com/labstack/echo/v4@latest
go get github.com/mattn/go-sqlite3@latest
go get google.golang.org/api@latest
```

```shell
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/air-verse/air@latest
```

---

## GCC 활성화

```shell
go env -w CGO_ENABLED=1
```

Window : [tdm-gcc](https://jmeubank.github.io/tdm-gcc/)

---

## 참고 글

https://ntorga.com/full-stack-go-app-with-htmx-and-alpinejs/
