#!/bin/bash

set -e  # 스크립트 실행 중 오류 발생 시 중단

echo "🔍 프로젝트 빌드 검사 중..."
go build ./...

echo "✅ 빌드 완료!"

# echo "🧪 테스트 실행 중..."
# go test ./...

# echo "✅ 모든 테스트 통과!"

echo "📝 코드 스타일 검사 중 (golangci-lint)..."
golangci-lint run ./...

echo "✅ 코드 스타일 검사 완료!"

echo "🎉 모든 검사가 완료되었습니다!"
