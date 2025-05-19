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
echo -e "${RED}=== $SERVICE_NAME 서비스 제거 명령어 목록 ====${NC}"
echo -e "${YELLOW}서버에 접속한 후 아래 명령어를 순서대로 실행하세요${NC}"

echo -e "\n${YELLOW}1. 서비스 중지 및 비활성화${NC}"
echo "sudo systemctl stop $SERVICE_NAME.service"
echo "sudo systemctl disable $SERVICE_NAME.service"

echo -e "\n${YELLOW}2. 어드민 서비스 중지 및 비활성화${NC}"
echo "sudo systemctl stop $SERVICE_NAME-admin.service"
echo "sudo systemctl disable $SERVICE_NAME-admin.service"

echo -e "\n${YELLOW}3. 서비스 파일 제거${NC}"
echo "sudo rm /etc/systemd/system/$SERVICE_NAME.service"
echo "sudo rm /etc/systemd/system/$SERVICE_NAME-admin.service"
echo "sudo systemctl daemon-reload"

echo -e "\n${YELLOW}4. 서비스 실행 파일 제거${NC}"
echo "rm /home/ubuntu/app/$SERVICE_NAME"
echo "rm /home/ubuntu/app/$SERVICE_NAME-admin"

echo -e "\n${YELLOW}5. 서비스 프로젝트 폴더 제거${NC}"
echo "rm -rf /home/ubuntu/app/projects/$SERVICE_NAME"

echo -e "\n${YELLOW}6. Caddyfile 수정${NC}"
echo "sudo nano /etc/caddy/Caddyfile"
echo -e "${GREEN}서브도메인 설정을 찾아서 삭제해주세요${NC}"

echo -e "\n${YELLOW}7. Caddy 서비스 재시작${NC}"
echo "sudo systemctl reload caddy"

echo -e "\n${BLUE}=== 서비스 상태 확인 명령어 ====${NC}"
echo "systemctl list-units | grep $SERVICE_NAME"

echo -e "\n${GREEN}위 명령어들을 복사하여 서버에서 순서대로 실행하세요.${NC}"