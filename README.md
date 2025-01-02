# Go 언어 풀스택 심플 서버

Go + Htmx + AlpineJS + Bulma + SQLite3 기반 심플 서버 프로젝트입니다.

## 실행

### 어드민 서버 실행

```shell
.\pocketbase.exe serve
```

--- 

```shell
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/air-verse/air@latest
```

```shell
go env -w CGO_ENABLED=1
```

Window : [tdm-gcc](https://jmeubank.github.io/tdm-gcc/)

## 참고 글

https://ntorga.com/full-stack-go-app-with-htmx-and-alpinejs/
