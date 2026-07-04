# Zpanel 安装教程 — Debian 12

> 轻量版 Linux 服务器运维面板，支持 LNMP 一键安装与 PHP / HTML / Go 站点管理。

| 节点 | 地址 |
|------|------|
| **中国大陆（推荐）** | [https://gitee.com/Ressss2023/Zpanel](https://gitee.com/Ressss2023/Zpanel) |
| 国际 | [https://github.com/Cyruss-top/Zpanel](https://github.com/Cyruss-top/Zpanel) |

---

## 一、环境要求

| 项目 | 要求 |
|------|------|
| 系统 | **Debian 12**（本教程）、Ubuntu 22.04+ |
| 架构 | x86_64 (amd64) / aarch64 (arm64) |
| 权限 | root 或 sudo |
| 内存 | 建议 1GB 以上 |
| 端口 | 默认 8888，可自定义 |

---

## 二、安装前准备

### 1. 更新系统

```bash
su -
apt update && apt upgrade -y
apt install -y curl wget ca-certificates
```

### 2. 确认防火墙（可选）

Debian 12 默认未启用 ufw。若已启用，需放行面板端口：

```bash
ufw allow 8888/tcp
ufw reload
```

---

## 三、一键安装

> **没有 curl？** 先执行 `apt update && apt install -y curl wget`，或直接用 **wget** 命令。  
> **中国大陆服务器** 请优先使用下方 **Gitee 节点**，速度更快、更稳定。

### 方式 A：交互式安装（推荐）

安装过程中可自定义 **端口、用户名、密码、安全入口后缀**。

#### 中国大陆节点（Gitee，推荐）

```bash
# wget 一行安装（推荐，无需 curl）
wget -qO- https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- --mirror gitee --interactive
```

```bash
# 下载后执行
wget -qO install.sh https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh
bash install.sh --mirror gitee --interactive
```

```bash
# 有 curl 时
curl -sSL https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- --mirror gitee --interactive
```

#### 国际节点（GitHub）

```bash
curl -sSL https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh -o install.sh
bash install.sh --interactive
```

```bash
wget -qO- https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh | bash -s -- --interactive
```

按提示输入：

| 提示 | 说明 | 示例 |
|------|------|------|
| 面板端口 | 访问面板的端口 | `8888` |
| 管理员用户名 | 登录账号 | `admin` |
| 管理员密码 | 留空则自动生成 | `MyPass123` |
| 安全入口后缀 | 类似宝塔安全入口，留空不启用 | `mypanel` |

设置后缀 `mypanel` 后，访问地址为：

```
http://你的服务器IP:8888/mypanel/
```

### 方式 B：命令行参数安装

#### 中国大陆（Gitee）

```bash
wget -qO- https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- \
  --mirror gitee \
  --port 9999 \
  --username myadmin \
  --password 'StrongPass123' \
  --entry mypanel
```

#### 国际（GitHub）

```bash
curl -sSL https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh | bash -s -- \
  --port 9999 \
  --username myadmin \
  --password 'StrongPass123' \
  --entry mypanel
```

| 参数 | 说明 |
|------|------|
| `--port` | 面板端口 |
| `--username` | 管理员用户名 |
| `--password` | 管理员密码 |
| `--entry` | 安全入口后缀（3~32 位字母数字） |
| `--mirror` | 下载镜像：`gitee`（国内）或 `github`（国际） |
| `--interactive` / `-i` | 进入交互配置 |

### 方式 C：从源码本地安装（开发者）

#### 中国大陆（Gitee）

```bash
apt install -y git golang nodejs npm
git clone https://gitee.com/Ressss2023/Zpanel.git
cd Zpanel
make build-all

ZPANEL_INSTALL_LOCAL=1 bash scripts/install.sh --interactive
```

#### 国际（GitHub）

```bash
apt install -y git golang nodejs npm
git clone https://github.com/Cyruss-top/Zpanel.git
cd Zpanel
make build-all

ZPANEL_INSTALL_LOCAL=1 bash scripts/install.sh --interactive
```

---

## 四、安装完成

成功后会显示：

```
============================================
  Zpanel 安装成功!
============================================
  面板地址: http://203.0.113.1:8888/mypanel/
  用户名:   myadmin
  密码:     xxxxx
  安全入口: /mypanel/
  管理命令: zpanel 或 zp
============================================
```

浏览器打开面板地址，使用账号密码登录。

---

## 五、命令行管理

安装后可使用 `zpanel` 或 `zp` 命令（对标宝塔 `bt`）：

```bash
zpanel                    # 交互式菜单
zpanel default            # 查看面板入口地址
zpanel version            # 查看版本

# 账号与端口
zpanel user password 新密码
zpanel user username 新用户名
zpanel port set 9999
zpanel entry set mypanel  # 设置安全入口
zpanel entry set clear    # 清除安全入口
zpanel restart            # 修改后重启生效

# 服务控制
zpanel start
zpanel stop
zpanel restart
zpanel status

# LNMP 环境
zpanel lnmp status
zpanel lnmp install

# 站点管理
zpanel site list
zpanel site delete example.com
```

---

## 六、安装 LNMP 环境

1. 登录面板
2. 左侧点击 **环境**
3. 点击 **一键安装 LNMP**

或使用命令：

```bash
zpanel lnmp install
```

安装完成后可创建站点：

- **HTML** 静态站
- **PHP** 动态站（WordPress 等）
- **Go** 反向代理

---

## 七、创建网站示例

### 面板操作

1. 进入 **网站** → **添加站点**
2. 填写站点名称、类型、域名
3. 点击创建

### 命令行查看

```bash
zpanel site list
```

---

## 八、自定义配置说明

### 安全入口后缀

防止面板被扫描，建议设置随机后缀：

```bash
zpanel entry set x7k9m2
zpanel restart
```

访问：`http://IP:端口/x7k9m2/`

### 修改端口

```bash
zpanel port set 9999
zpanel restart
```

### 修改密码

```bash
zpanel user password '新密码'
```

---

## 九、卸载

```bash
bash /usr/local/zpanel/scripts/uninstall.sh
# 或从仓库
bash scripts/uninstall.sh
```

默认保留 `/var/www` 网站数据。

---

## 十、常见问题

### 0. curl: command not found

Debian 最小化镜像可能未预装 curl，任选其一：

```bash
# 安装 curl（推荐）
apt update && apt install -y curl wget ca-certificates

# 或直接用 wget + Gitee 国内节点安装
wget -qO- https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- --mirror gitee --interactive
```

### 1. 国内下载 GitHub 很慢或失败

请改用 **Gitee 节点**：

```bash
wget -qO- https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- --mirror gitee --interactive
```

安装脚本会优先从 Gitee Releases 下载二进制，失败时自动回退 GitHub。

### 2. 无法访问面板

```bash
zpanel status          # 检查服务是否运行
ss -tlnp | grep 8888   # 检查端口监听
zpanel default         # 查看正确访问地址（含安全入口）
```

### 3. 忘记密码

```bash
zpanel user password 新密码
```

### 4. 忘记安全入口后缀

```bash
grep entry /etc/zpanel/config.yaml
# 或清除
zpanel entry set clear
zpanel restart
```

### 5. LNMP 安装失败

确保以 root 运行，且系统为 Debian/Ubuntu：

```bash
zpanel lnmp install
journalctl -u zpanel -n 50
```

---

## 十一、配置文件位置

| 文件 | 路径 |
|------|------|
| 主配置 | `/etc/zpanel/config.yaml` |
| 数据目录 | `/var/lib/zpanel/` |
| 日志 | `/var/log/zpanel/` |
| 网站根目录 | `/var/www/` |
| 二进制 | `/usr/local/bin/zpanel` |

---

## 十二、更新面板

```bash
# Gitee（国内）
cd Zpanel && git pull
# 若从 GitHub 克隆，可添加 Gitee 远程：
# git remote add gitee https://gitee.com/Ressss2023/Zpanel.git
# git pull gitee main

make build-all
install -m 755 bin/zpanel /usr/local/bin/zpanel
systemctl restart zpanel
```

---

*教程随项目更新。问题反馈：[Gitee Issues](https://gitee.com/Ressss2023/Zpanel/issues) · [GitHub Issues](https://github.com/Cyruss-top/Zpanel/issues)*
