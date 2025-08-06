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

# 완료 메시지
echo -e "\n${GREEN}=== 모든 의존성 업데이트가 완료되었습니다 ===${NC}"