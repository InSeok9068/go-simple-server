#!/bin/bash

# ==============================================================================
# Task Runner for go-simple-server
#
# 사용법: ./task.sh [명령어] [인자...]
# 예시:
#   ./task.sh help              - 모든 명령어 목록 보기
#   ./task.sh switch deario     - deario 서비스로 개발 환경 전환
#   ./task.sh check             - 프로젝트 전체 검사 (빌드, 테스트, 린트)
#   ./task.sh service create my-app - my-app 이라는 새 서비스 생성
# ==============================================================================

# 색상 정의
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 스크립트 경로
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/scripts"

# 메인 명령어
COMMAND=$1
shift # 첫 번째 인자를 제거하여 나머지 인자들을 쉽게 전달

# 도움말 함수
show_help() {
  echo -e "${BLUE}Go Simple Server - Task Runner${NC}"
  echo -e "--------------------------------"
  echo -e "이 스크립트는 프로젝트의 다양한 관리 작업을 자동화합니다."
  echo -e "\n${YELLOW}사용법:${NC} ./task.sh [명령어] [인자...]\n"
  echo -e "${YELLOW}주요 명령어:${NC}"
  echo -e "  ${GREEN}help${NC}                          - 이 도움말 메시지를 표시합니다."
  echo -e "  ${GREEN}switch${NC} [service]              - 개발 환경에서 실행할 서비스를 변경합니다. (예: deario, homepage)"
  echo -e "  ${GREEN}check${NC} [build|test|lint]       - 프로젝트의 오류를 검사합니다. (인자 없으면 전체 실행)"
  echo -e "  ${GREEN}deps${NC}                          - Go 의존성 및 도구를 업데이트합니다."
  echo -e "  ${GREEN}build:linux${NC}                   - Linux용 바이너리를 빌드합니다."
  echo -e "  ${GREEN}release${NC} [project]             - 릴리스 브랜치를 main에 병합합니다. (인자 없으면 전체 실행)"
  echo -e "  ${GREEN}install:tailwind${NC} [win|linux]  - TailwindCSS를 설치합니다."
  echo -e "\n${YELLOW}서비스 관리:${NC}"
  echo -e "  ${GREEN}service create${NC} [name] [port]  - 새 서비스를 생성합니다."
  echo -e "  ${GREEN}service deploy${NC} [name]         - 서비스 배포 가이드를 출력합니다."
  echo -e "  ${GREEN}service remove${NC} [name]         - 기존 서비스를 제거합니다."
  echo -e "  ${GREEN}service undeploy${NC} [name]       - 서비스 배포 제거 가이드를 출력합니다."
}

# 명령어가 없으면 도움말 표시
if [ -z "$COMMAND" ]; then
  show_help
  exit 0
fi

# 명령어에 따른 분기 처리
case "$COMMAND" in
  help)
    show_help
    ;;
  switch)
    bash "$SCRIPT_DIR/change.sh" "$@"
    ;;
  check)
    bash "$SCRIPT_DIR/error-check.sh" "$@"
    ;;
  deps)
    bash "$SCRIPT_DIR/update-deps.sh" "$@"
    ;;
  build:linux)
    bash "$SCRIPT_DIR/build-linux.sh" "$@"
    ;;
  release)
    bash "$SCRIPT_DIR/release-all.sh" "$@"
    ;;
  install:tailwind)
    bash "$SCRIPT_DIR/tailwindcss-install.sh" "$@"
    ;;
  service)
    bash "$SCRIPT_DIR/service-manager.sh" "$@"
    ;;
  *)
    echo -e "❌ 알 수 없는 명령어: ${YELLOW}$COMMAND${NC}"
    echo -e "사용 가능한 명령어 목록을 보려면 './task.sh help'를 실행하세요."
    exit 1
    ;;
esac
