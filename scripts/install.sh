#!/bin/bash
# Zpanel 一键安装脚本（适配 Debian 12 / Ubuntu）
# 用法:
#   # 国际节点 (GitHub)
#   curl -sSL https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh | bash -s -- --interactive
#   # 中国大陆节点 (Gitee，推荐)
#   wget -qO- https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- --interactive
#   bash scripts/install.sh --port 8888 --username admin --password 'yourpass' --entry mypanel
#   bash scripts/install.sh --mirror gitee --interactive
#   bash scripts/install.sh --package /root/zpanel-linux-amd64.tar.gz --interactive
#   bash install.sh --base-url https://your-server.com/zpanel --interactive
#   ZPANEL_BASE_URL=https://your-server.com/zpanel bash install.sh --interactive

set -euo pipefail

ZPANEL_VERSION="${ZPANEL_VERSION:-latest}"
ZPANEL_MIRROR="${ZPANEL_MIRROR:-github}"   # github | gitee
ZPANEL_BASE_URL="${ZPANEL_BASE_URL:-}"     # 自定义安装包地址，如 https://your-server.com/zpanel
GITHUB_REPO="${GITHUB_REPO:-Cyruss-top/Zpanel}"
GITEE_REPO="${GITEE_REPO:-Ressss2023/Zpanel}"
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

fetch_url() {
    local url=$1
    if command -v curl &>/dev/null; then
        curl -fsSL "$url"
    elif command -v wget &>/dev/null; then
        wget -qO- "$url"
    else
        error "需要 curl 或 wget，请执行: apt install -y curl wget"
    fi
}

# 交互输入：wget | bash 管道执行时 stdin 非终端，需从 /dev/tty 读取
read_prompt() {
    local var=$1 prompt=$2 secret=${3:-0} value=""
    if [[ -t 0 ]]; then
        if [[ "$secret" == "1" ]]; then
            read -rsp "$prompt" value
            echo ""
        else
            read -rp "$prompt" value
        fi
    elif [[ -r /dev/tty ]]; then
        if [[ "$secret" == "1" ]]; then
            read -rsp "$prompt" value < /dev/tty
            echo "" > /dev/tty
        else
            read -rp "$prompt" value < /dev/tty
        fi
    else
        warn "无法交互输入，请使用命令行参数，或先 wget 下载脚本再 bash install.sh"
        return 1
    fi
    printf -v "$var" '%s' "$value"
}

release_urls() {
    local arch=$1
    if [[ -n "${ZPANEL_BASE_URL:-}" ]]; then
        local base="${ZPANEL_BASE_URL%/}"
        echo "${base}/zpanel-linux-${arch}.tar.gz"
        return
    fi
    local gh gt
    if [[ "$ZPANEL_VERSION" == "latest" ]]; then
        gh="https://github.com/${GITHUB_REPO}/releases/latest/download/zpanel-linux-${arch}.tar.gz"
        gt="https://gitee.com/${GITEE_REPO}/releases/latest/download/zpanel-linux-${arch}.tar.gz"
    else
        gh="https://github.com/${GITHUB_REPO}/releases/download/${ZPANEL_VERSION}/zpanel-linux-${arch}.tar.gz"
        gt="https://gitee.com/${GITEE_REPO}/releases/download/${ZPANEL_VERSION}/zpanel-linux-${arch}.tar.gz"
    fi
    if [[ "$ZPANEL_MIRROR" == "gitee" ]]; then
        echo "$gt"
        echo "$gh"
    else
        echo "$gh"
        echo "$gt"
    fi
}

install_binary() {
    mkdir -p "$INSTALL_DIR/bin" "$INSTALL_DIR/scripts" "$INSTALL_DIR/templates"

    if [[ -n "${ZPANEL_PACKAGE:-}" ]]; then
        [[ -f "$ZPANEL_PACKAGE" ]] || error "安装包不存在: $ZPANEL_PACKAGE"
        info "使用本地安装包: $ZPANEL_PACKAGE"
        TMP=$(mktemp -d)
        tar xzf "$ZPANEL_PACKAGE" -C "$TMP"
        ARCH=$(uname -m)
        case "$ARCH" in
            x86_64)  ARCH="amd64" ;;
            aarch64) ARCH="arm64" ;;
        esac
        if [[ -f "$TMP/zpanel" ]]; then
            install -m 755 "$TMP/zpanel" "$BIN_PATH"
        elif [[ -f "$TMP/zpanel-linux-${ARCH}" ]]; then
            install -m 755 "$TMP/zpanel-linux-${ARCH}" "$BIN_PATH"
        else
            rm -rf "$TMP"
            error "压缩包内未找到 zpanel 二进制"
        fi
        rm -rf "$TMP"
        info "二进制已安装: $BIN_PATH"
        return
    fi

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

    info "下载 Zpanel ${ZPANEL_VERSION} (${ARCH})${ZPANEL_BASE_URL:+, 源: ${ZPANEL_BASE_URL}}..."
    TMP=$(mktemp -d)
    local downloaded=0 url
    while IFS= read -r url; do
        [[ -z "$url" ]] && continue
        info "尝试: $url"
        if fetch_url "$url" | tar xz -C "$TMP" 2>/dev/null; then
            downloaded=1
            break
        fi
        warn "下载失败: $url"
    done < <(release_urls "$ARCH")

    if [[ $downloaded -eq 1 ]]; then
        if [[ -f "$TMP/zpanel" ]]; then
            install -m 755 "$TMP/zpanel" "$BIN_PATH"
        elif [[ -f "$TMP/zpanel-linux-${ARCH}" ]]; then
            install -m 755 "$TMP/zpanel-linux-${ARCH}" "$BIN_PATH"
        else
            error "压缩包内未找到 zpanel 二进制"
        fi
        rm -rf "$TMP"
        info "二进制已安装: $BIN_PATH"
        return
    fi
    rm -rf "$TMP"
    warn "远程下载失败，尝试使用本地 bin/zpanel"
    SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    LOCAL_BIN="$(cd "$SCRIPT_DIR/.." && pwd)/bin/zpanel"
    [[ -f "$LOCAL_BIN" ]] || error "下载失败且无本地二进制，请先 make build 或使用 ZPANEL_INSTALL_LOCAL=1"
    install -m 755 "$LOCAL_BIN" "$BIN_PATH"
}

install_scripts() {
    mkdir -p "$INSTALL_DIR/scripts"
    local dst="$INSTALL_DIR/scripts/uninstall.sh"
    local src=""

    if [[ -n "${BASH_SOURCE[0]:-}" ]] && [[ -f "$(dirname "${BASH_SOURCE[0]}")/uninstall.sh" ]]; then
        src="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/uninstall.sh"
    fi

    if [[ -n "$src" && -f "$src" ]]; then
        cp "$src" "$dst"
        chmod 755 "$dst"
        info "已安装卸载脚本: $dst"
        return
    fi

    if [[ -n "${ZPANEL_BASE_URL:-}" ]]; then
        if fetch_url "${ZPANEL_BASE_URL%/}/uninstall.sh" > "$dst" 2>/dev/null; then
            chmod 755 "$dst"
            info "已安装卸载脚本: $dst"
            return
        fi
    fi

    local url
    for url in \
        "https://gitee.com/${GITEE_REPO}/raw/main/scripts/uninstall.sh" \
        "https://raw.githubusercontent.com/${GITHUB_REPO}/main/scripts/uninstall.sh"; do
        if fetch_url "$url" > "$dst" 2>/dev/null; then
            sed -i 's/\r$//' "$dst" 2>/dev/null || true
            chmod 755 "$dst"
            info "已安装卸载脚本: $dst"
            return
        fi
    done
    warn "未能下载 uninstall.sh，仍可使用 zpanel uninstall 命令卸载"
}

prompt_config() {
    local FLAG_PORT=0 FLAG_USER=0 FLAG_PASS=0 FLAG_ENTRY=0
    [[ "$PORT" != "$DEFAULT_PORT" ]] && FLAG_PORT=1
    [[ "$USERNAME" != "admin" ]] && FLAG_USER=1
    [[ -n "$PASSWORD" ]] && FLAG_PASS=1
    [[ -n "$ENTRY" ]] && FLAG_ENTRY=1

    if [[ "${INTERACTIVE:-0}" == "1" ]] || [[ -t 0 && $FLAG_PORT -eq 0 && $FLAG_USER -eq 0 && $FLAG_PASS -eq 0 ]]; then
        local input=""
        echo ""
        echo "========== 自定义面板配置（直接回车使用默认值）=========="
        read_prompt input "面板端口 [${PORT}]: " || return 0
        [[ -n "$input" ]] && PORT="$input"
        read_prompt input "管理员用户名 [${USERNAME}]: " || return 0
        [[ -n "$input" ]] && USERNAME="$input"
        if [[ -z "$PASSWORD" ]]; then
            read_prompt input "管理员密码（留空自动生成）: " 1 || return 0
            [[ -n "$input" ]] && PASSWORD="$input"
        fi
        read_prompt input "安全入口后缀（如 mypanel，留空不启用）: " || return 0
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

verify_service() {
    local PORT=$1 ENTRY=${2:-}
    local PATH_SUFFIX=""
    [[ -n "$ENTRY" ]] && PATH_SUFFIX="/${ENTRY}/"
    sleep 1
    if ! systemctl is-active --quiet zpanel 2>/dev/null; then
        warn "zpanel 服务未运行，请执行: journalctl -u zpanel -n 30"
        return 1
    fi
    if ! ss -tlnp 2>/dev/null | grep -q ":${PORT} "; then
        warn "端口 ${PORT} 未监听，请检查配置与服务日志"
        return 1
    fi
    if curl -fsS --connect-timeout 3 "http://127.0.0.1:${PORT}${PATH_SUFFIX}" >/dev/null 2>&1; then
        info "本地访问正常 (127.0.0.1:${PORT}${PATH_SUFFIX})"
        return 0
    fi
    warn "本地访问失败，请执行: journalctl -u zpanel -n 30"
    return 1
}

setup_symlink() {
    ln -sf "$BIN_PATH" /usr/bin/zp
}

is_private_ip() {
    local ip=$1
    [[ "$ip" =~ ^10\. ]] && return 0
    [[ "$ip" =~ ^192\.168\. ]] && return 0
    [[ "$ip" =~ ^172\.(1[6-9]|2[0-9]|3[0-1])\. ]] && return 0
    [[ "$ip" =~ ^127\. ]] && return 0
    return 1
}

detect_public_ip() {
    local ip url
    for url in \
        "https://api.ipify.org" \
        "https://ifconfig.me/ip" \
        "https://ip.sb" \
        "https://icanhazip.com" \
        "https://myip.ipip.net"; do
        ip=$(curl -fsSL --connect-timeout 3 "$url" 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' | head -1)
        if [[ -n "$ip" ]] && ! is_private_ip "$ip"; then
            echo "$ip"
            return 0
        fi
    done
    return 1
}

detect_local_ip() {
    local ip
    for ip in $(hostname -I 2>/dev/null); do
        [[ -n "$ip" ]] && { echo "$ip"; return 0; }
    done
    echo "127.0.0.1"
}

print_info() {
    local PORT=$1 USERNAME=$2 PASSWORD=$3 ENTRY=$4
    local GOT_PUBLIC="" PUBLIC_IP LOCAL_IP PATH_SUFFIX="" LABEL="面板地址"
    GOT_PUBLIC=$(detect_public_ip || true)
    LOCAL_IP=$(detect_local_ip)
    if [[ -n "$GOT_PUBLIC" ]]; then
        PUBLIC_IP="$GOT_PUBLIC"
        LABEL="面板地址(公网)"
    else
        PUBLIC_IP="$LOCAL_IP"
    fi
    [[ -n "$ENTRY" ]] && PATH_SUFFIX="/${ENTRY}/"
    echo ""
    echo "============================================"
    echo -e "  ${GREEN}Zpanel 安装成功!${NC}"
    echo "============================================"
    echo -e "  ${LABEL}: http://${PUBLIC_IP}:${PORT}${PATH_SUFFIX}"
    if [[ -n "$GOT_PUBLIC" && "$LOCAL_IP" != "$GOT_PUBLIC" ]]; then
        echo -e "  面板地址(内网): http://${LOCAL_IP}:${PORT}${PATH_SUFFIX}"
    fi
    echo -e "  用户名:   ${USERNAME}"
    echo -e "  密码:     ${PASSWORD}"
    [[ -n "$ENTRY" ]] && echo -e "  安全入口: /${ENTRY}/"
    echo -e "  管理命令: zpanel 或 zp"
    echo "============================================"
    echo -e "  ${YELLOW}外网无法访问? 请到云厂商控制台安全组放行 TCP ${PORT}${NC}"
    echo -e "  ${YELLOW}服务器自检: systemctl status zpanel && ss -tlnp | grep ${PORT}${NC}"
    echo "============================================"
    echo -e "  ${YELLOW}请妥善保存以上信息${NC}"
    echo "============================================"
}

install_scripts() {
    mkdir -p "$INSTALL_DIR/scripts"
    local dst="$INSTALL_DIR/scripts/uninstall.sh"
    local src=""

    if [[ -n "${BASH_SOURCE[0]:-}" ]] && [[ -f "$(dirname "${BASH_SOURCE[0]}")/uninstall.sh" ]]; then
        src="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/uninstall.sh"
    fi

    if [[ -n "$src" && -f "$src" ]]; then
        cp "$src" "$dst"
        chmod 755 "$dst"
        info "已安装卸载脚本: $dst"
        return
    fi

    if [[ -n "${ZPANEL_BASE_URL:-}" ]]; then
        if fetch_url "${ZPANEL_BASE_URL%/}/uninstall.sh" > "$dst" 2>/dev/null; then
            sed -i 's/\r$//' "$dst" 2>/dev/null || true
            chmod 755 "$dst"
            info "已安装卸载脚本: $dst"
            return
        fi
    fi

    local url
    for url in \
        "https://gitee.com/${GITEE_REPO}/raw/main/scripts/uninstall.sh" \
        "https://raw.githubusercontent.com/${GITHUB_REPO}/main/scripts/uninstall.sh"; do
        if fetch_url "$url" > "$dst" 2>/dev/null; then
            sed -i 's/\r$//' "$dst" 2>/dev/null || true
            chmod 755 "$dst"
            info "已安装卸载脚本: $dst"
            return
        fi
    done
    warn "未能下载 uninstall.sh，仍可使用 zpanel uninstall 命令卸载"
}

main() {
    local PORT=$DEFAULT_PORT USERNAME="admin" PASSWORD="" ENTRY="" INTERACTIVE=0

    while [[ $# -gt 0 ]]; do
        case "$1" in
            --port)       PORT="$2"; shift 2 ;;
            --username)   USERNAME="$2"; shift 2 ;;
            --password)   PASSWORD="$2"; shift 2 ;;
            --entry)      ENTRY="$2"; shift 2 ;;
            --mirror)     ZPANEL_MIRROR="$2"; shift 2 ;;
            --package)    ZPANEL_PACKAGE="$2"; shift 2 ;;
            --base-url)   ZPANEL_BASE_URL="$2"; shift 2 ;;
            --interactive|-i) INTERACTIVE=1; shift ;;
            *) shift ;;
        esac
    done

    detect_os
    prompt_config
    install_deps
    install_binary
    install_scripts
    setup_symlink
    PASSWORD=$(init_config "$PORT" "$USERNAME" "$PASSWORD" "$ENTRY")
    setup_service
    setup_firewall "$PORT"
    verify_service "$PORT" "$ENTRY" || true
    print_info "$PORT" "$USERNAME" "$PASSWORD" "$ENTRY"
}

main "$@"
