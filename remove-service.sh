#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 사용법 출력 함수
function show_usage {
  echo -e "${BLUE}사용법:${NC} $0 [삭제할 서비스명]"
  echo -e "  예시: $0 my-service"
  exit 1
}

# 오류 처리 함수
function error_exit {
  echo -e "${RED}오류:${NC} $1"
  exit 1
}

# 진행 메시지 함수
function step {
  echo -e "\n${GREEN}[$(date +%H:%M:%S)]${NC} ${BLUE}단계 $1:${NC} $2"
}

# 체크 표시 함수
function check {
  echo -e "  ${GREEN}✓${NC} $1"
}

# 경고 메시지 함수
function warn {
  echo -e "  ${YELLOW}!${NC} $1"
}

# 파라미터 확인
if [ $# -ne 1 ]; then
  show_usage
fi

SERVICE_NAME="$1"

# 서비스 이름 유효성 검사
if [[ ! $SERVICE_NAME =~ ^[a-z0-9\-]+$ ]]; then
  error_exit "서비스 이름은 소문자, 숫자, 하이픈(-)만 포함해야 합니다."
fi

# 메인 스크립트 시작
echo -e "${RED}=== $SERVICE_NAME 서비스 제거 스크립트 ===${NC}"
echo -e "서비스명: ${YELLOW}$SERVICE_NAME${NC}"
echo -e "${RED}경고: 이 작업은 되돌릴 수 없습니다!${NC}"
echo -e "${YELLOW}계속하려면 '$SERVICE_NAME'을 입력하세요: ${NC}"
read CONFIRM

if [ "$CONFIRM" != "$SERVICE_NAME" ]; then
  error_exit "취소되었습니다."
fi

# 1. change.sh 스크립트에서 서비스 제거
step "1" "change.sh 스크립트에서 서비스 제거"
if [ -f "change.sh" ]; then
  # 서비스 항목 찾기
  SERVICE_PATTERN="  $SERVICE_NAME)"
  if grep -q "$SERVICE_PATTERN" change.sh; then
    # 항목을 찾아서 3줄(case 항목, PORT 설정, ;;) 삭제
    sed -i "/$SERVICE_PATTERN/,/;;/d" change.sh
    check "change.sh에서 서비스 항목 제거 완료"
  else
    warn "change.sh에서 서비스 항목을 찾을 수 없습니다"
  fi
else
  warn "change.sh 파일을 찾을 수 없습니다"
fi

# 2. embed.go 파일에서 임베드 지시문 제거
step "2" "embed.go 파일에서 임베드 지시문 제거"
if [ -f "embed.go" ]; then
  EMBED_PATTERN="//go:embed projects/$SERVICE_NAME/static/\*"
  if grep -q "$EMBED_PATTERN" embed.go; then
    sed -i "\\|$EMBED_PATTERN|d" embed.go
    check "embed.go에서 임베드 지시문 제거 완료"
  else
    warn "embed.go에서 해당 서비스의 임베드 지시문을 찾을 수 없습니다"
  fi
else
  warn "embed.go 파일을 찾을 수 없습니다"
fi

# 3. Caddyfile에서 서비스 항목 제거 (완전히 개선된 버전)
step "3" "Caddyfile에서 서비스 항목 제거"
CADDYFILE=".linux/caddy/Caddyfile"
if [ -f "$CADDYFILE" ]; then
  # 백업 생성
  cp "$CADDYFILE" "${CADDYFILE}.bak"

  # Sed를 이용한 복잡한 패턴 처리 (여러 줄에 걸친 블록 제거)
  perl -i -0pe 's/# '"$SERVICE_NAME"' 서브도메인.*?}//gs' "$CADDYFILE"
  perl -i -0pe 's/# '"$SERVICE_NAME"' 어드민 서브도메인.*?}//gs' "$CADDYFILE"
  perl -i -0pe 's/'"$SERVICE_NAME"'\.toy-project\.n-e\.kr \{.*?}//gs' "$CADDYFILE"
  perl -i -0pe 's/'"$SERVICE_NAME"'-admin\.toy-project\.n-e\.kr \{.*?}//gs' "$CADDYFILE"

  # 중복된 빈 줄 정리 (3개 이상의 연속된 빈 줄을 2개로 줄임)
  perl -i -0pe 's/\n{3,}/\n\n/g' "$CADDYFILE"

  check "Caddyfile에서 서비스 관련 모든 블록 제거 완료"
else
  warn "Caddyfile을 찾을 수 없습니다"
fi

# 4. systemd 서비스 파일 제거
step "4" "systemd 서비스 파일 제거"
SERVICE_FILE=".linux/systemctl/$SERVICE_NAME.service"
ADMIN_SERVICE_FILE=".linux/systemctl/$SERVICE_NAME-admin.service"

if [ -f "$SERVICE_FILE" ]; then
  rm "$SERVICE_FILE"
  check "메인 서비스 파일 제거 완료"
else
  warn "메인 서비스 파일을 찾을 수 없습니다"
fi

if [ -f "$ADMIN_SERVICE_FILE" ]; then
  rm "$ADMIN_SERVICE_FILE"
  check "어드민 서비스 파일 제거 완료"
else
  warn "어드민 서비스 파일을 찾을 수 없습니다"
fi

# 5. main.go 파일 제거
step "5" "main.go 파일 제거"
MAIN_GO="cmd/$SERVICE_NAME/main.go"
MAIN_DIR="cmd/$SERVICE_NAME"

if [ -f "$MAIN_GO" ]; then
  rm "$MAIN_GO"
  check "main.go 파일 제거 완료"

  # 디렉토리가 비어있다면 제거
  if [ -d "$MAIN_DIR" ] && [ -z "$(ls -A "$MAIN_DIR")" ]; then
    rmdir "$MAIN_DIR"
    check "$SERVICE_NAME cmd 디렉토리 제거 완료"
  fi
else
  warn "main.go 파일을 찾을 수 없습니다"
fi

# 6. 서비스 프로젝트 폴더 제거
step "6" "서비스 프로젝트 폴더 제거"
SERVICE_DIR="projects/$SERVICE_NAME"

if [ -d "$SERVICE_DIR" ]; then
  if [ -z "$(ls -A "$SERVICE_DIR")" ]; then
    # 디렉토리가 비어있는 경우
    rmdir "$SERVICE_DIR"
    check "빈 프로젝트 폴더 제거 완료"
  else
    # 디렉토리에 파일이 있는 경우
    echo -e "${YELLOW}프로젝트 폴더에 파일이 있습니다. 모두 삭제하시겠습니까? (y/N)${NC}"
    read DELETE_FILES

    if [ "$DELETE_FILES" = "y" ] || [ "$DELETE_FILES" = "Y" ]; then
      rm -rf "$SERVICE_DIR"
      check "프로젝트 폴더 및 모든 파일 제거 완료"
    else
      warn "프로젝트 폴더 내용이 유지됩니다"
    fi
  fi
else
  warn "프로젝트 폴더를 찾을 수 없습니다"
fi

echo -e "\n${GREEN}=== $SERVICE_NAME 서비스 제거가 완료되었습니다 ===${NC}"
echo -e "${YELLOW}참고: 이 스크립트는 모든 변경 사항을 완전히 제거하지 못할 수 있습니다.${NC}"
echo -e "${YELLOW}      특히 커스텀 코드나 설정은 수동으로 확인하는 것이 좋습니다.${NC}"