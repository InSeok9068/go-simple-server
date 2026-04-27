#!/bin/bash

set -e  # 스크립트 실행 중 오류 발생 시 중단

# 기본값 설정: 인자 미지정 시 build만 실행
RUN_DEFAULT=true
RUN_BUILD=false
RUN_TEST=false
RUN_LINT=false

# 파라미터 처리
for arg in "$@"; do
  case $arg in
    build) RUN_BUILD=true; RUN_DEFAULT=false ;;
    test) RUN_TEST=true; RUN_DEFAULT=false ;;
    lint) RUN_LINT=true; RUN_DEFAULT=false ;;
    *)
      echo "사용법: $0 [build|test|lint]"
      echo "  build: 프로젝트 빌드만 실행"
      echo "  test: 테스트만 실행"
      echo "  lint: Lint 검사만 실행"
      echo "  미지정: 프로젝트 빌드만 실행"
      exit 1
      ;;
  esac
done

# 기본 check는 빠른 빌드 오류 확인만 실행
if [ "$RUN_DEFAULT" = true ]; then
  RUN_BUILD=true
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
  golangci-lint run --new-from-rev=origin/main ./...
  echo "✅ 코드 스타일 검사 완료!"
  echo ""
fi

echo "🎉 선택한 검사가 완료되었습니다!"
