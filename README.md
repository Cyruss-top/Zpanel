# Zpanel

超级轻量版 Linux 可视化运维面板 — Go 单二进制部署，支持 LNMP 一键安装与 PHP / HTML / Go 项目管理。

## 特性

- 单二进制 + 嵌入式前端，安装即用
- LNMP 一键安装（Nginx + MySQL + PHP-FPM）
- 站点管理：HTML 静态站、PHP 项目、Go 反向代理
- 系统监控、服务管理、SSL 证书、计划任务
- CLI 管理：`zpanel` / `zp` 命令（对标宝塔 `bt`）

## 快速安装

```bash
curl -sSL https://get.zpanel.io/install.sh | bash
```

## 管理命令

```bash
zpanel              # 交互式菜单
zpanel default      # 查看面板入口
zpanel lnmp install # 安装 LNMP
zpanel site list    # 站点列表
```

## 版本

当前版本：**v0.1.0**（见 [VERSION](VERSION)）

变更记录：[CHANGELOG.md](CHANGELOG.md)

版本管理与 Git 备份规范：[docs/DEVELOPMENT.md §13](docs/DEVELOPMENT.md#13-版本标记与-git-备份)

## 开发文档

详见 [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)

## 技术栈

- **后端：** Go 1.22+ / Gin
- **前端：** Vue 3 / Vite / TypeScript / Naive UI
- **部署：** go:embed 单二进制 + systemd

## License

MIT
