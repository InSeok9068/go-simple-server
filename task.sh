#!/bin/bash
set -euo pipefail

# ---- repo root 탐지: git repo면 최상단, 아니면 현재 파일 기준 ----
if ROOT="$(git rev-parse --show-toplevel 2>/dev/null)"; then
  ROOT_DIR="$ROOT"
else
  ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fi
SCRIPT_DIR="${ROOT_DIR}/scripts"

BLUE='\033[0;34m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
COMMAND="${1-}"; shift || true

show_help() {
  echo -e "${BLUE}Go Simple Server - Task Runner${NC}"
  echo -e "--------------------------------"
  echo -e "사용법: ./task.sh [명령어] [인자...]"
  echo -e ""
  echo -e "${YELLOW}주요 명령어:${NC}"
  echo -e "  ${GREEN}help${NC}"
  echo -e "  ${GREEN}switch${NC} [project]"
  echo -e "  ${GREEN}check${NC} [build|test|lint]"
  echo -e "  ${GREEN}deps${NC}"
  echo -e "  ${GREEN}build-linux${NC}"
  echo -e "  ${GREEN}release${NC} [project]"
  echo -e "  ${GREEN}install-tailwind${NC} [win|linux]"   # ← win으로 고정
  echo -e "  ${GREEN}sqlc-generate${NC} [project]"
  echo -e ""
  echo -e "${YELLOW}서비스 관리:${NC}"
  echo -e "  ${GREEN}service create${NC} [name] [port]"
  echo -e "  ${GREEN}service deploy${NC} [name]"
  echo -e "  ${GREEN}service remove${NC} [name]"
  echo -e "  ${GREEN}service undeploy${NC} [name]"
}

if [[ -z "${COMMAND}" ]]; then
  show_help; exit 0
fi

case "${COMMAND}" in
  help)             bash "${SCRIPT_DIR}/show-help.sh" "$@" 2>/dev/null || show_help ;;
  switch)           bash "${SCRIPT_DIR}/change.sh" "$@" ;;
  check)            bash "${SCRIPT_DIR}/error-check.sh" "$@" ;;
  deps)             bash "${SCRIPT_DIR}/update-deps.sh" "$@" ;;
  build-linux)      bash "${SCRIPT_DIR}/build-linux.sh" "$@" ;;
  release)          bash "${SCRIPT_DIR}/release-all.sh" "$@" ;;
  install-tailwind) bash "${SCRIPT_DIR}/tailwindcss-install.sh" "$@" ;;
  sqlc-generate)    bash "${SCRIPT_DIR}/sqlc-generate.sh" "$@" ;;
  service)          bash "${SCRIPT_DIR}/service-manager.sh" "$@" ;;
  *)
    echo -e "❌ 알 수 없는 명령어: ${YELLOW}${COMMAND}${NC}"
    echo -e "도움말: ./task.sh help"
    exit 1
    ;;
esac
