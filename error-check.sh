#!/bin/bash

set -e  # 스크립트 실행 중 오류 발생 시 중단

# 기본값 설정: 모든 작업 실행
RUN_ALL=true
RUN_BUILD=false
RUN_TEST=false
RUN_LINT=false

# 파라미터 처리
for arg in "$@"; do
  case $arg in
    build) RUN_BUILD=true; RUN_ALL=false ;;
    test) RUN_TEST=true; RUN_ALL=false ;;
    lint) RUN_LINT=true; RUN_ALL=false ;;
    *)
      echo "사용법: $0 [build|test|lint]"
      echo "  build: 프로젝트 빌드만 실행"
      echo "  test: 테스트만 실행"
      echo "  lint: Lint 검사만 실행"
      echo "  미지정: 모든 작업 실행"
      exit 1
      ;;
  esac
done

# 모든 작업 실행 옵션이 켜져있으면 모든 플래그를 true로 설정
if [ "$RUN_ALL" = true ]; then
  RUN_BUILD=true
  RUN_TEST=true
  RUN_LINT=true
fi

# 빌드 실행
if [ "$RUN_BUILD" = true ]; then
  echo "🔍 프로젝트 빌드 검사 중..."
  go build ./...
  echo "✅ 빌드 완료!"
  echo ""
fi

# 테스트 실행
if [ "$RUN_TEST" = true ]; then
  echo "🧪 테스트 실행 중..."
  go test ./...
  echo "✅ 모든 테스트 통과!"
  echo ""
fi

# Lint 실행
if [ "$RUN_LINT" = true ]; then
  echo "📝 코드 스타일 검사 중 (golangci-lint)..."
  golangci-lint run ./...
  echo "✅ 코드 스타일 검사 완료!"
  echo ""
fi

echo "🎉 모든 검사가 완료되었습니다!"
