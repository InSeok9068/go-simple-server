#!/bin/bash

# 사용법 표시
show_usage() {
  echo "사용법: $0 [프로젝트명]"
  echo "  프로젝트명이 지정되지 않으면 모든 프로젝트를 릴리스합니다."
  echo "  예: $0            # 모든 프로젝트 릴리스"
  echo "  예: $0 deario     # deario 프로젝트만 릴리스"
  exit 1
}

# 현재 브랜치가 main이라고 가정 (아니면 체크아웃)
git checkout main || exit 1
git pull origin main || exit 1

# release/* 브랜치 목록 가져오기
branches=$(git branch -r | grep 'origin/release/' | sed 's/origin\///')

# 특정 프로젝트만 필터링
if [ $# -eq 1 ]; then
  project=$1
  branches=$(echo "$branches" | grep "release/$project")
  
  if [ -z "$branches" ]; then
    echo "❌ 잘못된 프로젝트명: $project"
    show_usage
  fi
fi

for branch in $branches; do
  echo "🔄 병합 시작: main → $branch"

  git checkout "$branch" || continue
  git pull origin "$branch" || continue

  # main 병합
  if git merge --no-edit main; then
    echo "✅ 머지 성공, push 중: $branch"
    git push origin "$branch"
  else
    echo "❌ 병합 충돌 발생: $branch"
    git merge --abort
  fi
done

# 마지막에 main으로 돌아오기
git checkout main
