#!/bin/bash
# 卸载 Zpanel 面板
# 用法:
#   zpanel uninstall              # 交互确认，保留配置数据
#   zpanel uninstall --yes        # 跳过确认
#   zpanel uninstall --yes --purge  # 同时删除配置与数据库
#   bash /usr/local/zpanel/scripts/uninstall.sh
#   bash scripts/uninstall.sh --yes
#   wget -qO- https://www.mczybh.cn/zpanel-release/uninstall.sh | bash -s -- --yes

set -euo pipefail

INSTALL_DIR="/usr/local/zpanel"
BIN_PATH="/usr/local/bin/zpanel"
CONFIG_DIR="/etc/zpanel"
DATA_DIR="/var/lib/zpanel"
LOG_DIR="/var/log/zpanel"
WWW_ROOT="/var/www"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

[[ $EUID -ne 0 ]] && error "请使用 root 用户运行此脚本"

AUTO_YES=0
PURGE=0

while [[ $# -gt 0 ]]; do
    case "$1" in
        --yes|-y)   AUTO_YES=1; shift ;;
        --purge)    PURGE=1; shift ;;
        yes)        AUTO_YES=1; shift ;;  # 兼容旧用法
        no)         shift ;;             # 兼容旧用法
        *) shift ;;
    esac
done

confirm() {
    [[ $AUTO_YES -eq 1 ]] && return 0
    echo ""
    echo "============================================"
    echo "  即将卸载 Zpanel 面板"
    echo "============================================"
    echo "  将删除: 服务、二进制、面板程序目录"
    if [[ $PURGE -eq 1 ]]; then
        echo "  将删除: 配置、数据库、日志 (${CONFIG_DIR} 等)"
    else
        echo "  将保留: 配置与数据 (${CONFIG_DIR}, ${DATA_DIR})"
    fi
    echo "  不会删除: 网站目录 ${WWW_ROOT}"
    echo "============================================"
    read -rp "确认卸载? [y/N] " ans
    [[ "$ans" == "y" || "$ans" == "Y" ]]
}

stop_service() {
    if systemctl is-active --quiet zpanel 2>/dev/null; then
        info "停止 zpanel 服务..."
        systemctl stop zpanel
    fi
    if systemctl is-enabled --quiet zpanel 2>/dev/null; then
        systemctl disable zpanel
    fi
}

remove_files() {
    info "删除 systemd 服务..."
    rm -f /etc/systemd/system/zpanel.service

    info "删除二进制..."
    rm -f "$BIN_PATH" /usr/bin/zp

    if [[ -d "$INSTALL_DIR" ]]; then
        info "删除程序目录 ${INSTALL_DIR}..."
        rm -rf "$INSTALL_DIR"
    fi

    if [[ $PURGE -eq 1 ]]; then
        info "删除配置与数据..."
        rm -rf "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"
    fi

    systemctl daemon-reload 2>/dev/null || true
}

main() {
    confirm || { echo "已取消"; exit 0; }
    stop_service
    remove_files
    echo ""
    echo -e "${GREEN}Zpanel 已卸载完成${NC}"
    if [[ $PURGE -eq 0 ]]; then
        echo "配置保留于: ${CONFIG_DIR}"
    fi
    echo "网站数据保留于: ${WWW_ROOT}"
}

main "$@"
