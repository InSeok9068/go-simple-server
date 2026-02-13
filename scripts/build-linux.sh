#!/bin/bash
set -euo pipefail

list_projects() {
  local root projects_dir
  if root="$(git rev-parse --show-toplevel 2>/dev/null)"; then
    projects_dir="${root}/projects"
  else
    projects_dir="projects"
  fi

  shopt -s nullglob
  local dirs=("${projects_dir}"/*/)
  shopt -u nullglob

  local out=()
  local d
  for d in "${dirs[@]}"; do
    d="${d%/}"
    d="${d##*/}"
    out+=("${d}")
  done
  printf "%s\n" "${out[@]}"
}

select_project_interactive() {
  local projects=()
  while IFS= read -r p; do
    [[ -n "${p}" ]] && projects+=("${p}")
  done < <(list_projects)

  if [[ ${#projects[@]} -eq 0 ]]; then
    echo "❌ projects 폴더에서 프로젝트를 찾지 못했습니다."
    exit 1
  fi

  if [[ ! -t 0 ]]; then
    echo "❌ 비대화형 환경에서는 프로젝트 인자를 지정해야 합니다."
    echo "사용법: ./task.sh build-linux [project]"
    exit 1
  fi

  echo "빌드할 프로젝트를 선택하세요."
  local selected=""
  select candidate in "${projects[@]}"; do
    if [[ -n "${candidate:-}" ]]; then
      selected="${candidate}"
      break
    fi
    echo "유효한 번호를 선택하세요."
  done

  printf "%s\n" "${selected}"
}

PROJECT="${1:-}"
if [[ -z "${PROJECT}" ]]; then
  PROJECT="$(select_project_interactive)"
fi

BUILD_ENV="${BUILD_ENV:-prod}"

TARGET_PATH="./projects/${PROJECT}/cmd"
OUTPUT_PATH="./${PROJECT}"

if [[ ! -d "${TARGET_PATH}" ]]; then
  echo "❌ 프로젝트 실행 경로를 찾을 수 없습니다: ${TARGET_PATH}"
  echo "사용법: ./task.sh build-linux [project]"
  exit 1
fi

# 현재 셸에서만 ENV를 임시 오버라이드한다.
ORIGINAL_ENV="${ENV-}"
export ENV="${BUILD_ENV}"
trap 'export ENV="${ORIGINAL_ENV}"' EXIT

export GOOS=linux
export GOARCH=amd64

echo "🔧 Linux 빌드 시작"
echo "- 프로젝트: ${PROJECT}"
echo "- 출력 파일: ${OUTPUT_PATH}"
echo "- 임시 ENV: ${ENV}"

go build -ldflags "-s -w" -o "${OUTPUT_PATH}" "${TARGET_PATH}"

echo "✅ Linux용 바이너리 빌드 완료: ${OUTPUT_PATH}"
