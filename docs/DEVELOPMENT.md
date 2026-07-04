# Zpanel 开发文档

> 超级轻量版 Linux 可视化运维面板 — 对标宝塔核心能力，Go 单二进制部署，支持 LNMP 安装与 PHP / HTML / Go 项目管理。

**版本：** 见根目录 [VERSION](VERSION)（当前 v0.1.0）  
**最后更新：** 2026-07-04

---

## 目录

1. [项目概述](#1-项目概述)
2. [功能规划](#2-功能规划)
3. [技术架构](#3-技术架构)
4. [前端开发规范](#4-前端开发规范)
5. [后端开发规范](#5-后端开发规范)
6. [CLI 命令行（zpanel / zp）](#6-cli-命令行zpanel--zp)
7. [一键部署脚本](#7-一键部署脚本)
8. [UI 库与视觉设计](#8-ui-库与视觉设计)
9. [移动端适配策略](#9-移动端适配策略)
10. [目录结构](#10-目录结构)
11. [开发环境搭建](#11-开发环境搭建)
12. [安全规范](#12-安全规范)
13. [版本标记与 Git 备份](#13-版本标记与-git-备份)

---

## 1. 项目概述

### 1.1 定位

| 对比项 | 宝塔面板 | Zpanel |
|--------|----------|--------|
| 运行时 | PHP + 大量依赖 | Go 单二进制 |
| 内存占用 | 较高 | 目标空闲 < 50MB |
| 功能范围 | 全功能 + 插件商店 | LNMP + 三类站点管理 |
| 部署 | 安装脚本 + 面板服务 | 同左，但更轻 |
| CLI | `bt` 命令 | `zpanel` / `zp` 命令 |

### 1.2 支持的项目类型（v1）

| 类型 | 说明 | 面板职责 |
|------|------|----------|
| **HTML** | 静态站点 | Nginx root 指向站点目录 |
| **PHP** | WordPress、Discuz 等 | Nginx + PHP-FPM pool + 站点目录 |
| **Go** | 编译后的 Web 服务 | systemd 守护进程 + Nginx 反向代理 |

### 1.3 设计原则

- **单二进制部署**：前端构建产物通过 `go:embed` 嵌入，用户无需 Node 环境
- **不重复造轮子**：LNMP 用系统包管理器安装，面板只做配置编排与生命周期管理
- **API First**：所有 Web UI 操作均有对应 REST / WebSocket API
- **安全默认收紧**：HTTPS、路径沙箱、命令白名单、审计日志

---

## 2. 功能规划

### 2.1 功能模块总览

```
Zpanel
├── 系统
│   ├── 概览（CPU / 内存 / 磁盘 / 网络 / 负载）
│   ├── 进程管理
│   └── 服务管理（systemd）
├── 环境
│   ├── LNMP 一键安装 / 卸载 / 修复
│   ├── 组件状态（Nginx / MySQL / PHP-FPM）
│   └── 版本切换（PHP 多版本，可选）
├── 网站
│   ├── 站点列表（HTML / PHP / Go）
│   ├── 添加 / 删除 / 启停站点
│   ├── 域名绑定
│   ├── SSL 证书（Let's Encrypt）
│   └── 配置文件编辑（Nginx server block）
├── 数据库（PHP 配套）
│   ├── 库 / 用户 CRUD
│   └── 备份 / 还原（mysqldump）
├── 文件
│   ├── 目录浏览（沙箱路径）
│   ├── 上传 / 下载 / 重命名 / 删除
│   └── 在线编辑（文本，大小限制）
├── 计划任务
│   └── Crontab 可视化管理
├── 防火墙
│   └── ufw / firewalld 规则管理
├── 日志
│   ├── Nginx access / error
│   ├── PHP-FPM / 应用日志
│   └── 面板操作审计日志
└── 设置
    ├── 面板端口 / 绑定地址
    ├── 管理员账号 / 2FA
    └── 备份面板配置
```

### 2.2 版本分期

#### Phase 1 — MVP（4~6 周）

- [ ] 面板登录 + JWT 认证
- [ ] 系统概览监控
- [ ] systemd 服务列表与启停
- [ ] LNMP 一键安装脚本对接
- [ ] HTML / PHP / Go 站点添加与删除
- [ ] Nginx 配置自动生成与 reload
- [ ] `zpanel` CLI 基础命令
- [ ] 一键安装 / 卸载脚本

#### Phase 2 — 日常运维（4 周）

- [ ] SSL 自动申请与续期（ACME）
- [ ] 文件管理（沙箱）
- [ ] MySQL 库管理 + 备份
- [ ] 计划任务管理
- [ ] 日志查看与 tail

#### Phase 3 — 增强（按需）

- [ ] PHP 多版本切换
- [ ] 站点备份一键打包
- [ ] Web 终端（WebSocket + PTY）
- [ ] 多用户 / 子账号权限

### 2.3 站点模型（统一数据结构）

```yaml
# 站点配置示例 configs/sites/example.com.yaml
id: "uuid"
name: "example.com"
type: "php"          # html | php | go
status: "running"    # running | stopped | error
domains:
  - "example.com"
  - "www.example.com"
root: "/var/www/example.com"
ssl:
  enabled: true
  cert_path: "/etc/zpanel/ssl/example.com/fullchain.pem"
  key_path: "/etc/zpanel/ssl/example.com/privkey.pem"
  auto_renew: true
php:
  version: "8.2"
  pool: "www"
go:
  binary: "/var/www/example.com/app"
  port: 8081
  systemd_unit: "zpanel-site-example.service"
nginx:
  config_path: "/etc/nginx/sites-available/example.com.conf"
created_at: "2026-07-04T10:00:00Z"
updated_at: "2026-07-04T10:00:00Z"
```

---

## 3. 技术架构

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                      Browser                            │
│              Vue 3 SPA（embed 进二进制）                  │
└────────────────────────┬────────────────────────────────┘
                         │ HTTPS
                         ▼
┌─────────────────────────────────────────────────────────┐
│                   Zpanel Server (Go)                    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────────┐  │
│  │  Auth   │ │ Monitor │ │  Site   │ │   LNMP       │  │
│  │  JWT    │ │ /proc   │ │ Manager │ │   Installer  │  │
│  └─────────┘ └─────────┘ └─────────┘ └──────────────┘  │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────────┐  │
│  │  File   │ │  MySQL  │ │  Cron   │ │   Firewall   │  │
│  │ Sandbox │ │ Manager │ │ Manager │ │   Manager    │  │
│  └─────────┘ └─────────┘ └─────────┘ └──────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         ▼               ▼               ▼
   ┌──────────┐   ┌────────────┐   ┌──────────┐
   │ systemd  │   │ nginx      │   │ mysql    │
   │ journal  │   │ php-fpm    │   │ certbot  │
   └──────────┘   └────────────┘   └──────────┘
```

### 3.2 技术栈

| 层级 | 选型 | 说明 |
|------|------|------|
| 后端语言 | Go 1.22+ | 单二进制、低内存 |
| Web 框架 | **Gin** | 生态成熟、中间件丰富 |
| 前端框架 | **Vue 3 + Vite + TypeScript** | 组件化，适合复杂后台 |
| UI 库 | **Naive UI** | 见 [§8](#8-ui-库与视觉设计) |
| 配置 | YAML | 人类可读，可 Git 管理 |
| 数据库 | SQLite（面板自身） | 存用户、站点元数据、审计日志 |
| 系统交互 | `os/exec` + systemd D-Bus | 不引入 shell 注入 |
| SSL | lego（Go ACME 客户端） | 纯 Go，无 certbot 依赖 |
| 实时通信 | WebSocket（gorilla/websocket） | 监控推送、日志 tail |

### 3.3 端口与路径约定

| 项目 | 默认值 | 说明 |
|------|--------|------|
| 面板端口 | `8888` | 安装时可改 |
| 面板配置 | `/etc/zpanel/config.yaml` | 主配置 |
| 面板数据 | `/var/lib/zpanel/` | SQLite、站点 YAML、SSL |
| 站点根目录 | `/var/www/` | 与宝塔习惯一致 |
| Nginx 配置 | `/etc/nginx/sites-available/` | 软链到 sites-enabled |
| 日志 | `/var/log/zpanel/` | 面板日志 + 审计 |
| 二进制 | `/usr/local/bin/zpanel` | CLI 与 service 共用 |

---

## 4. 前端开发规范

### 4.1 项目结构

```
web/
├── index.html
├── vite.config.ts
├── package.json
├── tsconfig.json
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/
│   │   └── index.ts
│   ├── stores/              # Pinia 状态
│   │   ├── auth.ts
│   │   └── system.ts
│   ├── api/                 # API 封装
│   │   ├── client.ts
│   │   ├── site.ts
│   │   └── monitor.ts
│   ├── views/
│   │   ├── Dashboard.vue    # 系统概览
│   │   ├── Sites/
│   │   │   ├── List.vue
│   │   │   └── Create.vue
│   │   ├── Environment.vue  # LNMP 管理
│   │   ├── Files.vue
│   │   ├── Database.vue
│   │   ├── Cron.vue
│   │   ├── Logs.vue
│   │   └── Settings.vue
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Sidebar.vue
│   │   │   └── Header.vue
│   │   ├── SiteTypeTag.vue
│   │   └── StatCard.vue
│   └── styles/
│       └── global.css
└── dist/                    # build 输出，由 Go embed
```

### 4.2 路由规划

| 路径 | 页面 | 说明 |
|------|------|------|
| `/login` | 登录 | 无布局 |
| `/` | 概览 | 监控仪表盘 |
| `/sites` | 网站列表 | 支持类型筛选 |
| `/sites/create` | 添加站点 | 分步向导 |
| `/sites/:id` | 站点详情 | 域名、SSL、配置 |
| `/environment` | 环境管理 | LNMP 安装与状态 |
| `/database` | 数据库 | MySQL 管理 |
| `/files` | 文件管理 | 沙箱目录 |
| `/cron` | 计划任务 | |
| `/logs` | 日志 | |
| `/settings` | 设置 | 面板配置 |

### 4.3 API 调用规范

```typescript
// web/src/api/client.ts
const BASE = '/api/v1'

export async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const token = localStorage.getItem('token')
  const res = await fetch(`${BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options?.headers,
    },
  })
  if (res.status === 401) {
    router.push('/login')
    throw new Error('Unauthorized')
  }
  const data = await res.json()
  if (!data.ok) throw new Error(data.message)
  return data.data
}
```

### 4.4 构建与嵌入

```bash
# 开发
cd web && npm run dev          # 代理到 Go :8888

# 生产构建
cd web && npm run build        # 输出到 web/dist/

# Go embed（internal/web/embed.go）
//go:embed all:dist
var StaticFS embed.FS
```

```typescript
// vite.config.ts — 生产 base 路径
export default defineConfig({
  base: '/',
  build: { outDir: 'dist', emptyOutDir: true },
  server: {
    proxy: { '/api': 'http://127.0.0.1:8888' }
  }
})
```

### 4.5 前端开发约定

- 使用 **Composition API** + `<script setup lang="ts">`
- 状态管理：**Pinia**（仅全局状态：auth、system、sidebar）
- 图表：**ECharts**（监控曲线）或 **uPlot**（更轻）
- 代码编辑器：**CodeMirror 6**（Nginx 配置、文件编辑）
- 禁止在前端存储敏感信息（仅 JWT token）
- 所有表单提交需 loading 态 + 错误提示
- 列表页统一：分页、搜索、空状态、操作确认弹窗

---

## 5. 后端开发规范

### 5.1 项目结构

```
zpanel/
├── cmd/
│   └── zpanel/
│       └── main.go              # 入口：server / cli 子命令分发
├── internal/
│   ├── app/
│   │   └── server.go            # HTTP 服务启动
│   ├── cli/                     # CLI 命令实现
│   │   ├── root.go
│   │   ├── install.go
│   │   ├── default.go
│   │   └── ...
│   ├── auth/
│   │   ├── jwt.go
│   │   └── middleware.go
│   ├── config/
│   │   └── config.go
│   ├── handler/                 # HTTP handlers
│   │   ├── monitor.go
│   │   ├── site.go
│   │   ├── lnmp.go
│   │   └── ...
│   ├── service/                 # 业务逻辑
│   │   ├── monitor/
│   │   ├── site/
│   │   ├── nginx/
│   │   ├── lnmp/
│   │   ├── mysql/
│   │   ├── ssl/
│   │   └── systemd/
│   ├── model/                   # 数据模型
│   ├── store/                   # SQLite 存储
│   └── web/
│       └── embed.go             # 前端静态资源
├── templates/                   # Nginx / systemd 模板
│   ├── nginx/
│   │   ├── html.conf.tmpl
│   │   ├── php.conf.tmpl
│   │   └── go-proxy.conf.tmpl
│   └── systemd/
│       └── go-site.service.tmpl
├── scripts/
│   ├── install.sh
│   └── uninstall.sh
├── configs/
│   └── config.example.yaml
├── docs/
│   └── DEVELOPMENT.md
├── web/                         # 前端源码
├── go.mod
└── Makefile
```

### 5.2 API 设计

**Base URL：** `/api/v1`

**统一响应格式：**

```json
{
  "ok": true,
  "message": "",
  "data": {}
}
```

**错误响应：**

```json
{
  "ok": false,
  "message": "site not found",
  "code": "SITE_NOT_FOUND"
}
```

#### 核心 API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/auth/login` | 登录 |
| GET | `/monitor/overview` | 系统概览 |
| GET | `/monitor/processes` | 进程列表 |
| GET | `/services` | systemd 服务列表 |
| POST | `/services/:name/:action` | start/stop/restart |
| GET | `/lnmp/status` | LNMP 组件状态 |
| POST | `/lnmp/install` | 触发 LNMP 安装 |
| GET | `/sites` | 站点列表 |
| POST | `/sites` | 创建站点 |
| GET | `/sites/:id` | 站点详情 |
| PUT | `/sites/:id` | 更新站点 |
| DELETE | `/sites/:id` | 删除站点 |
| POST | `/sites/:id/ssl` | 申请 SSL |
| GET | `/files?path=` | 列出目录 |
| POST | `/files/upload` | 上传文件 |
| GET | `/database/databases` | 数据库列表 |
| POST | `/database/backup` | 备份数据库 |
| GET | `/cron` | 计划任务列表 |
| POST | `/cron` | 添加任务 |
| WS | `/ws/logs` | 日志实时推送 |
| WS | `/ws/monitor` | 监控实时推送 |

### 5.3 LNMP 安装逻辑

面板**不自行编译** LNMP，而是调用安装脚本 + 包管理器：

```
用户点击「安装 LNMP」或执行 zpanel lnmp install
        │
        ▼
Go 后端校验 root 权限 + 检测 OS（Ubuntu/Debian/CentOS）
        │
        ▼
执行 scripts/lnmp-install.sh
        │
        ├── apt/dnf 安装 nginx mysql-server php-fpm php-mysql ...
        ├── 写入默认 php.ini / pool 配置
        ├── systemctl enable --now nginx mysql php*-fpm
        └── 输出 JSON 状态供面板读取
        │
        ▼
Go 后端解析结果 → 更新 SQLite lnmp_status 表 → 前端展示
```

### 5.4 站点创建流程

#### HTML 站点

```
1. 创建 /var/www/{domain}/
2. 渲染 templates/nginx/html.conf.tmpl
3. 写入 /etc/nginx/sites-available/{domain}.conf
4. ln -s 到 sites-enabled
5. nginx -t && systemctl reload nginx
6. 保存站点 YAML + SQLite 记录
```

#### PHP 站点

```
1. 创建站点目录 + 设置 owner 为 www-data
2. 渲染 php.conf.tmpl（含 fastcgi_pass unix socket）
3. Nginx 配置 + reload
4. （可选）创建 MySQL 库和用户
5. 保存记录
```

#### Go 站点

```
1. 指定二进制路径 + 监听端口
2. 渲染 systemd/go-site.service.tmpl → /etc/systemd/system/
3. systemctl enable --now zpanel-site-{name}
4. 渲染 go-proxy.conf.tmpl（Nginx 反代到 127.0.0.1:port）
5. Nginx reload + 保存记录
```

### 5.5 命令执行安全

**禁止：** 将用户输入直接拼接到 shell 命令。

**正确做法：**

```go
// 白名单命令
var allowedCommands = map[string][]string{
    "nginx_test":  {"nginx", "-t"},
    "nginx_reload": {"systemctl", "reload", "nginx"},
}

func runCommand(name string, extraArgs ...string) error {
    base, ok := allowedCommands[name]
    if !ok {
        return ErrCommandNotAllowed
    }
    cmd := exec.Command(base[0], append(base[1:], extraArgs...)...)
    return cmd.Run()
}
```

---

## 6. CLI 命令行（zpanel / zp）

对标宝塔 `bt` 命令，提供无 UI 管理能力。安装后创建软链：

```bash
ln -sf /usr/local/bin/zpanel /usr/bin/zp
```

### 6.1 命令总览

```
zpanel [command] [flags]

Commands:
  default         查看面板入口信息（IP、端口、账号）— 对标 bt default
  info            显示面板运行状态
  start           启动面板服务
  stop            停止面板服务
  restart         重启面板服务
  reload          重载配置（不中断服务）
  status          服务状态

  user            管理员账号管理
    user show                     显示当前管理员
    user password [新密码]         修改密码
    user username [新用户名]       修改用户名

  port            修改面板端口
    port show                     显示当前端口
    port set <port>               设置端口

  ssl             面板 SSL 管理
    ssl enable                    启用面板 HTTPS
    ssl disable                   禁用面板 HTTPS

  lnmp            LNMP 环境管理
    lnmp status                   查看 Nginx/MySQL/PHP 状态
    lnmp install                  一键安装 LNMP
    lnmp uninstall                卸载 LNMP（需确认）
    lnmp repair                   修复常见配置问题

  site            站点管理
    site list                     列出所有站点
    site add                      交互式添加站点
    site delete <domain>          删除站点
    site start <domain>           启动站点
    site stop <domain>            停止站点

  firewall        防火墙管理
    firewall status               查看状态
    firewall open <port>          开放端口
    firewall close <port>         关闭端口

  logs            日志
    logs panel                    面板日志
    logs audit                    审计日志
    logs nginx                    Nginx 错误日志

  update          更新面板
    update check                  检查新版本
    update install                安装更新

  uninstall       卸载面板（保留网站数据）
  help            显示帮助
  version         显示版本号
```

### 6.2 交互式菜单（对标 bt 无参数执行）

无参数执行 `zpanel` 或 `zp` 时，进入交互菜单：

```
===============================================
  Zpanel 命令行管理工具
  Version: 1.0.0
===============================================
  1. 查看面板入口信息
  2. 修改面板密码
  3. 修改面板端口
  4. 启动面板
  5. 停止面板
  6. 重启面板
  7. 查看 LNMP 状态
  8. 一键安装 LNMP
  9. 查看站点列表
  10. 卸载面板
  0. 退出
===============================================
请输入命令编号：
```

### 6.3 CLI 实现（cobra）

```go
// internal/cli/root.go
var rootCmd = &cobra.Command{
    Use:   "zpanel",
    Short: "Zpanel Linux 面板管理工具",
    Run:   runInteractiveMenu, // 无子命令时进入菜单
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

```go
// cmd/zpanel/main.go
func main() {
    if len(os.Args) > 1 && os.Args[1] == "server" {
        app.RunServer()
        return
    }
    cli.Execute()
}
```

### 6.4 systemd 服务

```ini
# /etc/systemd/system/zpanel.service
[Unit]
Description=Zpanel Linux Admin Panel
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/zpanel server
Restart=on-failure
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

---

## 7. 一键部署脚本

### 7.1 安装脚本 `scripts/install.sh`

对标宝塔安装方式：`curl -sSO url && bash install.sh`

```bash
#!/bin/bash
# Zpanel 一键安装脚本
# 用法: curl -sSL https://get.zpanel.io/install.sh | bash
#   或: bash install.sh --port 8888 --username admin

set -euo pipefail

# ---------- 配置 ----------
ZPANEL_VERSION="${ZPANEL_VERSION:-latest}"
INSTALL_DIR="/usr/local/zpanel"
BIN_PATH="/usr/local/bin/zpanel"
CONFIG_DIR="/etc/zpanel"
DATA_DIR="/var/lib/zpanel"
LOG_DIR="/var/log/zpanel"
WWW_ROOT="/var/www"
DEFAULT_PORT=8888

# ---------- 颜色 ----------
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

# ---------- 权限检查 ----------
[[ $EUID -ne 0 ]] && error "请使用 root 用户运行此脚本"

# ---------- 系统检测 ----------
detect_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
    else
        error "不支持的操作系统"
    fi
    case "$OS" in
        ubuntu|debian|centos|rocky|almalinux) info "检测到系统: $OS $VER" ;;
        *) error "暂不支持 $OS，目前支持: Ubuntu, Debian, CentOS, Rocky" ;;
    esac
}

# ---------- 安装依赖 ----------
install_deps() {
    info "安装基础依赖..."
    case "$OS" in
        ubuntu|debian)
            apt-get update -qq
            apt-get install -y -qq curl wget sqlite3
            ;;
        centos|rocky|almalinux)
            yum install -y curl wget sqlite
            ;;
    esac
}

# ---------- 下载二进制 ----------
download_binary() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)  ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        *) error "不支持的架构: $ARCH" ;;
    esac

    URL="https://github.com/your-org/zpanel/releases/${ZPANEL_VERSION}/zpanel-linux-${ARCH}.tar.gz"
    info "下载 Zpanel ${ZPANEL_VERSION} (${ARCH})..."

    TMP=$(mktemp -d)
    curl -sSL "$URL" | tar xz -C "$TMP"
    install -m 755 "$TMP/zpanel" "$BIN_PATH"
    rm -rf "$TMP"
    ln -sf "$BIN_PATH" /usr/bin/zp
    info "二进制已安装: $BIN_PATH"
}

# ---------- 初始化配置 ----------
init_config() {
    local PORT=${1:-$DEFAULT_PORT}
    local USERNAME=${2:-admin}
    local PASSWORD=${3:-$(openssl rand -base64 12)}

    mkdir -p "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR" "$WWW_ROOT"

    if [[ ! -f "$CONFIG_DIR/config.yaml" ]]; then
        cat > "$CONFIG_DIR/config.yaml" <<EOF
panel:
  port: ${PORT}
  bind: "0.0.0.0"
  ssl: false
auth:
  username: "${USERNAME}"
  # 密码 bcrypt hash，由 zpanel 初始化命令写入
paths:
  www: "${WWW_ROOT}"
  data: "${DATA_DIR}"
  logs: "${LOG_DIR}"
EOF
        zpanel user password "$PASSWORD" 2>/dev/null || true
    fi
}

# ---------- 注册服务 ----------
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

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
    systemctl enable zpanel
    systemctl start zpanel
}

# ---------- 防火墙 ----------
setup_firewall() {
    local PORT=$1
    if command -v ufw &>/dev/null && ufw status | grep -q "Status: active"; then
        ufw allow "$PORT/tcp"
        info "ufw 已开放端口 $PORT"
    elif command -v firewall-cmd &>/dev/null; then
        firewall-cmd --permanent --add-port="${PORT}/tcp"
        firewall-cmd --reload
        info "firewalld 已开放端口 $PORT"
    fi
}

# ---------- 打印安装信息 ----------
print_info() {
    local PORT=$1 PASSWORD=$2
    local IP
    IP=$(curl -s --connect-timeout 3 ip.sb 2>/dev/null || hostname -I | awk '{print $1}')
    echo ""
    echo "============================================"
    echo -e "  ${GREEN}Zpanel 安装成功!${NC}"
    echo "============================================"
    echo -e "  面板地址: http://${IP}:${PORT}"
    echo -e "  用户名:   admin"
    echo -e "  密码:     ${PASSWORD}"
    echo -e "  管理命令: zpanel 或 zp"
    echo "============================================"
    echo -e "  ${YELLOW}请妥善保存以上信息${NC}"
    echo "============================================"
}

# ---------- 主流程 ----------
main() {
    local PORT=$DEFAULT_PORT USERNAME="admin" PASSWORD=""

    while [[ $# -gt 0 ]]; do
        case "$1" in
            --port)     PORT="$2"; shift 2 ;;
            --username) USERNAME="$2"; shift 2 ;;
            --password) PASSWORD="$2"; shift 2 ;;
            *) shift ;;
        esac
    done

    [[ -z "$PASSWORD" ]] && PASSWORD=$(openssl rand -base64 12)

    detect_os
    install_deps
    download_binary
    init_config "$PORT" "$USERNAME" "$PASSWORD"
    setup_service
    setup_firewall "$PORT"
    print_info "$PORT" "$PASSWORD"
}

main "$@"
```

### 7.2 LNMP 安装脚本 `scripts/lnmp-install.sh`

```bash
#!/bin/bash
# LNMP 一键安装 — 由 zpanel lnmp install 或面板 UI 调用
set -euo pipefail

install_ubuntu() {
    apt-get update -qq
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq \
        nginx \
        mysql-server \
        php-fpm php-mysql php-cli php-curl php-gd php-mbstring php-xml php-zip \
        unzip curl

    systemctl enable --now nginx mysql "php$(php -r 'echo PHP_MAJOR_VERSION.".".PHP_MINOR_VERSION;')-fpm"
}

install_centos() {
    yum install -y nginx mysql-server \
        php-fpm php-mysqlnd php-cli php-curl php-gd php-mbstring php-xml php-zip
    systemctl enable --now nginx mysqld php-fpm
}

# 输出 JSON 供 Go 解析
print_status() {
    php -r 'echo PHP_MAJOR_VERSION.".".PHP_MINOR_VERSION;' 2>/dev/null || echo "unknown"
}

case "$(. /etc/os-release && echo $ID)" in
    ubuntu|debian) install_ubuntu ;;
    centos|rocky|almalinux) install_centos ;;
    *) echo '{"ok":false,"message":"unsupported os"}'; exit 1 ;;
esac

echo "{\"ok\":true,\"nginx\":\"$(nginx -v 2>&1 | awk -F/ '{print $2}')\",\"php\":\"$(print_status)\",\"mysql\":\"$(mysql --version | awk '{print $3}')\"}"
```

### 7.3 卸载脚本 `scripts/uninstall.sh`

```bash
#!/bin/bash
# 卸载面板，可选保留网站数据
set -euo pipefail

KEEP_DATA=${1:-"yes"}

systemctl stop zpanel 2>/dev/null || true
systemctl disable zpanel 2>/dev/null || true
rm -f /etc/systemd/system/zpanel.service
rm -f /usr/local/bin/zpanel /usr/bin/zp

if [[ "$KEEP_DATA" != "yes" ]]; then
    rm -rf /etc/zpanel /var/lib/zpanel /var/log/zpanel
fi

systemctl daemon-reload
echo "Zpanel 已卸载。网站数据目录 /var/www 未删除。"
```

### 7.4 Makefile 构建发布

```makefile
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build frontend release

frontend:
	cd web && npm ci && npm run build

build: frontend
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/zpanel ./cmd/zpanel

release: build
	@for arch in amd64 arm64; do \
		tar czf bin/zpanel-linux-$$arch.tar.gz -C bin zpanel; \
	done
```

---

## 8. UI 库与视觉设计

### 8.1 UI 库：Naive UI（已确定）

| 项 | 说明 |
|----|------|
| 组件库 | Naive UI 2.x |
| 图标 | `@vicons/ionicons5`（线性图标，outline 风格） |
| 图表 | ECharts（低饱和配色，见 §8.4） |
| 代码编辑 | CodeMirror 6 |
| 状态 | Pinia |
| 断点 | `@vueuse/core` useBreakpoints |

选型理由：TypeScript 原生、tree-shaking、主题变量覆盖方便、组件覆盖后台全场景。

### 8.2 设计风格：运维控制台（Ops Console）

**一句话定位：** 像 GitHub Settings + 传统服务器面板，冷静、克制、信息优先。不是营销落地页，不是「AI 产品」视觉。

#### 设计气质

| 要 | 不要 |
|----|------|
| 中性灰白 + 单一功能色 | 紫蓝渐变、霓虹青、彩虹背景 |
| 扁平、细边框、小圆角 | 大圆角卡片、玻璃拟态、阴影堆叠 |
| 数据与表格优先 | 大插画、空状态表情包 |
| 文字 + 线性图标 | Emoji、3D 图标、卡通素材 |
| 状态色仅用于运行/停止/告警 | 到处滥用高饱和标签色 |
| 深色侧栏 + 浅色工作区 | 全屏渐变 Hero、发光按钮 |

#### 参考对象（学结构，不抄皮）

- **GitHub** — 设置页布局、表格密度、中性色
- **Vercel Dashboard** — 留白克制、层级清晰
- **Datadog / Grafana** — 监控数据呈现方式
- **传统 cPanel** — 运维功能的信息架构

#### 默认主题模式

**深色侧栏 + 浅色内容区**（默认），支持切换全暗色。侧栏深色让导航稳定，内容区浅色保证长时间阅读的舒适度。

```
┌──────────────────────────────────────────────────────────┐
│ ■ Zpanel                              admin ▾  退出     │  Header：白底、底边框
├──────────┬───────────────────────────────────────────────┤
│          │  概览                                         │
│  概览    │  ┌─────────┐ ┌─────────┐ ┌─────────┐         │
│  网站    │  │ CPU 12% │ │ 内存 45%│ │ 磁盘 60%│         │  统计卡片：白底细边框
│  环境    │  └─────────┘ └─────────┘ └─────────┘         │
│  数据库  │  ┌─────────────────────────────────────┐     │
│  文件    │  │ 站点列表（NDataTable）               │     │
│  计划任务│  └─────────────────────────────────────┘     │
│  日志    │                                               │
│  设置    │                                               │
│          │                                               │
│  # 深灰侧栏│                                               │
└──────────┴───────────────────────────────────────────────┘
```

### 8.3 色彩系统

只用 **中性色 + 1 个主色 + 4 个语义色**，禁止额外装饰色。

#### 主色（操作 / 链接 / 选中）

| 角色 | 色值 | 用途 |
|------|------|------|
| Primary | `#2563EB` | 主按钮、链接、菜单选中 |
| Primary Hover | `#1D4ED8` | 悬停 |
| Primary Pressed | `#1E40AF` | 按下 |

> 选标准蓝，不用靛紫（`#6366F1`）、不用青色（`#06B6D4`）。一个面板只需要一种「可点击」颜色。

#### 中性色（Light 内容区）

| 角色 | 色值 |
|------|------|
| 页面背景 | `#F4F4F5` |
| 卡片 / 面板 | `#FFFFFF` |
| 边框 | `#E4E4E7` |
| 分割线 | `#F4F4F5` |
| 主文字 | `#18181B` |
| 次要文字 | `#71717A` |
| 禁用文字 | `#A1A1AA` |

#### 侧栏（Dark）

| 角色 | 色值 |
|------|------|
| 侧栏背景 | `#18181B` |
| 侧栏边框 | `#27272A` |
| 菜单文字 | `#A1A1AA` |
| 菜单选中背景 | `#27272A` |
| 菜单选中文字 | `#FAFAFA` |

#### 语义色（仅状态场景）

| 状态 | 色值 | 用法 |
|------|------|------|
| 运行 / 成功 | `#16A34A` | 服务 running、SSL 有效 |
| 停止 / 错误 | `#DC2626` | 服务 stopped、操作失败 |
| 警告 | `#D97706` | 磁盘 >80%、证书即将过期 |
| 信息 | `#52525B` | 中性提示，不用蓝色重复 |

#### 图表配色（ECharts）

按顺序使用，不用渐变填充：

```
#2563EB  #16A34A  #D97706  #DC2626  #52525B
```

面积图可用 10% 透明度填充，**禁止** 紫→蓝→青渐变。

### 8.4 字体与排版

```css
/* 界面正文 — 不用 Inter，避免「泛 AI 感」 */
font-family: "Source Sans 3", "PingFang SC", "Microsoft YaHei", system-ui, sans-serif;

/* 日志、配置、路径、端口 */
font-family: "IBM Plex Mono", "SF Mono", "Consolas", monospace;
```

| 级别 | 大小 | 字重 | 场景 |
|------|------|------|------|
| 页面标题 | 20px | 600 | 「网站管理」「系统概览」 |
| 区块标题 | 16px | 600 | 卡片标题 |
| 正文 | 14px | 400 | 表格、表单、描述 |
| 辅助 | 12px | 400 | 时间戳、路径、版本号 |
| 代码/日志 | 13px | 400 | monospace |

行高：正文 `1.5`，表格 `1.4`。不使用过大的标题（避免 `32px+` 营销风标题）。

### 8.5 组件规范

| 组件 | 规范 |
|------|------|
| 圆角 | 统一 `4px`（`borderRadius: '4px'`），按钮/输入框/卡片一致 |
| 阴影 | 默认不用；弹窗/抽屉仅 `0 1px 3px rgba(0,0,0,.08)` |
| 边框 | 卡片用 `1px solid #E4E4E7`，不用纯阴影卡片 |
| 按钮 | 主操作 `primary`，次要 `default`，危险 `error`+二次确认 |
| 表格 | `NDataTable`，`size="small"`，斑马纹关闭，行悬停浅灰 |
| 标签 | 站点类型用 `NTag bordered`，不用 `strong` 荧光色 |
| 图标 | 仅导航和关键操作，16~18px，`outline` 风格 |
| 空状态 | 文字「暂无站点」+ 操作按钮，**不用** emoji 和插画 |

#### 站点类型标签色

| 类型 | 边框/文字色 | 背景 |
|------|-------------|------|
| HTML | `#52525B` | `#F4F4F5` |
| PHP | `#2563EB` | `#EFF6FF` |
| Go | `#16A34A` | `#F0FDF4` |

低饱和、带边框，不用实心大色块。

### 8.6 禁止清单

开发时 **不得** 出现以下内容：

- Emoji（包括空状态、通知、菜单）
- 紫蓝渐变背景（`linear-gradient(135deg, #667eea, #764ba2)` 类）
- 霓虹 / 赛博朋克配色
- 全屏 hero 大图、3D 插图、AI 生成装饰图
- 过大圆角（`rounded-2xl` / `16px+` 泛滥）
- 玻璃拟态（`backdrop-blur` 装饰）
- 彩虹色状态标签
- 「欢迎使用！」「让我们一起…」类营销文案
- Inter / Space Grotesk 作为唯一字体（可选辅助，不作主字体）

### 8.7 Naive UI 主题覆盖

```typescript
// web/src/styles/theme.ts
import type { GlobalThemeOverrides } from 'naive-ui'

export const lightThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#2563EB',
    primaryColorHover: '#1D4ED8',
    primaryColorPressed: '#1E40AF',
    primaryColorSuppl: '#3B82F6',
    borderRadius: '4px',
    borderRadiusSmall: '4px',
    fontFamily: '"Source Sans 3", "PingFang SC", "Microsoft YaHei", system-ui, sans-serif',
    fontFamilyMono: '"IBM Plex Mono", "SF Mono", Consolas, monospace',
    textColorBase: '#18181B',
    textColor1: '#18181B',
    textColor2: '#71717A',
    textColor3: '#A1A1AA',
    bodyColor: '#F4F4F5',
    dividerColor: '#E4E4E7',
    borderColor: '#E4E4E7',
  },
  Layout: {
    color: '#F4F4F5',
    siderColor: '#18181B',
    headerColor: '#FFFFFF',
    footerColor: '#F4F4F5',
  },
  Menu: {
    itemTextColor: '#A1A1AA',
    itemTextColorHover: '#FAFAFA',
    itemTextColorActive: '#FAFAFA',
    itemColorActive: '#27272A',
    itemColorHover: '#27272A',
    borderRadius: '4px',
  },
  Card: {
    color: '#FFFFFF',
    borderColor: '#E4E4E7',
    borderRadius: '4px',
  },
  Button: {
    borderRadiusSmall: '4px',
    borderRadiusMedium: '4px',
    borderRadiusLarge: '4px',
  },
  DataTable: {
    borderRadius: '4px',
    thColor: '#FAFAFA',
    tdColor: '#FFFFFF',
  },
  Tag: {
    borderRadius: '4px',
  },
}

export const darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#3B82F6',
    primaryColorHover: '#2563EB',
    primaryColorPressed: '#1D4ED8',
    borderRadius: '4px',
    bodyColor: '#09090B',
    cardColor: '#18181B',
    textColorBase: '#FAFAFA',
    textColor1: '#FAFAFA',
    textColor2: '#A1A1AA',
    borderColor: '#27272A',
    dividerColor: '#27272A',
  },
  Layout: {
    color: '#09090B',
    siderColor: '#18181B',
    headerColor: '#18181B',
  },
}
```

```typescript
// web/src/main.ts
import { create, NConfigProvider } from 'naive-ui'
import { lightThemeOverrides } from './styles/theme'

// App.vue
<NConfigProvider :theme-overrides="lightThemeOverrides">
  <router-view />
</NConfigProvider>
```

字体通过 `index.html` 引入（Google Fonts 或本地 self-host）：

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=Source+Sans+3:wght@400;500;600&display=swap" rel="stylesheet">
```

### 8.8 辅助库

| 用途 | 库 | 说明 |
|------|-----|------|
| 图标 | `@vicons/ionicons5` | 线性 outline，禁止 emoji |
| 图表 | `echarts` + `vue-echarts` | 配色见 §8.3 |
| 代码编辑 | `codemirror` 6 | 主题用 `oneDark` 或自定义灰调 |
| HTTP | 原生 fetch | — |
| 日期 | `date-fns` | — |
| 断点 | `@vueuse/core` | 响应式布局 |

---

## 9. 移动端适配策略

### 9.1 结论：**响应式适配，桌面优先**

| 策略 | 说明 |
|------|------|
| **桌面优先** | 主要使用场景是 PC 浏览器管理服务器 |
| **响应式布局** | 同一套代码，CSS 媒体查询适配平板 / 手机 |
| **不做独立 App** | 不做原生 App / 小程序 |
| **不做移动端专属功能** | 文件管理、代码编辑在手机上体验差，可提示「请使用 PC」 |

### 9.2 适配范围

| 页面 | 手机适配 | 说明 |
|------|----------|------|
| 登录 | ✅ 必须 | 简单表单 |
| 系统概览 | ✅ 必须 | 查看 CPU/内存，卡片纵向堆叠 |
| 站点列表 | ✅ 必须 | 列表改卡片，操作收进菜单 |
| 服务启停 | ✅ 必须 | 紧急运维场景 |
| LNMP 状态 | ✅ 必须 | 查看状态 |
| 添加站点 | ⚠️ 简化 | 分步表单单列布局 |
| 文件管理 | ❌ 仅提示 | 屏幕太小，提示用 PC |
| 代码/配置编辑 | ❌ 仅提示 | 同左 |
| 数据库管理 | ⚠️ 只读 | 手机可查看列表，复杂操作用 PC |
| 计划任务 | ⚠️ 只读 | 同左 |

### 9.3 断点与布局切换

| 断点 | 宽度 | 布局 |
|------|------|------|
| Desktop | ≥ 1200px | 固定侧栏 220px + 内容区 |
| Tablet | 768 ~ 1199px | 侧栏可折叠为图标模式（64px） |
| Mobile | < 768px | 无侧栏，顶部栏 + 底部 Tab + 抽屉菜单 |

#### Desktop / Tablet

- 侧栏始终深色（`#18181B`），与内容区形成稳定对比
- 内容区 `max-width` 不限制，表格可全宽
- Tablet 折叠时仅显示图标，hover 显示 tooltip

#### Mobile

```
┌─────────────────────────┐
│  ≡  Zpanel        admin │  顶栏：菜单按钮 + 标题 + 用户
├─────────────────────────┤
│                         │
│     主内容（单列）        │
│     卡片全宽堆叠          │
│                         │
├─────────────────────────┤
│  概览  网站  环境  更多   │  底部 Tab（4 项）
└─────────────────────────┘
```

- 「更多」打开 `NDrawer` 抽屉，内含完整菜单
- 表格改卡片列表：每行显示域名、类型标签、状态点、操作 `...` 菜单
- 统计数字放大到 24px，方便扫一眼
- 表单全部单列，`label-placement="top"`

```vue
<script setup lang="ts">
import { useBreakpoints } from '@vueuse/core'

const bp = useBreakpoints({ mobile: 768, tablet: 1200 })
const isMobile = bp.smaller('mobile')
const isTablet = bp.between('mobile', 'tablet')
const siderCollapsed = computed(() => isTablet.value)
</script>

<template>
  <NLayout has-sider style="min-height: 100vh">
    <!-- Desktop / Tablet 侧栏 -->
    <NLayoutSider
      v-if="!isMobile"
      :collapsed="siderCollapsed"
      :collapsed-width="64"
      :width="220"
      bordered
    />
    <NLayout>
      <NLayoutHeader bordered />
      <NLayoutContent :style="{ padding: isMobile ? '12px' : '24px' }">
        <router-view />
      </NLayoutContent>
      <!-- Mobile 底部 Tab -->
      <MobileTabBar v-if="isMobile" />
    </NLayout>
  </NLayout>
  <NDrawer v-if="isMobile" v-model:show="drawerOpen" placement="left" />
</template>
```

使用 `@vueuse/core` 切换布局，**同一套组件、同一套配色**，仅改变排列方式。

### 9.4 与宝塔对比

宝塔面板有移动端适配，但文件管理、插件安装等复杂操作体验一般。Zpanel 策略更务实：**保证监控和紧急运维可在手机完成，复杂配置引导到 PC**。

---

## 10. 目录结构

完整仓库结构（开发期）：

```
Zpanel/
├── cmd/zpanel/main.go
├── internal/
│   ├── app/
│   ├── auth/
│   ├── cli/
│   ├── config/
│   ├── handler/
│   ├── service/
│   ├── model/
│   ├── store/
│   └── web/embed.go
├── templates/
│   ├── nginx/
│   └── systemd/
├── scripts/
│   ├── install.sh
│   ├── uninstall.sh
│   └── lnmp-install.sh
├── configs/
│   └── config.example.yaml
├── docs/
│   └── DEVELOPMENT.md          ← 本文档
├── web/                        ← Vue 3 前端
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 11. 开发环境搭建

### 11.1 前置要求

| 工具 | 版本 |
|------|------|
| Go | 1.22+ |
| Node.js | 20+ |
| npm / pnpm | 最新 |
| Git | 任意 |
| Linux VM | Ubuntu 22.04 推荐（开发 systemd/nginx 调试） |

### 11.2 本地启动

```bash
# 1. 克隆
git clone https://github.com/your-org/zpanel.git
cd zpanel

# 2. 后端
cp configs/config.example.yaml configs/config.yaml
go run ./cmd/zpanel server

# 3. 前端（另开终端）
cd web
npm install
npm run dev
# 访问 http://localhost:5173，API 代理到 :8888
```

### 11.3 交叉编译

```bash
# Linux amd64（在 Windows/Mac 上编译）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/zpanel ./cmd/zpanel
```

---

## 12. 安全规范

### 12.1 必做项

- [ ] 首次安装强制修改默认密码
- [ ] 密码 bcrypt 存储，cost >= 12
- [ ] JWT 有效期 <= 24h，支持 refresh
- [ ] 登录失败限流：5 次 / 15 分钟
- [ ] 所有 API 需认证（除 login）
- [ ] 文件 API 路径沙箱，禁止 `../` 穿越
- [ ] 命令白名单，禁止任意 shell
- [ ] 面板操作写审计日志
- [ ] 生产环境强制 HTTPS 或限制 IP
- [ ] 敏感配置（DB 密码）加密存储

### 12.2 文件沙箱白名单

```yaml
# config.yaml
files:
  allowed_paths:
    - "/var/www"
    - "/etc/nginx/sites-available"
    - "/etc/nginx/sites-enabled"
    - "/var/log/nginx"
  max_upload_size: 50MB
  editable_extensions:
    - ".html" ".htm" ".css" ".js" ".json"
    - ".php" ".go" ".yaml" ".yml" ".conf"
    - ".txt" ".md" ".env" ".sql" ".sh"
```

---

## 13. 版本标记与 Git 备份

> **原则：每个可识别版本都必须有 Git 提交 + Tag，便于回滚与追溯。**

### 13.1 版本号规范

遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)：`MAJOR.MINOR.PATCH`

| 段 | 何时 +1 | 示例 |
|----|---------|------|
| MAJOR | 不兼容的大改、API 破坏性变更 | `1.0.0` 首个稳定版 |
| MINOR | 新功能，向后兼容 | `0.2.0` 新增 SSL 管理 |
| PATCH | Bug 修复，向后兼容 | `0.1.1` 修复 nginx reload |

**预发布阶段**（`0.x.y`）：主版本为 0 表示尚未稳定，可频繁迭代。

**版本存放位置（单一来源）：**

| 文件 | 用途 |
|------|------|
| `VERSION` | 纯文本版本号，脚本与 CI 读取 |
| `CHANGELOG.md` | 人类可读的变更记录 |
| Git Tag | `v0.1.0` 格式，不可变发布锚点 |
| Go 构建 | `-ldflags "-X main.version=$(cat VERSION)"` 写入二进制 |

### 13.2 Git 备份策略

#### 必须提交的情况

| 时机 | 操作 | 说明 |
|------|------|------|
| 功能完成 | `git commit` | 每个功能点一次提交，便于 bisect |
| 每日结束 | `git push` | 推送到远程，异地备份 |
| 版本发布 | `commit` + `tag` + `push --tags` | 每个版本一个 tag |
| 文档/设计变更 | `git commit` | 与代码同等重要 |
| 修复线上问题 | `commit` + 可选 `patch` tag | 如 `v0.1.1` |

#### 禁止入库

- `.env`、密钥、证书私钥
- `configs/config.yaml`（本地配置）
- `web/node_modules/`、`bin/` 构建产物
- SQLite 数据库、日志文件

见根目录 `.gitignore`。

#### 远程备份

```bash
# 首次关联远程（GitHub / Gitee / 自建 GitLab）
git remote add origin git@github.com:your-org/zpanel.git
git push -u origin main

# 每次发布
git push origin main
git push origin v0.1.0
```

**至少保留一个远程仓库**，本地 + 远程双备份。

### 13.3 分支策略

```
main              ← 稳定可发布，仅 merge 已测内容
develop           ← 日常开发主线（可选，小团队可直接 main）
feature/xxx       ← 功能分支，完成后 merge
fix/xxx           ← 修复分支
release/v0.2.0    ← 发布准备（可选）
```

小团队简化流程：**直接在 main 开发，功能 commit 勤提交，发布打 tag**。

### 13.4 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/)：

```
<type>(<scope>): <description>

feat(site): 支持 Go 站点反向代理
fix(nginx): 修复 reload 后配置未生效
docs: 更新版本管理章节
chore: release v0.2.0
refactor(auth): 简化 JWT 中间件
test(lnmp): 添加安装脚本单元测试
```

| type | 含义 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档 |
| `chore` | 构建、版本、杂项 |
| `refactor` | 重构 |
| `test` | 测试 |

### 13.5 版本发布流程（标准操作）

每次发版按顺序执行，**不可跳过 Git 步骤**：

```
1. 开发完成，所有改动已 commit
2. 更新 CHANGELOG.md（[Unreleased] → [X.Y.Z] - 日期）
3. 更新 VERSION 文件
4. git add VERSION CHANGELOG.md [其他变更]
5. git commit -m "chore: release vX.Y.Z"
6. git tag -a vX.Y.Z -m "Release vX.Y.Z"
7. git push origin main
8. git push origin vX.Y.Z
9. （可选）GitHub Release 上传二进制
```

或使用脚本：

```bash
# 递增版本号
bash scripts/bump-version.sh minor

# 手动完善 CHANGELOG 后
git add -A
git commit -m "chore: release v0.2.0"
bash scripts/release.sh 0.2.0
git push origin main --tags
```

### 13.6 回滚与恢复

| 场景 | 命令 |
|------|------|
| 查看所有版本 | `git tag -l 'v*'` |
| 切换到某版本只读查看 | `git checkout v0.1.0` |
| 基于旧版本修 bug | `git checkout -b fix/xxx v0.1.0` |
| 回退 main 到某 tag（慎用） | `git revert` 或 `git reset --hard v0.1.0` + force push |
| 对比两版本差异 | `git diff v0.1.0 v0.2.0` |

**推荐**：用 `git revert` 撤销错误提交，保留历史；避免 force push 到 main。

### 13.7 开发阶段里程碑与 Tag 计划

| 版本 | 里程碑 | Tag |
|------|--------|-----|
| v0.1.0 | 项目文档 + 规范 | `v0.1.0` |
| v0.2.0 | Go 骨架 + 登录 + 监控 API | `v0.2.0` |
| v0.3.0 | Vue 前端 + 系统概览页 | `v0.3.0` |
| v0.4.0 | LNMP 安装 + 站点管理 | `v0.4.0` |
| v0.5.0 | CLI + install.sh 可用 | `v0.5.0` |
| v1.0.0 | MVP 完整可用 | `v1.0.0` |

每个里程碑完成后：**commit → tag → push**，即使只是文档阶段也要 tag。

### 13.8 发布 Checklist

- [ ] 所有改动已 commit，工作区干净
- [ ] `VERSION` 与 tag 一致
- [ ] `CHANGELOG.md` 已更新该版本条目
- [ ] 前端 `npm run build` 无报错（有前端时）
- [ ] Go 测试通过 `go test ./...`（有后端时）
- [ ] 交叉编译 amd64 + arm64（发布二进制时）
- [ ] 在干净 Ubuntu VM 测试 install.sh（发布安装脚本时）
- [ ] `git tag -a vX.Y.Z` 已创建
- [ ] `git push origin main --tags` 已执行
- [ ] GitHub Release 已创建（可选）

### 13.9 日常开发节奏建议

```
上午/开始：git pull
开发中：    每完成一个小功能 → git commit
下午/结束：git push（远程备份）
发版：      走 §13.5 完整流程
```

**Commit 宜小不宜大**：一个功能、一个修复、一次文档更新 = 一次 commit。出问题时能精确回退到某一个改动。

---

## 附录 A：Nginx 模板示例

### PHP 站点 `templates/nginx/php.conf.tmpl`

```nginx
server {
    listen 80;
    server_name {{ range .Domains }}{{ . }} {{ end }};
    root {{ .Root }};
    index index.php index.html;

    access_log /var/log/nginx/{{ .Name }}_access.log;
    error_log  /var/log/nginx/{{ .Name }}_error.log;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass unix:/run/php/php{{ .PHPVersion }}-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }

    location ~ /\.(ht|git|env) {
        deny all;
    }
}
```

### Go 反向代理 `templates/nginx/go-proxy.conf.tmpl`

```nginx
server {
    listen 80;
    server_name {{ range .Domains }}{{ . }} {{ end }};

    access_log /var/log/nginx/{{ .Name }}_access.log;
    error_log  /var/log/nginx/{{ .Name }}_error.log;

    location / {
        proxy_pass http://127.0.0.1:{{ .GoPort }};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 附录 B：开发任务分工建议

| 角色 | 负责模块 | 预估工时 |
|------|----------|----------|
| 后端 | auth, monitor, systemd, CLI | 2 周 |
| 后端 | site, nginx, lnmp, ssl | 3 周 |
| 后端 | file, mysql, cron, firewall | 2 周 |
| 前端 | 布局 + 概览 + 站点管理 | 2 周 |
| 前端 | 环境 + 文件 + 数据库 + 设置 | 2 周 |
| 运维 | install.sh, lnmp-install.sh, 测试 | 1 周 |

---

*文档随项目演进持续更新。如有疑问，在 GitHub Issues 讨论。*
