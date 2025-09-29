#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 함수: 진행 메시지 출력
print_step() {
    echo -e "\n${BLUE}==>${NC} ${GREEN}$1${NC}"
}

# 함수: 완료 메시지 출력
print_done() {
    echo -e "${GREEN}✅ $1 완료${NC}"
}

# 시작 메시지
echo -e "\n${YELLOW}=== Go 의존성 업데이트를 시작합니다 ===${NC}"

# 1. Go 모듈 업데이트
print_step "1. Go 모듈 업데이트 중..."

modules=(
    "cloud.google.com/go/storage@latest"
    "firebase.google.com/go/v4@latest"
    "github.com/AlecAivazis/survey/v2@latest"
    "github.com/Blank-Xu/sql-adapter@latest"
    "github.com/casbin/casbin/v2@latest"
    "github.com/crazy-max/echo-ipfilter@latest"
    "github.com/doganarif/govisual@latest"
    "github.com/eduardolat/gomponents-lucide@latest"
    "github.com/glsubri/gomponents-alpine@latest"
    "github.com/go-rod/rod@latest"
    "github.com/gorilla/sessions@latest"
    "github.com/joho/godotenv@latest"
    "github.com/labstack/echo-contrib@latest"
    "github.com/labstack/echo/v4@latest"
    "github.com/lmittmann/tint@latest"
    "github.com/pressly/goose/v3@latest"
    "github.com/qustavo/sqlhooks/v2@latest"
    "github.com/robfig/cron/v3@latest"
    "github.com/willoma/bulma-gomponents@latest"
    "github.com/willoma/gomplements@latest"
    "go.opentelemetry.io/otel@latest"
    "go.opentelemetry.io/otel/sdk@latest"
    "go.opentelemetry.io/otel/trace@latest"
    "golang.org/x/time@latest"
    "google.golang.org/api@latest"
    "google.golang.org/genai@latest"
    "maragu.dev/gomponents@latest"
    "maragu.dev/gomponents-htmx@latest"
    "maragu.dev/goqite@latest"
    "modernc.org/sqlite@latest"
)

for module in "${modules[@]}"; do
    echo -n "- ${module} 업데이트 중... "
    if go get $module; then
        echo -e "${GREEN}성공${NC}"
    else
        echo -e "${YELLOW}경고: 업데이트 실패${NC}"
    fi
done

print_done "모든 Go 모듈 업데이트"

# 2. Go 도구 설치
print_step "2. Go 도구 설치 중..."

tools=(
    "github.com/air-verse/air@latest"
    "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"
    "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
    "golang.org/x/vuln/cmd/govulncheck@latest"
)

for tool in "${tools[@]}"; do
    echo -n "- ${tool} 설치 중... "
    if go install $tool; then
        echo -e "${GREEN}성공${NC}"
    else
        echo -e "${YELLOW}경고: 설치 실패${NC}"
    fi
done

print_done "모든 Go 도구 설치"

# 3. 모듈 정리
print_step "3. 불필요한 의존성 정리 중..."
go mod tidy

if [ $? -eq 0 ]; then
    print_done "의존성 정리 완료"
else
    echo -e "${YELLOW}경고: 의존성 정리 중 문제가 발생했습니다.${NC}"
fi

# 4. 프론트엔드 라이브러리 버전 확인
print_step "4. 프론트엔드 라이브러리 버전 확인 중..."

frontend_libs=(
    "htmx.org|2.0.7|https://cdn.jsdelivr.net/npm/htmx.org@%s/dist/htmx.min.js"
    "alpinejs|3.x.x|https://cdn.jsdelivr.net/npm/alpinejs@%s/dist/cdn.min.js"
    "beercss|3.12.8|https://cdn.jsdelivr.net/npm/beercss@%s/dist/cdn/beer.min.js"
    "open-props|1.7.16|https://cdn.jsdelivr.net/npm/open-props@%s/open-props.min.css"
    "bulma|1.0.4|https://cdn.jsdelivr.net/npm/bulma@%s/css/bulma.min.css"
    "@picocss/pico|2|https://cdn.jsdelivr.net/npm/@picocss/pico@%s/css/pico.classless.min.css"
    "material-dynamic-colors|1.1.2|https://cdn.jsdelivr.net/npm/material-dynamic-colors@%s/dist/cdn/material-dynamic-colors.min.js"
    "flatpickr|4.6.13|https://cdn.jsdelivr.net/npm/flatpickr@%s/dist/flatpickr.min.css"
)

for entry in "${frontend_libs[@]}"; do
    IFS='|' read -r pkg current url <<< "$entry"
    latest=$(npm info "$pkg" version 2>/dev/null)
    if [ -z "$latest" ]; then
        latest="알 수 없음"
        latest_url="알 수 없음"
    else
        latest_url=$(printf "$url" "$latest")
    fi
    echo "- $pkg: 현재 $current, 최신 $latest"
    echo "  최신 URL: $latest_url"
done

print_done "프론트엔드 라이브러리 버전 확인"

# 완료 메시지
echo -e "\n${GREEN}=== 모든 의존성 업데이트가 완료되었습니다 ===${NC}"
