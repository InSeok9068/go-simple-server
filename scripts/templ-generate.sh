#!/bin/bash
set -euo pipefail

BLUE='\033[0;34m'; GREEN='\033[0;32m'; NC='\033[0m'

echo -e "${BLUE}프로젝트에 고정된 Templ 버전으로 템플릿 생성 중...${NC}"
go tool templ generate

echo -e "${GREEN}Templ 생성 완료!${NC}"
