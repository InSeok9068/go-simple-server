#!/bin/bash

if [[ $# -lt 1 || $# -gt 2 ]]; then
  echo "Usage: $0 {window|linux} [-u]"
  exit 1
fi

OS="$1"
UPDATE="${2:-}"
FILE=""
EXE=""
TARGET_DIR="tailwind"
TARGET_FILE="$TARGET_DIR/tailwindcss$EXE"

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

# EXE 변수 설정 이후에 타깃 파일 경로 갱신
TARGET_FILE="$TARGET_DIR/tailwindcss$EXE"
UPDATE_FLAG=false
if [[ "$UPDATE" == "-u" ]]; then
  UPDATE_FLAG=true
fi

mkdir -p "$TARGET_DIR"

if [[ ! -f "$TARGET_FILE" || "$UPDATE_FLAG" == true ]]; then
  if [[ -f "$TARGET_FILE" && "$UPDATE_FLAG" == true ]]; then
    echo "기존 tailwindcss$EXE 파일을 최신 버전으로 갱신합니다."
  else
    echo "tailwindcss$EXE 파일이 없어 새로 설치합니다."
  fi
  echo "Downloading tailwindcss-$FILE..."
  curl -sLO "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$FILE"
  mv "tailwindcss-$FILE" "$TARGET_FILE"
  chmod a+x "$TARGET_FILE"
  echo "tailwindcss$EXE 설치가 완료되었습니다."
else
  echo "tailwindcss$EXE 파일이 이미 설치되어 있습니다. -u 옵션으로 최신 버전을 설치할 수 있습니다."
fi
