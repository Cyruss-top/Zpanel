# Zpanel

超级轻量版 Linux 可视化运维面板 — Go 单二进制部署，支持 LNMP 一键安装与 PHP / HTML / Go 项目管理。

| 节点 | 地址 |
|------|------|
| **中国大陆（推荐）** | [https://gitee.com/Ressss2023/Zpanel](https://gitee.com/Ressss2023/Zpanel) |
| 国际 | [https://github.com/Cyruss-top/Zpanel](https://github.com/Cyruss-top/Zpanel) |

## 特性

- 单二进制 + 嵌入式前端，安装即用
- LNMP 一键安装（Nginx + MySQL + PHP-FPM）
- 站点管理：HTML 静态站、PHP 项目、Go 反向代理
- 系统监控、服务管理
- 安全入口后缀（防扫描）
- CLI 管理：`zpanel` / `zp` 命令（对标宝塔 `bt`）

## 快速安装（Debian 12 / Ubuntu）

### 中国大陆（Gitee，推荐）

```bash
# 交互式安装
wget -qO- https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- --mirror gitee --interactive

# 指定端口、账号、密码、安全入口
wget -qO- https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- \
  --mirror gitee --port 8888 --username admin --password 'yourpass' --entry mypanel
```

### 国际（GitHub）

```bash
curl -sSL https://raw.githubusercontent.com/Cyruss-top/Zpanel/main/scripts/install.sh | bash -s -- --interactive
```

详细教程：[docs/INSTALL-DEBIAN12.md](docs/INSTALL-DEBIAN12.md)

当前版本：**v0.7.0**（见 [VERSION](VERSION)）

## 本地开发安装

```bash
make build-all
ZPANEL_INSTALL_LOCAL=1 sudo bash scripts/install.sh --interactive
```

## 管理命令

```bash
zpanel                    # 交互式菜单
zpanel default            # 查看面板入口
zpanel user password 新密码
zpanel user username 新用户名
zpanel port set 9999
zpanel entry set mypanel  # 设置安全入口后缀
zpanel entry set clear    # 清除安全入口
zpanel lnmp install       # 安装 LNMP
zpanel site list          # 站点列表
zpanel uninstall              # 交互选择卸载模式
zpanel uninstall -y --keep-www   # 保留 /var/www 网站
zpanel uninstall -y --all        # 彻底删除干净
```

## 版本

变更记录：[CHANGELOG.md](CHANGELOG.md)

版本管理与 Git 备份规范：[docs/DEVELOPMENT.md §13](docs/DEVELOPMENT.md#13-版本标记与-git-备份)

## 开发文档

- [docs/INSTALL-DEBIAN12.md](docs/INSTALL-DEBIAN12.md) — Debian 12 安装教程（含 Gitee 国内节点）
- [docs/GITEE-SYNC.md](docs/GITEE-SYNC.md) — Gitee 自动同步配置
- [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) — 完整技术规范
- [docs/PLAN.md](docs/PLAN.md) — 开发计划与 Sprint 排期

Agent Skill：`.cursor/skills/zpanel-dev/`（Zpanel 开发时自动加载）

## 技术栈

- **后端：** Go 1.22+ / Gin
- **前端：** Vue 3 / Vite / TypeScript / Naive UI
- **部署：** go:embed 单二进制 + systemd

## License

MIT
