#!/bin/bash
# Zpanel 一键安装脚本（适配 Debian 12 / Ubuntu）
# 用法:
#   curl -sSL https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh | bash
#   bash scripts/install.sh --port 8888 --username admin --password 'yourpass' --entry mypanel
#   bash scripts/install.sh --interactive
#   ZPANEL_INSTALL_LOCAL=1 bash scripts/install.sh --interactive

set -euo pipefail

ZPANEL_VERSION="${ZPANEL_VERSION:-latest}"
GITHUB_REPO="${GITHUB_REPO:-Cyruss-top/Zpanel}"
INSTALL_DIR="/usr/local/zpanel"
BIN_PATH="/usr/local/bin/zpanel"
CONFIG_DIR="/etc/zpanel"
DATA_DIR="/var/lib/zpanel"
LOG_DIR="/var/log/zpanel"
WWW_ROOT="/var/www"
DEFAULT_PORT=8888

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

[[ $EUID -ne 0 ]] && error "请使用 root 用户运行此脚本"

detect_os() {
    if [[ -f /etc/os-release ]]; then
        # shellcheck source=/dev/null
        . /etc/os-release
        case "$ID" in
            ubuntu|debian|centos|rocky|almalinux) info "检测到系统: $ID $VERSION_ID" ;;
            *) error "暂不支持 $ID" ;;
        esac
    else
        error "不支持的操作系统"
    fi
}

install_deps() {
    info "安装基础依赖..."
    if command -v apt-get &>/dev/null; then
        apt-get update -qq
        apt-get install -y -qq curl wget sqlite3 openssl ca-certificates
    elif command -v yum &>/dev/null; then
        yum install -y curl wget sqlite openssl ca-certificates
    fi
}

install_binary() {
    mkdir -p "$INSTALL_DIR/bin" "$INSTALL_DIR/scripts" "$INSTALL_DIR/templates"

    if [[ "${ZPANEL_INSTALL_LOCAL:-}" == "1" ]]; then
        SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
        ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
        LOCAL_BIN="$ROOT_DIR/bin/zpanel"
        [[ -f "$LOCAL_BIN" ]] || error "本地二进制不存在: $LOCAL_BIN (先执行 make build-all)"
        install -m 755 "$LOCAL_BIN" "$BIN_PATH"
        cp -r "$ROOT_DIR/scripts/"* "$INSTALL_DIR/scripts/" 2>/dev/null || true
        cp -r "$ROOT_DIR/templates/"* "$INSTALL_DIR/templates/" 2>/dev/null || true
        info "已安装本地二进制"
        return
    fi

    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)  ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        *) error "不支持的架构: $ARCH" ;;
    esac

    URL="https://github.com/${GITHUB_REPO}/releases/download/${ZPANEL_VERSION}/zpanel-linux-${ARCH}.tar.gz"
    info "下载 Zpanel ${ZPANEL_VERSION} (${ARCH})..."
    TMP=$(mktemp -d)
    if curl -fsSL "$URL" | tar xz -C "$TMP" 2>/dev/null; then
        install -m 755 "$TMP/zpanel" "$BIN_PATH"
        rm -rf "$TMP"
        info "二进制已安装: $BIN_PATH"
        return
    fi
    warn "远程下载失败，尝试使用本地 bin/zpanel"
    SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    LOCAL_BIN="$(cd "$SCRIPT_DIR/.." && pwd)/bin/zpanel"
    [[ -f "$LOCAL_BIN" ]] || error "下载失败且无本地二进制，请先 make build 或使用 ZPANEL_INSTALL_LOCAL=1"
    install -m 755 "$LOCAL_BIN" "$BIN_PATH"
}

prompt_config() {
    local FLAG_PORT=0 FLAG_USER=0 FLAG_PASS=0 FLAG_ENTRY=0
    [[ "$PORT" != "$DEFAULT_PORT" ]] && FLAG_PORT=1
    [[ "$USERNAME" != "admin" ]] && FLAG_USER=1
    [[ -n "$PASSWORD" ]] && FLAG_PASS=1
    [[ -n "$ENTRY" ]] && FLAG_ENTRY=1

    if [[ "${INTERACTIVE:-0}" == "1" ]] || [[ -t 0 && $FLAG_PORT -eq 0 && $FLAG_USER -eq 0 && $FLAG_PASS -eq 0 ]]; then
        echo ""
        echo "========== 自定义面板配置（直接回车使用默认值）=========="
        read -rp "面板端口 [${PORT}]: " input
        [[ -n "$input" ]] && PORT="$input"
        read -rp "管理员用户名 [${USERNAME}]: " input
        [[ -n "$input" ]] && USERNAME="$input"
        if [[ -z "$PASSWORD" ]]; then
            read -rsp "管理员密码（留空自动生成）: " input
            echo ""
            [[ -n "$input" ]] && PASSWORD="$input"
        fi
        read -rp "安全入口后缀（如 mypanel，留空不启用）: " input
        [[ -n "$input" ]] && ENTRY="$input"
        echo "======================================================"
        echo ""
    fi
}

init_config() {
    local PORT=${1:-$DEFAULT_PORT}
    local USERNAME=${2:-admin}
    local PASSWORD=${3:-}
    local ENTRY=${4:-}

    [[ -z "$PASSWORD" ]] && PASSWORD=$(openssl rand -base64 12)

    mkdir -p "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR" "$WWW_ROOT"

    local ENTRY_YAML=""
    if [[ -n "$ENTRY" ]]; then
        ENTRY_YAML="  entry: \"${ENTRY}\""
    fi

    cat > "$CONFIG_DIR/config.yaml" <<EOF
panel:
  port: ${PORT}
  bind: "0.0.0.0"
  ssl: false
${ENTRY_YAML}
auth:
  username: "${USERNAME}"
paths:
  www: "${WWW_ROOT}"
  data: "${DATA_DIR}"
  logs: "${LOG_DIR}"
  nginx_sites: "/etc/nginx/sites-available"
  nginx_enabled: "/etc/nginx/sites-enabled"
files:
  allowed_paths:
    - "${WWW_ROOT}"
  max_upload_size: 52428800
database:
  sqlite: "zpanel.db"
EOF

    export ZPANEL_CONFIG="$CONFIG_DIR/config.yaml"
    "$BIN_PATH" user password "$PASSWORD" >/dev/null 2>&1 || true
    printf '%s\n' "$PASSWORD"
}

setup_service() {
    cat > /etc/systemd/system/zpanel.service <<EOF
[Unit]
Description=Zpanel Linux Admin Panel
After=network.target

[Service]
Type=simple
ExecStart=${BIN_PATH} server
Restart=on-failure
RestartSec=5
LimitNOFILE=65535
Environment=ZPANEL_CONFIG=${CONFIG_DIR}/config.yaml

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
    systemctl enable zpanel
    systemctl restart zpanel
}

setup_firewall() {
    local PORT=$1
    if command -v ufw &>/dev/null && ufw status 2>/dev/null | grep -q "Status: active"; then
        ufw allow "$PORT/tcp"
        info "ufw 已开放端口 $PORT"
    fi
}

setup_symlink() {
    ln -sf "$BIN_PATH" /usr/bin/zp
}

print_info() {
    local PORT=$1 USERNAME=$2 PASSWORD=$3 ENTRY=$4
    local IP PATH_SUFFIX=""
    IP=$(curl -s --connect-timeout 3 ip.sb 2>/dev/null || hostname -I 2>/dev/null | awk '{print $1}')
    [[ -n "$ENTRY" ]] && PATH_SUFFIX="/${ENTRY}/"
    echo ""
    echo "============================================"
    echo -e "  ${GREEN}Zpanel 安装成功!${NC}"
    echo "============================================"
    echo -e "  面板地址: http://${IP}:${PORT}${PATH_SUFFIX}"
    echo -e "  用户名:   ${USERNAME}"
    echo -e "  密码:     ${PASSWORD}"
    [[ -n "$ENTRY" ]] && echo -e "  安全入口: /${ENTRY}/"
    echo -e "  管理命令: zpanel 或 zp"
    echo "============================================"
    echo -e "  ${YELLOW}请妥善保存以上信息${NC}"
    echo "============================================"
}

main() {
    local PORT=$DEFAULT_PORT USERNAME="admin" PASSWORD="" ENTRY="" INTERACTIVE=0

    while [[ $# -gt 0 ]]; do
        case "$1" in
            --port)       PORT="$2"; shift 2 ;;
            --username)   USERNAME="$2"; shift 2 ;;
            --password)   PASSWORD="$2"; shift 2 ;;
            --entry)      ENTRY="$2"; shift 2 ;;
            --interactive|-i) INTERACTIVE=1; shift ;;
            *) shift ;;
        esac
    done

    detect_os
    prompt_config
    install_deps
    install_binary
    setup_symlink
    PASSWORD=$(init_config "$PORT" "$USERNAME" "$PASSWORD" "$ENTRY")
    setup_service
    setup_firewall "$PORT"
    print_info "$PORT" "$USERNAME" "$PASSWORD" "$ENTRY"
}

main "$@"
