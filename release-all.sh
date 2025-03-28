#!/bin/bash

# 현재 브랜치가 main이라고 가정 (아니면 체크아웃)
git checkout main || exit 1
git pull origin main || exit 1

# release/* 브랜치 목록 가져오기
branches=$(git branch -r | grep 'origin/release/' | sed 's/origin\///')

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
