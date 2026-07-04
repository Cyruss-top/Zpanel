# Changelog

本文件记录 Zpanel 每个版本的变更。格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

## [Unreleased]

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

[Unreleased]: https://github.com/your-org/zpanel/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/your-org/zpanel/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/your-org/zpanel/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/your-org/zpanel/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/your-org/zpanel/releases/tag/v0.1.0
