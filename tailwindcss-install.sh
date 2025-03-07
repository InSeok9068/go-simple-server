#!/bin/bash

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 {window|linux}"
  exit 1
fi

OS="$1"
FILE=""
EXE=""

case "$OS" in
  window)
    FILE="windows-x64.exe"
    EXE=".exe"
    ;;
  linux)
    FILE="linux-x64"
    ;;
  *)
    echo "❌ Unknown OS: $OS"
    exit 1
    ;;
esac

if [ ! -f "tailwindcss$EXE" ]; then
  echo "Downloading tailwindcss-$FILE..."
  curl -sLO "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$FILE"
  mv "tailwindcss-$FILE" "tailwindcss$EXE"
  chmod a+x "tailwindcss$EXE"
  mkdir -p node_modules/tailwindcss/lib
  echo '{"devDependencies": {"tailwindcss": "latest"}}' >package.json
  echo ""
else
  echo "tailwindcss$EXE already exists. Skipping download."
fi