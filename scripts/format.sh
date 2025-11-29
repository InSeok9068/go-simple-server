#!/bin/bash
set -euo pipefail

BLUE='\033[0;34m'; GREEN='\033[0;32m'; NC='\033[0m'


# Function to run tailwind formatting
run_tailwind() {
  echo -e "${BLUE}Running tailwind sorter...${NC}"
  go-tailwind-sorter . --fix
}

# Function to run prettier
run_prettier() {
  echo -e "${BLUE}Running prettier...${NC}"
  npx -y prettier --write .
}

# Function to run go formatting
run_go() {
  echo -e "${BLUE}Running go fmt...${NC}"
  go fmt ./...
}

# Function to run templ formatting
run_templ() {
  echo -e "${BLUE}Running templ fmt...${NC}"
  templ fmt .
}

# Check for arguments
if [ $# -eq 0 ]; then
  # No arguments, run all
  run_tailwind
  run_prettier
  run_go
  run_templ
else
  case "$1" in
    tailwind)
      run_tailwind
      ;;
    prettier)
      run_prettier
      ;;
    go)
      run_go
      ;;
    templ)
      run_templ
      ;;
    *)
      echo -e "${BLUE}Usage: $0 [go|templ|tailwind|prettier]${NC}"
      exit 1
      ;;
  esac
fi

echo -e "${GREEN}Format complete!${NC}"
