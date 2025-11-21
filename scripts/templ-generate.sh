#!/bin/bash
set -euo pipefail

BLUE='\033[0;34m'; GREEN='\033[0;32m'; NC='\033[0m'

echo -e "${BLUE}Running templ generate...${NC}"
templ generate

echo -e "${GREEN}Templ generation complete!${NC}"
