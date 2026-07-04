#!/bin/bash
# 版本号 bump 脚本
# 用法:
#   ./scripts/bump-version.sh patch   # 0.1.0 -> 0.1.1
#   ./scripts/bump-version.sh minor   # 0.1.0 -> 0.2.0
#   ./scripts/bump-version.sh major   # 0.1.0 -> 1.0.0
#   ./scripts/bump-version.sh 1.2.3   # 指定版本

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION_FILE="$ROOT/VERSION"
CHANGELOG="$ROOT/CHANGELOG.md"

current=$(cat "$VERSION_FILE" | tr -d '[:space:]')
IFS='.' read -r major minor patch <<< "$current"

case "${1:-}" in
  major)
    major=$((major + 1)); minor=0; patch=0
    ;;
  minor)
    minor=$((minor + 1)); patch=0
    ;;
  patch)
    patch=$((patch + 1))
    ;;
  "")
    echo "用法: $0 {major|minor|patch|X.Y.Z}"
    exit 1
    ;;
  *)
    if [[ "$1" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
      IFS='.' read -r major minor patch <<< "$1"
    else
      echo "无效版本: $1"
      exit 1
    fi
    ;;
esac

new_version="${major}.${minor}.${patch}"
today=$(date +%Y-%m-%d)

echo "$new_version" > "$VERSION_FILE"
echo "VERSION -> $new_version"

# 提示更新 CHANGELOG（不自动写内容，避免覆盖手工记录）
echo ""
echo "请手动更新 CHANGELOG.md："
echo "  1. 将 [Unreleased] 下内容移到 [$new_version] - $today"
echo "  2. 更新底部链接"
echo ""
echo "然后执行发布："
echo "  git add VERSION CHANGELOG.md"
echo "  git commit -m \"chore: release v$new_version\""
echo "  git tag -a v$new_version -m \"Release v$new_version\""
echo "  git push origin main --tags"
