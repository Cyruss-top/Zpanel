# Changelog

本文件记录 Zpanel 每个版本的变更。格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

## [Unreleased]

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

[Unreleased]: https://github.com/your-org/zpanel/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/your-org/zpanel/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/your-org/zpanel/releases/tag/v0.1.0
