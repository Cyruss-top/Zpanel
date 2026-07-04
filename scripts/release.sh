#!/bin/bash
# 发布前检查 + 打 tag
# 用法: ./scripts/release.sh [version]
# 示例: ./scripts/release.sh          # 使用 VERSION 文件中的版本
#       ./scripts/release.sh 0.2.0    # 指定版本

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
ok()   { echo -e "${GREEN}[OK]${NC} $*"; }
fail() { echo -e "${RED}[FAIL]${NC} $*"; exit 1; }
warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }

VERSION="${1:-$(cat VERSION | tr -d '[:space:]')}"
TAG="v${VERSION}"

echo "=== Zpanel Release v${VERSION} ==="

# 1. 工作区必须干净（VERSION/CHANGELOG 除外若正在发布）
if [[ -n "$(git status --porcelain | grep -v '^.. VERSION$\|^.. CHANGELOG.md$')" ]]; then
  git status --short
  fail "工作区有未提交改动，请先 commit"
fi

# 2. VERSION 文件一致
file_ver=$(cat VERSION | tr -d '[:space:]')
[[ "$file_ver" == "$VERSION" ]] || fail "VERSION 文件 ($file_ver) 与参数 ($VERSION) 不一致"

# 3. CHANGELOG 必须包含该版本
grep -q "\[${VERSION}\]" CHANGELOG.md || fail "CHANGELOG.md 缺少 [${VERSION}] 条目"

# 4. Go 测试（有代码后启用）
if [[ -f go.mod ]]; then
  go test ./... || fail "go test 失败"
  ok "go test 通过"
else
  warn "尚无 go.mod，跳过 Go 测试"
fi

# 5. 前端构建（有前端后启用）
if [[ -f web/package.json ]]; then
  (cd web && npm run build) || fail "前端构建失败"
  ok "前端构建通过"
else
  warn "尚无前端，跳过 npm build"
fi

# 6. 打 tag
if git rev-parse "$TAG" >/dev/null 2>&1; then
  fail "tag $TAG 已存在"
fi

echo ""
echo "即将创建 tag: $TAG"
read -rp "确认? [y/N] " confirm
[[ "$confirm" == "y" || "$confirm" == "Y" ]] || exit 0

git tag -a "$TAG" -m "Release $TAG"
ok "已创建 tag $TAG"

echo ""
echo "下一步:"
echo "  git push origin main"
echo "  git push origin $TAG"
