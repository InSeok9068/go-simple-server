# Go 언어 기반 풀스택 서버

**Host** : https://toy-project.n-e.kr

| 구성                  | 사용 기술        |
| --------------------- | ---------------- |
| 언어                  | Go               |
| 백엔드 프레임워크     | Echo             |
| 프론트엔드 프레임워크 | Htmx + Alpine.js |
| CSS 프레임워크        | Bulma            |
| Template 엔진         | Templ            |
| 데이터베이스          | SQLite3          |

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
./pocketbase.exe serve --dir ./projects/sample/pb_data
```

#### 리눅스

```shell
chmod +x pocketbase
./pocketbase serve --dir ./projects/sample/pb_data
```

---

## 종속성 최신화

### 프로젝트 종속성 최신화

```shell
go get firebase.google.com/go/v4@latest
go get github.com/a-h/templ@latest
go get github.com/chromedp/chromedp@latest
go get github.com/joho/godotenv@latest
go get github.com/labstack/echo/v4@latest
go get github.com/mattn/go-sqlite3@latest
go get github.com/robfig/cron/v3@latest
go get google.golang.org/api@latest
```

### 개발 도구 종속성 최신화

```shell
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/air-verse/air@latest
CGO_ENABLED=0 go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### 불필요한 종속성 제거

```shell
go mod tidy
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
