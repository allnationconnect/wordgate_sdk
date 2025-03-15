#!/bin/bash

# 确保脚本在错误时停止
set -e

echo "创建apt-repository分支..."

# 检查当前分支
CURRENT_BRANCH=$(git branch --show-current)
echo "当前分支: $CURRENT_BRANCH"

# 检查是否有未提交的更改
if [[ -n $(git status -s) ]]; then
  echo "错误: 有未提交的更改，请先提交所有更改。"
  exit 1
fi

# 检查是否已存在apt-repository分支
if git show-ref --verify --quiet refs/heads/apt-repository; then
  echo "apt-repository分支已存在，跳过创建步骤。"
else
  # 创建一个空的分支
  echo "创建全新的apt-repository分支..."
  git checkout --orphan apt-repository
  git rm -rf .
  
  # 创建初始结构
  mkdir -p apt-repo
  touch apt-repo/.gitkeep
  
  # 提交初始结构
  git add apt-repo
  git commit -m "初始化APT仓库结构"
  
  # 推送到远程
  git push -u origin apt-repository
  
  echo "apt-repository分支已创建并推送到远程。"
  
  # 切回原始分支
  git checkout "$CURRENT_BRANCH"
fi

echo "完成！"
echo "你现在可以按照RELEASE.md中的说明发布新版本。" 