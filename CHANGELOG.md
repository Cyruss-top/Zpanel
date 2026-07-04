# Changelog

本文件记录 Zpanel 每个版本的变更。格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### Added
- Gin HTTP 服务：`zpanel server` 启动面板
- `GET /api/v1/health` 健康检查 API
- 配置加载：`config.Load` + `ZPANEL_CONFIG` 环境变量 + 自动创建 data 目录
- Go 项目目录结构（`cmd/` `internal/` `templates/` `web/`）
- `go.mod` 模块 `github.com/zex/zpanel`
- 统一 API 响应模型 `internal/model/response.go`
- Nginx / systemd 配置模板占位
- Makefile 构建入口

---

## [0.1.0] - 2026-07-04

### Added
- 初始开发文档 `docs/DEVELOPMENT.md`
- 项目 README
- 版本标记与 Git 备份规范

[Unreleased]: https://github.com/your-org/zpanel/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/your-org/zpanel/releases/tag/v0.1.0
