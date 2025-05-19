#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 사용법 출력 함수
function show_usage {
  echo -e "${BLUE}사용법:${NC} $0 [서비스명] [포트번호(선택)] [어드민포트(선택)]"
  echo -e "  예시: $0 my-service 8003 9003"
  echo -e "  ${YELLOW}포트를 지정하지 않으면 자동으로 할당됩니다.${NC}"
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

# 파라미터 확인
if [ $# -lt 1 ]; then
  show_usage
fi

SERVICE_NAME="$1"
SERVICE_PORT="$2"
ADMIN_PORT="$3"

# 서비스 이름 유효성 검사
if [[ ! $SERVICE_NAME =~ ^[a-z0-9\-]+$ ]]; then
  error_exit "서비스 이름은 소문자, 숫자, 하이픈(-)만 포함해야 합니다."
fi

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
echo -e "${BLUE}=== $SERVICE_NAME 서비스 배포 자동화 스크립트 ===${NC}"
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