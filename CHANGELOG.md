# Changelog

本文件记录 Zpanel 每个版本的变更。格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

## [Unreleased]

## [0.7.4] - 2026-07-04

### Changed
- 卸载支持两种模式：保留站点数据 / 彻底删除干净（含 /var/www）
- `uninstall.sh` 与 `zpanel uninstall` 交互式二选一菜单

## [0.7.3] - 2026-07-04

### Fixed
- 安全入口（如 `/admin/`）下静态资源 404 导致白屏：注入 `<base href>` 并重写 assets 路径
- `/admin` 无尾斜杠时重定向到 `/admin/`

## [0.7.2] - 2026-07-04

### Added
- `zpanel uninstall` 命令（支持 `-y` `--purge`）
- 交互式菜单「12. 卸载面板」
- 完善 `scripts/uninstall.sh`（确认提示、purge 选项）
- 安装时自动部署卸载脚本到 `/usr/local/zpanel/scripts/`

## [0.7.1] - 2026-07-04

### Fixed
- 安装完成与 `zpanel default` 优先显示公网 IP（多源探测 + 内网地址单独展示）

## [0.7.0] - 2026-07-04

### Added
- 安装脚本 `--base-url` / `ZPANEL_BASE_URL`：从自建服务器一键安装
- 安装脚本 `--package`：从本地 Release 包安装
- Gitee 国内镜像与 `docs/GITEE-SYNC.md` 同步指南
- GitHub Actions：Release 发布 + Gitee 代码同步

### Changed
- 安装脚本支持 wget、Gitee/GitHub 双镜像回退
- 站点管理、LNMP 环境相关功能完善

## [0.6.0] - 2026-07-04

### Added
- 安全入口后缀（`panel.entry`）：访问路径如 `/mypanel/`，防面板扫描
- 安装脚本支持自定义：`--port` `--username` `--password` `--entry` `--interactive`
- CLI：`zpanel entry show/set`、`zpanel user username`
- 交互式菜单新增：修改安全入口、修改管理员用户名
- Debian 12 中文安装教程 `docs/INSTALL-DEBIAN12.md`
- 前端动态 base 路径（`__ZPANEL_ENTRY__` 注入）

## [0.5.0] - 2026-07-04

### Added
- cobra CLI：`zpanel` / `zp` 命令（对标宝塔 `bt`）
- 交互式菜单（无参数执行 `zpanel`）
- 命令：`default` `start/stop/restart/status` `user password` `port set` `lnmp` `site list/delete` `version`
- `scripts/install.sh` 一键安装（支持 `ZPANEL_INSTALL_LOCAL=1` 本地安装）
- `scripts/uninstall.sh` 卸载脚本
- Makefile 交叉编译：`make release-amd64` `release-arm64` `release`

## [0.4.0] - 2026-07-04

### Added
- `scripts/lnmp-install.sh` 一键安装 LNMP（Ubuntu/Debian/CentOS）
- LNMP 状态检测与安装 API：`GET/POST /api/v1/lnmp/*`
- 站点 CRUD API：HTML / PHP / Go 三类站点
- Nginx 配置模板渲染、启用、reload（`nginx -t` 校验）
- Go 站点 systemd unit 管理
- 前端「环境管理」页 + 「网站」列表与创建页
- 开发环境 Nginx 配置目录（`./data/nginx/`）

## [0.3.0] - 2026-07-04

### Added
- Vue 3 + Vite + TypeScript + Naive UI 前端
- 运维控制台主题（深侧栏 + 浅内容区）
- 登录页 + JWT 路由守卫
- 系统概览页（CPU / 内存 / 磁盘 / 负载）
- 响应式布局（Desktop / Tablet / Mobile）
- Go 嵌入前端静态资源 + SPA 路由

## [0.2.0] - 2026-07-04

### Added
- SQLite 存储：用户表、审计日志表
- JWT 登录 `POST /api/v1/auth/login`（bcrypt + 限流）
- 系统监控 `GET /api/v1/monitor/overview`（Linux 读 /proc）
- 首次启动自动创建 admin 用户（默认密码 admin，日志警告）
- Gin HTTP 服务：`zpanel server` 启动面板
- `GET /api/v1/health` 健康检查 API
- 配置加载：`config.Load` + `ZPANEL_CONFIG` 环境变量 + 自动创建 data 目录
- Go 项目目录结构（`cmd/` `internal/` `templates/` `web/`）

### Added (v0.1.x)
- 开发计划 `docs/PLAN.md`（里程碑 v0.2.0 ~ v1.0.0）
- Agent Skill `.cursor/skills/zpanel-dev/`
- 本地开发配置示例 `configs/config.example.yaml`

---

## [0.1.0] - 2026-07-04

### Added
- 初始开发文档 `docs/DEVELOPMENT.md`
- 项目 README
- 版本标记与 Git 备份规范

[Unreleased]: https://github.com/your-org/zpanel/compare/v0.5.0...HEAD
[0.5.0]: https://github.com/your-org/zpanel/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/your-org/zpanel/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/your-org/zpanel/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/your-org/zpanel/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/your-org/zpanel/releases/tag/v0.1.0
