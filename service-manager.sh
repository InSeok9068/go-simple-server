#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 사용법 출력 함수
function show_usage {
  echo -e "${BLUE}사용법:${NC} $0 [명령어] [서비스명] [옵션...]"
  echo -e "명령어:"
  echo -e "  ${GREEN}create${NC} [서비스명] [포트번호(선택)] [어드민포트(선택)] - 새 서비스 생성"
  echo -e "  ${GREEN}deploy${NC} [서비스명] - 서비스 배포 명령어 출력"
  echo -e "  ${GREEN}remove${NC} [서비스명] - 서비스 제거"
  echo -e "  ${GREEN}undeploy${NC} [서비스명] - 서비스 제거 명령어 출력"
  echo -e "예시:"
  echo -e "  $0 create my-service 8003 9003"
  echo -e "  $0 deploy my-service"
  echo -e "  $0 remove my-service"
  echo -e "  $0 undeploy my-service"
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
if [ $# -lt 2 ]; then
  show_usage
fi

COMMAND="$1"
SERVICE_NAME="$2"

# 서비스 이름 유효성 검사
if [[ ! $SERVICE_NAME =~ ^[a-z0-9\-]+$ ]]; then
  error_exit "서비스 이름은 소문자, 숫자, 하이픈(-)만 포함해야 합니다."
fi

# 명령어에 따른 분기
case "$COMMAND" in
  "create")
    # new-service.sh 기능
    SERVICE_PORT="$3"
    ADMIN_PORT="$4"

    # 포트 자동 할당 (제공되지 않은 경우)
    if [ -z "$SERVICE_PORT" ]; then
      # 기존 서비스 포트 찾기 (sample 제외)
      LAST_PORT=$(awk '/^  [a-z\-]+\)$/ && !/^  sample\)$/ {
        getline;
        if ($0 ~ /PORT=[0-9]+/) {
          match($0, /[0-9]+/);
          print substr($0, RSTART, RLENGTH);
        }
      }' change.sh | sort -nr | head -1)

      if [ -n "$LAST_PORT" ]; then
        SERVICE_PORT=$((LAST_PORT + 1))
      else
        SERVICE_PORT=8010 # 기본 시작 포트
      fi
    fi

    # 어드민 포트 자동 할당 (제공되지 않은 경우)
    if [ -z "$ADMIN_PORT" ]; then
      ADMIN_PORT=$((SERVICE_PORT + 1000)) # 서비스 포트 + 1000
    fi

    # 메인 스크립트 시작
    echo -e "${BLUE}=== $SERVICE_NAME 서비스 생성 스크립트 ===${NC}"
    echo -e "서비스명: ${YELLOW}$SERVICE_NAME${NC}"
    echo -e "서비스 포트: ${YELLOW}$SERVICE_PORT${NC}"
    echo -e "어드민 포트: ${YELLOW}$ADMIN_PORT${NC}"
    echo -e "${YELLOW}계속하려면 Enter 키를 누르세요... (취소하려면 Ctrl+C)${NC}"
    read

    # 1. 서비스 폴더 생성
    step "1" "projects/$SERVICE_NAME 폴더 생성"
    mkdir -p "projects/$SERVICE_NAME" || error_exit "폴더 생성 실패"
    check "폴더 구조 생성 완료"

    # 2. main.go 파일 생성
    step "2" "cmd/$SERVICE_NAME/main.go 파일 생성"
    mkdir -p "cmd/$SERVICE_NAME" || error_exit "cmd 폴더 생성 실패"

    cat > "cmd/$SERVICE_NAME/main.go" << EOF
package main

func main() {
}
EOF
    check "main.go 파일 생성 완료"

    # 3. embed.go 업데이트
    step "3" "embed.go 파일 업데이트"
    if [ -f "embed.go" ]; then
      # 마지막 //go:embed 줄 찾기
      LAST_EMBED_LINE=$(grep -n "//go:embed projects/" embed.go | tail -1 | cut -d':' -f1)

      # 마지막 //go:embed 줄 뒤에 새 임베드 지시문 추가
      sed -i "${LAST_EMBED_LINE}a//go:embed projects/$SERVICE_NAME/static/*" embed.go || error_exit "embed.go 업데이트 실패"
      check "embed.go 파일에 새 임베드 지시문 추가 완료"
    else
      error_exit "embed.go 파일을 찾을 수 없습니다"
    fi

    # 4. change.sh 스크립트 업데이트
    step "4" "change.sh 스크립트 업데이트"
    if [ -f "change.sh" ]; then
      # sample) 앞에 새 서비스 추가
      sed -i '/sample)/i\  '"$SERVICE_NAME"')\n    PORT='"$SERVICE_PORT"'\n    ;;' change.sh || error_exit "change.sh 업데이트 실패"
      check "change.sh 파일에 새 서비스 추가 완료"
    else
      error_exit "change.sh 파일을 찾을 수 없습니다"
    fi

    # 5. 서비스 파일 생성
    step "5" "systemd 서비스 파일 생성"
    mkdir -p ".linux/systemctl" || error_exit "systemctl 폴더 생성 실패"

    # 메인 서비스 파일 생성
    cat > ".linux/systemctl/$SERVICE_NAME.service" << EOF
[Unit]
Description=$SERVICE_NAME Service
After=network.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
WorkingDirectory=/home/ubuntu/app
ExecStart=/home/ubuntu/app/$SERVICE_NAME
Restart=on-failure
RestartSec=5s
Environment=PORT=$SERVICE_PORT

[Install]
WantedBy=multi-user.target
EOF
    check "메인 서비스 파일 생성 완료"

    # 어드민 서비스 파일 생성
    cat > ".linux/systemctl/$SERVICE_NAME-admin.service" << EOF
[Unit]
Description=$SERVICE_NAME Admin Service
After=network.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
WorkingDirectory=/home/ubuntu/app
ExecStart=/home/ubuntu/app/$SERVICE_NAME-admin serve --dir /home/ubuntu/app/projects/$SERVICE_NAME/pb_data --http=127.0.0.1:$ADMIN_PORT
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF
    check "어드민 서비스 파일 생성 완료"

    # 6. Caddyfile 업데이트
    step "6" "Caddyfile 업데이트"
    mkdir -p ".linux/caddy" || error_exit "caddy 폴더 생성 실패"

    if [ -f ".linux/caddy/Caddyfile" ]; then
      # 존재하는 Caddyfile의 마지막에 새 서비스 설정 추가
      cat >> ".linux/caddy/Caddyfile" << EOF


# $SERVICE_NAME 서브도메인
$SERVICE_NAME.toy-project.n-e.kr {
    reverse_proxy 127.0.0.1:$SERVICE_PORT
}

# $SERVICE_NAME 어드민 서브도메인
$SERVICE_NAME-admin.toy-project.n-e.kr {
    reverse_proxy 127.0.0.1:$ADMIN_PORT
}
EOF
      check "Caddyfile 업데이트 완료"
    else
      error_exit "Caddyfile을 찾을 수 없습니다"
    fi
    ;;

  "deploy")
    # new-service-deploy.sh 기능
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
    ;;

  "remove")
    # remove-service.sh 기능
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
    ;;

  "undeploy")
    # remove-service-deploy.sh 기능
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
    ;;

  *)
    echo -e "${RED}알 수 없는 명령어:${NC} $COMMAND"
    show_usage
    ;;
esac