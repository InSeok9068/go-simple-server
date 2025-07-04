#!/bin/bash

#sudo apt update
#sudo apt install -y build-essential gcc-multilib

# Linux 환경 변수 설정
export GOOS=linux
export GOARCH=amd64
# export CGO_ENABLED=1

# 빌드 실행
go build -o ./main ./cmd/homepage

# 결과 출력
echo "Linux용 바이너리가 빌드되었습니다"