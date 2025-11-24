#!/bin/bash
set -euo pipefail

PORT="${1-}"

if [[ -z "${PORT}" ]]; then
  echo "❌ 포트 번호를 입력해주세요."
  echo "사용법: ./task.sh kill <port>"
  exit 1
fi

echo "🔍 포트 ${PORT} 점유 프로세스 찾는 중..."

# OS 탐지
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
  # Windows (Git Bash)
  # netstat -ano 출력 예: TCP    0.0.0.0:8080           0.0.0.0:0              LISTENING       1234
  # 5번째 컬럼이 PID
  PID=$(netstat -ano | grep ":${PORT} " | grep "LISTENING" | awk '{print $5}' | head -n 1 || true)

  if [[ -n "${PID}" ]]; then
    echo "✅ PID 발견: ${PID}"
    taskkill //F //PID "${PID}"
    echo "🗑️  프로세스(PID: ${PID})가 종료되었습니다."
  else
    echo "⚠️  포트 ${PORT}를 사용하는 프로세스가 없습니다."
  fi

else
  # Linux / macOS
  PID=$(lsof -t -i :${PORT} 2>/dev/null || true)

  if [[ -n "${PID}" ]]; then
    echo "✅ PID 발견: ${PID}"
    kill -9 ${PID}
    echo "🗑️  프로세스(PID: ${PID})가 종료되었습니다."
  else
    echo "⚠️  포트 ${PORT}를 사용하는 프로세스가 없습니다."
  fi
fi
