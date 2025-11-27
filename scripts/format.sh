#!/bin/bash
set -euo pipefail

BLUE='\033[0;34m'; GREEN='\033[0;32m'; NC='\033[0m'

echo -e "${BLUE}Running go fmt...${NC}"
go fmt ./...

echo -e "${BLUE}Running prettier...${NC}"
npx -y prettier --write .

echo -e "${BLUE}Running tailwind sorter...${NC}"
go-tailwind-sorter . --fix

echo -e "${BLUE}Running templ fmt...${NC}"
templ fmt .

echo -e "${GREEN}Format complete!${NC}"
