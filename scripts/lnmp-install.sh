#!/bin/bash
# LNMP 一键安装 — 由 zpanel lnmp install 或面板 UI 调用
set -euo pipefail

install_ubuntu() {
    export DEBIAN_FRONTEND=noninteractive
    apt-get update -qq
    apt-get install -y -qq nginx mysql-server \
        php-fpm php-mysql php-cli php-curl php-gd php-mbstring php-xml php-zip unzip curl

    PHP_VER=$(php -r 'echo PHP_MAJOR_VERSION.".".PHP_MINOR_VERSION;')
    systemctl enable --now nginx mysql "php${PHP_VER}-fpm"
}

install_centos() {
    yum install -y nginx mysql-server \
        php-fpm php-mysqlnd php-cli php-curl php-gd php-mbstring php-xml php-zip
    systemctl enable --now nginx mysqld php-fpm
}

OS_ID=""
if [[ -f /etc/os-release ]]; then
    # shellcheck source=/dev/null
    . /etc/os-release
    OS_ID=$ID
fi

case "$OS_ID" in
    ubuntu|debian) install_ubuntu ;;
    centos|rocky|almalinux) install_centos ;;
    *) echo '{"ok":false,"message":"unsupported os"}'; exit 1 ;;
esac

NGINX_VER=$(nginx -v 2>&1 | awk -F/ '{print $2}')
PHP_VER=$(php -r 'echo PHP_MAJOR_VERSION.".".PHP_MINOR_VERSION;' 2>/dev/null || echo "unknown")
MYSQL_VER=$(mysql --version 2>/dev/null | awk '{print $3}' | tr -d ',' || echo "unknown")

echo "{\"ok\":true,\"nginx\":\"${NGINX_VER}\",\"php\":\"${PHP_VER}\",\"mysql\":\"${MYSQL_VER}\"}"
