#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 사용법 출력 함수
function show_usage {
  echo -e "${BLUE}사용법:${NC} $0 [배포모드] [서비스명] [호스트IP] [SSH키경로] [포트(기본값:22)]"
  echo -e "배포 모드:"
  echo -e "  ${GREEN}deploy${NC} - 서비스 배포"
  echo -e "  ${GREEN}undeploy${NC} - 서비스 배포 제거"
  echo -e "예시:"
  echo -e "  $0 deploy my-service 123.456.789.012 ~/.ssh/my_key.pem"
  echo -e "  $0 deploy my-service 123.456.789.012 ~/.ssh/my_key.pem 2222"
  echo -e "  $0 undeploy my-service 123.456.789.012 ~/.ssh/my_key.pem"
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
if [ $# -lt 4 ]; then
  show_usage
fi

DEPLOY_MODE="$1"
SERVICE_NAME="$2"
HOST_IP="$3"
SSH_KEY="$4"
SSH_PORT="${5:-22}"  # 포트가 제공되지 않으면 기본값 22 사용

# 배포 모드 검증
if [ "$DEPLOY_MODE" != "deploy" ] && [ "$DEPLOY_MODE" != "undeploy" ]; then
  error_exit "알 수 없는 배포 모드입니다. 'deploy' 또는 'undeploy'를 사용하세요."
fi

# SSH 연결 문자열
SSH_CMD="ssh -i $SSH_KEY -p $SSH_PORT ubuntu@$HOST_IP"
SCP_CMD="scp -i $SSH_KEY -P $SSH_PORT"

# 파일 존재 여부 확인
if [ ! -f "$SSH_KEY" ]; then
  error_exit "SSH 키 파일을 찾을 수 없습니다: $SSH_KEY"
fi

# 1. 서버 연결 테스트 (모든 모드에서 공통)
step "1" "서버 연결 테스트"
if $SSH_CMD "echo '연결 성공!'" > /dev/null 2>&1; then
  check "서버 연결 테스트 성공"
else
  error_exit "서버 연결에 실패했습니다. SSH 키와 IP를 확인하세요."
fi

# 배포 모드에 따른 분기
if [ "$DEPLOY_MODE" = "deploy" ]; then
  # 필요한 파일들 확인
  if [ ! -d "cmd/$SERVICE_NAME" ]; then
    error_exit "서비스 디렉토리를 찾을 수 없습니다: cmd/$SERVICE_NAME"
  fi

  if [ ! -f ".linux/systemctl/$SERVICE_NAME.service" ]; then
    error_exit "서비스 파일을 찾을 수 없습니다: .linux/systemctl/$SERVICE_NAME.service"
  fi

  if [ ! -f ".linux/systemctl/$SERVICE_NAME-admin.service" ]; then
    error_exit "어드민 서비스 파일을 찾을 수 없습니다: .linux/systemctl/$SERVICE_NAME-admin.service"
  fi

  # 메인 스크립트 시작
  echo -e "${BLUE}=== $SERVICE_NAME 서비스 배포 스크립트 ===${NC}"
  echo -e "서비스명: ${YELLOW}$SERVICE_NAME${NC}"
  echo -e "대상 서버: ${YELLOW}$HOST_IP${NC}"
  echo -e "SSH 키: ${YELLOW}$SSH_KEY${NC}"
  echo -e "SSH 포트: ${YELLOW}$SSH_PORT${NC}"
  echo -e "${YELLOW}계속하려면 Enter 키를 누르세요... (취소하려면 Ctrl+C)${NC}"
  read

  # 2. 원격 서버에 디렉토리 생성
  step "2" "원격 서버에 필요한 디렉토리 생성"
  $SSH_CMD "mkdir -p ~/app/projects/$SERVICE_NAME/static ~/app/projects/$SERVICE_NAME/pb_data ~/app/.linux/systemctl"
  check "디렉토리 생성 완료"

  # 3. 서비스 파일 전송
  step "3" "서비스 파일 전송"
  # 서비스 파일 전송
  $SCP_CMD ".linux/systemctl/$SERVICE_NAME.service" "ubuntu@$HOST_IP:~/app/.linux/systemctl/"
  check "메인 서비스 파일 전송 완료"

  $SCP_CMD ".linux/systemctl/$SERVICE_NAME-admin.service" "ubuntu@$HOST_IP:~/app/.linux/systemctl/"
  check "어드민 서비스 파일 전송 완료"

  # 4. Caddyfile 업데이트
  step "4" "Caddyfile 전송 및 업데이트"
  # 먼저 현재 Caddyfile 가져오기
  TEMP_DIR=$(mktemp -d)
  $SCP_CMD "ubuntu@$HOST_IP:/etc/caddy/Caddyfile" "$TEMP_DIR/"
  check "현재 Caddyfile 가져오기 완료"

  # 로컬 Caddyfile에서 서비스 관련 설정 추출
  SERVICE_CONFIG=$(cat .linux/caddy/Caddyfile | grep -A10 "$SERVICE_NAME" | grep -B10 -m1 "}" | grep -B10 "}")

  # 원격 Caddyfile에 설정 추가 (중복 방지)
  if grep -q "$SERVICE_NAME.toy-project.n-e.kr" "$TEMP_DIR/Caddyfile"; then
    warn "$SERVICE_NAME 설정이 이미 Caddyfile에 존재합니다. 업데이트를 건너뜁니다."
  else
    echo -e "\n$SERVICE_CONFIG" >> "$TEMP_DIR/Caddyfile"
    $SCP_CMD "$TEMP_DIR/Caddyfile" "ubuntu@$HOST_IP:~/caddy_update"
    $SSH_CMD "sudo cp ~/caddy_update /etc/caddy/Caddyfile"
    check "Caddyfile 업데이트 완료"
  fi

  # 임시 디렉토리 정리
  rm -rf "$TEMP_DIR"

  # 5. 서비스 바이너리 빌드 및 전송
  step "5" "서비스 바이너리 빌드 및 전송"
  # 빌드 환경 확인
  if command -v go &> /dev/null; then
    echo "Go 컴파일러를 사용하여 바이너리를 빌드합니다."
    
    # 서비스 바이너리 빌드
    go build -o "$SERVICE_NAME" "cmd/$SERVICE_NAME/main.go"
    check "메인 서비스 바이너리 빌드 완료"
    
    # 바이너리 전송
    $SCP_CMD "$SERVICE_NAME" "ubuntu@$HOST_IP:~/app/"
    check "메인 서비스 바이너리 전송 완료"
    
    # 로컬 바이너리 정리
    rm "$SERVICE_NAME"
  else
    warn "Go 컴파일러를 찾을 수 없습니다. 바이너리 빌드를 건너뜁니다."
    warn "빌드된 바이너리를 수동으로 서버에 업로드해야 합니다."
  fi

  # 6. PocketBase 바이너리 확인 및 전송
  step "6" "PocketBase 바이너리 확인 및 전송"
  if [ -f "pocketbase" ]; then
    $SCP_CMD "pocketbase" "ubuntu@$HOST_IP:~/app/$SERVICE_NAME-admin"
    check "PocketBase 바이너리 전송 완료"
  else
    warn "PocketBase 바이너리를 찾을 수 없습니다. 수동으로 업로드해야 합니다."
  fi

  # 7. 정적 파일 전송 (projects/$SERVICE_NAME/static 디렉토리 내용)
  step "7" "정적 파일 전송"
  if [ -d "projects/$SERVICE_NAME/static" ] && [ "$(ls -A "projects/$SERVICE_NAME/static" 2>/dev/null)" ]; then
    $SCP_CMD -r "projects/$SERVICE_NAME/static/"* "ubuntu@$HOST_IP:~/app/projects/$SERVICE_NAME/static/"
    check "정적 파일 전송 완료"
  else
    warn "정적 파일 디렉토리가 비어있거나 존재하지 않습니다."
  fi

  # 8. 서비스 설치 및 활성화
  step "8" "서비스 설치 및 활성화"
  $SSH_CMD "
    # 서비스 실행 파일 권한 부여
    chmod +x ~/app/$SERVICE_NAME
    chmod +x ~/app/$SERVICE_NAME-admin

    # 서비스 파일 설치 및 활성화
    sudo cp ~/app/.linux/systemctl/$SERVICE_NAME.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl start $SERVICE_NAME.service
    sudo systemctl enable $SERVICE_NAME.service

    # 어드민 서비스 파일 설치 및 활성화
    sudo cp ~/app/.linux/systemctl/$SERVICE_NAME-admin.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl start $SERVICE_NAME-admin.service
    sudo systemctl enable $SERVICE_NAME-admin.service

    # Caddy 서비스 재시작
    sudo systemctl reload caddy
  "
  check "서비스 설치 및 활성화 완료"

  # 9. 서비스 상태 확인
  step "9" "서비스 상태 확인"
  echo -e "${YELLOW}메인 서비스 상태:${NC}"
  $SSH_CMD "sudo systemctl status $SERVICE_NAME.service --no-pager"

  echo -e "\n${YELLOW}어드민 서비스 상태:${NC}"
  $SSH_CMD "sudo systemctl status $SERVICE_NAME-admin.service --no-pager"

  echo -e "\n${YELLOW}Caddy 서비스 상태:${NC}"
  $SSH_CMD "sudo systemctl status caddy --no-pager"

  # 완료 메시지
  echo -e "\n${GREEN}=== $SERVICE_NAME 서비스 배포가 완료되었습니다 ===${NC}"
  echo -e "${BLUE}접속 정보:${NC}"
  echo -e "  메인 서비스: ${YELLOW}http://$SERVICE_NAME.toy-project.n-e.kr${NC}"
  echo -e "  어드민 서비스: ${YELLOW}http://$SERVICE_NAME-admin.toy-project.n-e.kr${NC}"

elif [ "$DEPLOY_MODE" = "undeploy" ]; then
  # 메인 스크립트 시작
  echo -e "${RED}=== $SERVICE_NAME 서비스 디배포(배포 제거) 스크립트 ===${NC}"
  echo -e "서비스명: ${YELLOW}$SERVICE_NAME${NC}"
  echo -e "대상 서버: ${YELLOW}$HOST_IP${NC}"
  echo -e "SSH 키: ${YELLOW}$SSH_KEY${NC}"
  echo -e "SSH 포트: ${YELLOW}$SSH_PORT${NC}"
  echo -e "${RED}경고: 이 작업은 서버에서 해당 서비스를 완전히 제거합니다!${NC}"
  echo -e "${YELLOW}계속하려면 '$SERVICE_NAME'을 입력하세요: ${NC}"
  read CONFIRM

  if [ "$CONFIRM" != "$SERVICE_NAME" ]; then
    error_exit "취소되었습니다."
  fi

  # 2. 서비스 중지 및 비활성화
  step "2" "서비스 중지 및 비활성화"
  $SSH_CMD "
    # 서비스 중지 및 비활성화
    sudo systemctl stop $SERVICE_NAME.service 2>/dev/null || true
    sudo systemctl disable $SERVICE_NAME.service 2>/dev/null || true
    sudo systemctl stop $SERVICE_NAME-admin.service 2>/dev/null || true
    sudo systemctl disable $SERVICE_NAME-admin.service 2>/dev/null || true
    
    echo '서비스가 중지되고 비활성화되었습니다.'
  "
  check "서비스 중지 및 비활성화 완료"

  # 3. 서비스 파일 제거
  step "3" "서비스 파일 제거"
  $SSH_CMD "
    # 서비스 파일 제거
    sudo rm -f /etc/systemd/system/$SERVICE_NAME.service
    sudo rm -f /etc/systemd/system/$SERVICE_NAME-admin.service
    sudo systemctl daemon-reload
    
    echo '서비스 파일이 제거되었습니다.'
  "
  check "서비스 파일 제거 완료"

  # 4. Caddyfile 업데이트
  step "4" "Caddyfile에서 서비스 설정 제거"
  # 현재 Caddyfile 가져오기
  TEMP_DIR=$(mktemp -d)
  $SCP_CMD "ubuntu@$HOST_IP:/etc/caddy/Caddyfile" "$TEMP_DIR/Caddyfile"
  check "현재 Caddyfile 가져오기 완료"

  # 백업 생성
  cp "$TEMP_DIR/Caddyfile" "$TEMP_DIR/Caddyfile.bak"
  
  # 서비스 관련 설정 제거
  sed -i "/$SERVICE_NAME\.toy-project\.n-e\.kr/,/}/d" "$TEMP_DIR/Caddyfile"
  sed -i "/$SERVICE_NAME-admin\.toy-project\.n-e\.kr/,/}/d" "$TEMP_DIR/Caddyfile"
  
  # 연속된 빈 줄 처리
  perl -i -0pe 's/\n{3,}/\n\n/g' "$TEMP_DIR/Caddyfile"
  
  # 수정된 Caddyfile 업로드 및 적용
  $SCP_CMD "$TEMP_DIR/Caddyfile" "ubuntu@$HOST_IP:~/caddy_update"
  $SSH_CMD "sudo cp ~/caddy_update /etc/caddy/Caddyfile && sudo systemctl reload caddy"
  check "Caddyfile에서 서비스 설정 제거 완료"

  # 임시 디렉토리 정리
  rm -rf "$TEMP_DIR"

  # 5. 서비스 바이너리 및 데이터 제거
  step "5" "서비스 바이너리 및 데이터 제거"
  echo -e "${YELLOW}서비스 데이터를 백업하시겠습니까? (y/N)${NC}"
  read BACKUP_DATA
  
  if [ "$BACKUP_DATA" = "y" ] || [ "$BACKUP_DATA" = "Y" ]; then
    $SSH_CMD "
      # 서비스 데이터 백업
      BACKUP_DIR=~/backups/$SERVICE_NAME-\$(date +%Y%m%d-%H%M%S)
      mkdir -p \$BACKUP_DIR
      
      if [ -d ~/app/projects/$SERVICE_NAME ]; then
        cp -r ~/app/projects/$SERVICE_NAME \$BACKUP_DIR/
        echo '서비스 데이터가 \$BACKUP_DIR에 백업되었습니다.'
      else
        echo '백업할 서비스 데이터가 없습니다.'
      fi
    "
    check "서비스 데이터 백업 완료"
  fi
  
  # 서비스 파일 및 디렉토리 제거
  $SSH_CMD "
    # 서비스 바이너리 제거
    rm -f ~/app/$SERVICE_NAME
    rm -f ~/app/$SERVICE_NAME-admin
    
    # 서비스 데이터 제거
    rm -rf ~/app/projects/$SERVICE_NAME
    
    echo '서비스 파일 및 디렉토리가 제거되었습니다.'
  "
  check "서비스 바이너리 및 데이터 제거 완료"

  # 6. Caddy 서비스 재시작 및 상태 확인
  step "6" "서비스 제거 완료 확인"
  $SSH_CMD "sudo systemctl reload caddy"
  check "Caddy 서비스 재시작 완료"
  
  echo -e "\n${YELLOW}제거된 서비스 확인:${NC}"
  $SSH_CMD "systemctl list-units --type=service | grep $SERVICE_NAME || echo '서비스가 더 이상 시스템에 존재하지 않습니다.'"

  # 완료 메시지
  echo -e "\n${GREEN}=== $SERVICE_NAME 서비스 디배포(배포 제거)가 완료되었습니다 ===${NC}"
fi
