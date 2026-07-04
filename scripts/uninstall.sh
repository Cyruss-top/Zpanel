#!/bin/bash
# 卸载 Zpanel 面板
# 用法:
#   bash uninstall.sh                    # 交互选择卸载模式
#   bash uninstall.sh --keep-www --yes   # 保留 /var/www 网站数据
#   bash uninstall.sh --all --yes        # 彻底删除（含网站数据）
#   zpanel uninstall
#   zpanel uninstall -y --keep-www
#   zpanel uninstall -y --all

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

# keep-www | all
MODE=""
AUTO_YES=0

while [[ $# -gt 0 ]]; do
    case "$1" in
        --keep-www) MODE="keep-www"; shift ;;
        --all|--purge-all) MODE="all"; shift ;;
        --purge) MODE="all"; shift ;;  # 兼容旧参数
        --yes|-y) AUTO_YES=1; shift ;;
        yes) AUTO_YES=1; shift ;;
        *) shift ;;
    esac
done

read_choice() {
    local prompt=$1
    local value=""
    if [[ -t 0 ]]; then
        read -rp "$prompt" value
    elif [[ -r /dev/tty ]]; then
        read -rp "$prompt" value < /dev/tty
    fi
    printf '%s' "$value"
}

choose_mode() {
    [[ -n "$MODE" ]] && return 0
    echo ""
    echo "============================================"
    echo "  Zpanel 卸载"
    echo "============================================"
    echo "  1. 保留站点数据"
    echo "     删除面板程序、配置、数据库、日志"
    echo "     保留网站目录: ${WWW_ROOT}"
    echo ""
    echo "  2. 彻底删除干净"
    echo "     删除面板及 ${WWW_ROOT} 下全部网站文件"
    echo "     此操作不可恢复！"
    echo "============================================"
    local choice
    choice=$(read_choice "请选择 [1/2]: ")
    case "$choice" in
        1) MODE="keep-www" ;;
        2) MODE="all" ;;
        *) echo "已取消"; exit 0 ;;
    esac
}

confirm_mode() {
    [[ $AUTO_YES -eq 1 ]] && return 0
    echo ""
    local ans=""
    if [[ "$MODE" == "all" ]]; then
        warn "即将彻底删除 Zpanel 及 ${WWW_ROOT} 全部数据"
        ans=$(read_choice "确认彻底删除? 输入 yes: ")
        [[ "$ans" == "yes" ]]
    else
        echo "即将卸载 Zpanel，保留 ${WWW_ROOT} 网站数据"
        ans=$(read_choice "确认卸载? [y/N] ")
        [[ "$ans" == "y" || "$ans" == "Y" ]]
    fi
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

remove_panel() {
    info "删除 systemd 服务..."
    rm -f /etc/systemd/system/zpanel.service

    info "删除二进制..."
    rm -f "$BIN_PATH" /usr/bin/zp

    if [[ -d "$INSTALL_DIR" ]]; then
        info "删除程序目录 ${INSTALL_DIR}..."
        rm -rf "$INSTALL_DIR"
    fi

    info "删除配置与数据..."
    rm -rf "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"

    systemctl daemon-reload 2>/dev/null || true
}

remove_www() {
    if [[ -d "$WWW_ROOT" ]]; then
        warn "删除网站目录 ${WWW_ROOT}..."
        rm -rf "$WWW_ROOT"
    fi
}

print_done() {
    echo ""
    if [[ "$MODE" == "all" ]]; then
        echo -e "${GREEN}Zpanel 已彻底卸载，所有数据已删除${NC}"
    else
        echo -e "${GREEN}Zpanel 已卸载完成${NC}"
        echo "网站数据保留于: ${WWW_ROOT}"
    fi
}

main() {
    choose_mode
    confirm_mode || { echo "已取消"; exit 0; }
    stop_service
    remove_panel
    if [[ "$MODE" == "all" ]]; then
        remove_www
    fi
    print_done
}

main "$@"
