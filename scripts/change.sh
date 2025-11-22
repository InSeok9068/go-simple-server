#!/bin/bash

# 사용법 안내
if [[ $# -ne 1 ]]; then
  echo "Usage: $0 {homepage|ai-study|sample}"
  exit 1
fi

SERVICE="$1"
PORT="$2"
CONFIG_FILE=".air.toml"

# 서비스명에 따른 포트 자동 설정
case "$SERVICE" in
  homepage)
    PORT=8000
    ;;
  ai-study)
    PORT=8001
    ;;
  deario)
    PORT=8002
    ;;
  closet)
    PORT=8003
    ;;
  sample)
    PORT=8999
    ;;
  *)
    echo "❌ Unknown service: $SERVICE"
    exit 1
    ;;
esac

# `cmd` 라인을 정확히 한 번만 변경 (중복 방지)
sed -i'' -E 's|(cmd = "templ generate ; go build -o ./main.exe ./projects/)[^"]+(")|\1'"$SERVICE"'/cmd/\2|' "$CONFIG_FILE"
#sed -i'' -E 's|(cmd = "tailwind/tailwindcss.exe -i tailwind/tailwindcss.css -o shared/static/tailwindcss.css  --minify && go build -o ./main.exe ./cmd/)[^"]+(")|\1'"$SERVICE"'\2|' "$CONFIG_FILE"

# `app_port` 값을 정확히 변경
sed -i'' -E 's|^(app_port = )[0-9]+$|\1'"$PORT"'|' "$CONFIG_FILE"

echo "✅ air.toml 업데이트 완료 (cmd: $SERVICE, app_port: $PORT)"
