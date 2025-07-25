#!/bin/bash

# notepad ~/.bash_profile
# source "scripts/task-completion.sh"

# Git Bash for Windows용 task.sh 자동 완성 스크립트

_task_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # 메인 명령어 목록
    # ./task.sh [TAB] 을 누를 때 (두 번째 단어 완성)
    if [ ${COMP_CWORD} -eq 1 ]; then
        opts="help switch check deps build-linux release install-tailwind sqlc-generate service"
        COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
        return 0
    fi

    # 하위 명령어 목록 (이전 단어를 기반으로)
    case "${prev}" in
        switch|release|sqlc-generate)
            # projects/ 디렉터리 아래의 서비스 이름을 동적으로 찾아 제안
            local project_dirs
            # 프로젝트 루트에서 실행한다고 가정
            if [ -d "projects" ]; then
                project_dirs=$(find projects -maxdepth 1 -mindepth 1 -type d -exec basename {} \;)
                COMPREPLY=( $(compgen -W "${project_dirs}" -- "${cur}") )
            fi
            ;;
        check)
            opts="build test lint"
            COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
            ;;
        install-tailwind)
            opts="win linux"
            COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
            ;;
        service)
            opts="create deploy remove undeploy"
            COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
            ;;
        *)
            # 다른 명령어에 대한 기본 완성은 없음
            COMPREPLY=()
            ;;
    esac
}

# 'task.sh'와 './task.sh' 명령어에 대해 _task_completion 함수를 등록
complete -F _task_completion task.sh
complete -F _task_completion ./task.sh
