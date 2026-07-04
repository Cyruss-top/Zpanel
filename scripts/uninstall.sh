#!/bin/bash
# 卸载 Zpanel 面板（保留网站数据）
# 用法: bash scripts/uninstall.sh [yes]
set -euo pipefail

KEEP_DATA=${1:-yes}

systemctl stop zpanel 2>/dev/null || true
systemctl disable zpanel 2>/dev/null || true
rm -f /etc/systemd/system/zpanel.service
rm -f /usr/local/bin/zpanel /usr/bin/zp

if [[ "$KEEP_DATA" != "yes" ]]; then
    rm -rf /etc/zpanel /var/lib/zpanel /var/log/zpanel
fi

systemctl daemon-reload 2>/dev/null || true
echo "Zpanel 已卸载。网站目录 /var/www 未删除。"
