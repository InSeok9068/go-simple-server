#!/bin/bash

# notepad ~/.bash_profile
## task.sh 자동 완성 스크립트 불러오기
#if [ -f "scripts/task-completion.sh" ]; then
#  source "scripts/task-completion.sh"
#fi

_task_projects_list() {
  # projects 디렉토리 기준: find 대신 glob 사용 (Windows Git Bash에서 훨씬 빠름)
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
  for d in "${dirs[@]}"; do
    # 디렉토리명만
    d="${d%/}"; d="${d##*/}"
    out+=("$d")
  done
  printf "%s\n" "${out[@]}"
}

_task_completion() {
  local cur prev prev2 cmd
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  prev2="${COMP_WORDS[COMP_CWORD-2]:-}"
  cmd="${COMP_WORDS[1]:-}"

  # 1단계: 메인 명령어
  if [[ ${COMP_CWORD} -eq 1 ]]; then
    local opts="help switch check deps build-linux release install-tailwind sqlc-generate templ-generate fmt kill"
    COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
    return 0
  fi

  # 공통 프로젝트 후보
  local projects; projects="$(_task_projects_list 2>/dev/null)"

  # 2단계: 하위 문맥
  case "${cmd}" in
    switch|release|sqlc-generate)
      COMPREPLY=( $(compgen -W "${projects}" -- "${cur}") )
      ;;

    check)
      COMPREPLY=( $(compgen -W "build test lint" -- "${cur}") )
      ;;

    install-tailwind)
      COMPREPLY=( $(compgen -W "win linux" -- "${cur}") )   # ← win 으로 수정
      ;;

    fmt)
      COMPREPLY=( $(compgen -W "go templ tailwind prettier" -- "${cur}") )
      ;;
    *)
      COMPREPLY=()
      ;;
  esac
}

# 등록
complete -F _task_completion task.sh
complete -F _task_completion ./task.sh
