#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 사용법 출력 함수
function show_usage {
  echo -e "${BLUE}사용법:${NC} $0 [서비스명]"
  echo -e "  예시: $0 my-service"
  exit 1
}

# 오류 처리 함수
function error_exit {
  echo -e "${RED}오류:${NC} $1"
  exit 1
}

# 파라미터 확인
if [ $# -lt 1 ]; then
  show_usage
fi

SERVICE_NAME="$1"

# 서비스 이름 유효성 검사
if [[ ! $SERVICE_NAME =~ ^[a-z0-9\-]+$ ]]; then
  error_exit "서비스 이름은 소문자, 숫자, 하이픈(-)만 포함해야 합니다."
fi

# 명령어 출력 시작
echo -e "${BLUE}=== $SERVICE_NAME 서비스 배포 명령어 목록 ====${NC}"
echo -e "${YELLOW}서버에 접속한 후 아래 명령어를 순서대로 실행하세요${NC}"

echo -e "\n${YELLOW}1. 서비스 실행 파일 권한 부여${NC}"
echo "chmod +x /home/ubuntu/app/$SERVICE_NAME"

echo -e "\n${YELLOW}2. 어드민 실행 파일 권한 부여${NC}"
echo "chmod +x /home/ubuntu/app/$SERVICE_NAME-admin"

echo -e "\n${YELLOW}3. 서비스 파일 설치 및 활성화${NC}"
echo "sudo cp .linux/systemctl/$SERVICE_NAME.service /etc/systemd/system/"
echo "sudo systemctl daemon-reload"
echo "sudo systemctl start $SERVICE_NAME.service"
echo "sudo systemctl enable $SERVICE_NAME.service"

echo -e "\n${YELLOW}4. 어드민 서비스 파일 설치 및 활성화${NC}"
echo "sudo cp .linux/systemctl/$SERVICE_NAME-admin.service /etc/systemd/system/"
echo "sudo systemctl daemon-reload"
echo "sudo systemctl start $SERVICE_NAME-admin.service"
echo "sudo systemctl enable $SERVICE_NAME-admin.service"

echo -e "\n${YELLOW}5. Caddyfile 수정${NC}"
echo "sudo nano /etc/caddy/Caddyfile"
echo -e "${GREEN}서브도메인과 포트 설정을 추가해주세요${NC}"

echo -e "\n${YELLOW}6. Caddy 서비스 재시작${NC}"
echo "sudo systemctl reload caddy"

echo -e "\n${BLUE}=== 서비스 상태 확인 명령어 ====${NC}"
echo "sudo systemctl status $SERVICE_NAME.service"
echo "sudo systemctl status $SERVICE_NAME-admin.service"
echo "sudo systemctl status caddy"

echo -e "\n${YELLOW}7. 어드민 계정 생성 (필요한 경우)${NC}"
echo "# 브라우저에서 $SERVICE_NAME-admin.toy-project.n-e.kr 접속"
echo "# 어드민 계정을 설정하세요"

echo -e "\n${GREEN}위 명령어들을 복사하여 서버에서 순서대로 실행하세요.${NC}"