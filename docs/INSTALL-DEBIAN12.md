# Zpanel 安装教程 — Debian 12

> 轻量版 Linux 服务器运维面板，支持 LNMP 一键安装与 PHP / HTML / Go 站点管理。  
> 项目地址：[https://github.com/Cyruss-top/Zpanel](https://github.com/Cyruss-top/Zpanel)

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

> **没有 curl？** 先执行 `apt update && apt install -y curl wget`，或直接用下面的 **wget** 命令。

### 方式 A：交互式安装（推荐）

安装过程中可自定义 **端口、用户名、密码、安全入口后缀**：

```bash
# 使用 curl
curl -sSL https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh -o install.sh
bash install.sh --interactive
```

```bash
# 没有 curl 时用 wget（Debian 最小化镜像常见）
wget -qO install.sh https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh
bash install.sh --interactive
```

```bash
# 一行命令（wget 管道）
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

```bash
# 使用 curl
curl -sSL https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh | bash -s -- \
  --port 9999 \
  --username myadmin \
  --password 'StrongPass123' \
  --entry mypanel
```

```bash
# 没有 curl 时用 wget
wget -qO- https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh | bash -s -- \
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
| `--interactive` / `-i` | 进入交互配置 |

### 方式 C：从源码本地安装（开发者）

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

# 或直接用 wget 安装面板
wget -qO- https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh | bash -s -- --interactive
```

### 1. 无法访问面板

```bash
zpanel status          # 检查服务是否运行
ss -tlnp | grep 8888   # 检查端口监听
zpanel default         # 查看正确访问地址（含安全入口）
```

### 2. 忘记密码

```bash
zpanel user password 新密码
```

### 3. 忘记安全入口后缀

```bash
grep entry /etc/zpanel/config.yaml
# 或清除
zpanel entry set clear
zpanel restart
```

### 4. LNMP 安装失败

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
# 拉取新版本后
cd Zpanel && git pull
make build-all
install -m 755 bin/zpanel /usr/local/bin/zpanel
systemctl restart zpanel
```

---

*教程随项目更新，问题反馈请提交 [GitHub Issues](https://github.com/Cyruss-top/Zpanel/issues)。*
