# Zpanel 开发计划

> 基于 [DEVELOPMENT.md](./DEVELOPMENT.md) 制定的可执行开发计划。  
> **当前版本：** v0.1.0 | **计划起始：** 2026-07-04 | **目标 MVP：** v1.0.0

---

## 1. 总览

### 1.1 目标

构建超级轻量版 Linux 可视化运维面板，核心能力：

- 一键安装 LNMP（Nginx + MySQL + PHP-FPM）
- 管理 HTML / PHP / Go 三类站点
- Go 单二进制 + Vue 3 嵌入部署
- CLI 管理（`zpanel` / `zp`，对标宝塔 `bt`）

### 1.2 里程碑与 Git Tag

| 版本 | 目标 | 预计周期 | Tag |
|------|------|----------|-----|
| **v0.1.0** | 文档 + 规范 + 开发计划 | 完成 | `v0.1.0` |
| **v0.2.0** | Go 后端骨架 + Auth + Monitor API | 第 1~2 周 | `v0.2.0` |
| **v0.3.0** | Vue 前端骨架 + 登录 + 概览页 | 第 2~3 周 | `v0.3.0` |
| **v0.4.0** | LNMP 安装 + 站点 CRUD + Nginx | 第 3~5 周 | `v0.4.0` |
| **v0.5.0** | CLI + install.sh + systemd | 第 5~6 周 | `v0.5.0` |
| **v1.0.0** | MVP 联调通过，可生产试用 | 第 6~8 周 | `v1.0.0` |
| **v1.1.0+** | Phase 2：SSL、文件、数据库、Cron | 第 9~12 周 | 按需 |

**规则：** 每个里程碑完成后必须 `commit → tag → push`。

---

## 2. 技术栈（不可随意变更）

| 层级 | 选型 |
|------|------|
| 后端 | Go 1.22+ / Gin / SQLite / cobra |
| 前端 | Vue 3 / Vite / TypeScript / Naive UI / Pinia |
| UI 风格 | 运维控制台：深侧栏 + 浅内容区，主色 `#2563EB`，禁止 emoji 和 AI 渐变配色 |
| 部署 | go:embed + systemd + install.sh |
| 配置 | YAML + 模板（Nginx / systemd） |

---

## 3. 阶段详细计划

### Phase 0 — 准备（v0.1.0，已完成）

- [x] 开发文档 `docs/DEVELOPMENT.md`
- [x] 版本规范 `VERSION` / `CHANGELOG.md` / `.gitignore`
- [x] Git 初始化 + tag `v0.1.0`
- [x] 开发计划 `docs/PLAN.md`
- [x] Agent Skill `.cursor/skills/zpanel-dev/`

---

### Phase A — 后端骨架（v0.2.0，第 1~2 周）

**目标：** 可运行的 Go HTTP 服务，提供认证与系统监控 API。

#### 第 1 周：项目脚手架

| 任务 | 产出 | 优先级 |
|------|------|--------|
| 初始化 `go.mod`，目录结构 | `cmd/` `internal/` 骨架 | P0 |
| 配置模块 | `internal/config/` 读取 YAML | P0 |
| HTTP Server | Gin 启动、健康检查 `/api/v1/health` | P0 |
| 统一响应格式 | `{ ok, message, data }` | P0 |
| SQLite 存储 | 用户表、审计日志表 | P0 |
| 静态资源 embed 占位 | `internal/web/embed.go` | P1 |
| Makefile | `make build` `make dev` | P1 |

**交付标准：**
```bash
go run ./cmd/zpanel server
curl http://127.0.0.1:8888/api/v1/health  # {"ok":true}
```

#### 第 2 周：认证 + 监控

| 任务 | 产出 | 优先级 |
|------|------|--------|
| JWT 登录 | `POST /api/v1/auth/login` | P0 |
| Auth 中间件 | Bearer token 校验 | P0 |
| 密码 bcrypt | 初始化默认 admin | P0 |
| 系统概览 API | `GET /api/v1/monitor/overview` | P0 |
| CPU/内存/磁盘 | 读 `/proc` `syscall.Statfs` | P0 |
| 进程列表 API | `GET /api/v1/monitor/processes` | P1 |
| systemd 服务 API | 列表 + start/stop/restart | P1 |
| 命令白名单 | `internal/service/systemd/` 安全封装 | P0 |

**Git：** 完成后更新 CHANGELOG → `v0.2.0` tag

---

### Phase B — 前端骨架（v0.3.0，第 2~3 周）

**目标：** Vue SPA 登录 + 布局 + 系统概览页，对接后端 API。

#### 第 2~3 周：前端初始化

| 任务 | 产出 | 优先级 |
|------|------|--------|
| Vite + Vue 3 + TS 初始化 | `web/` 目录 | P0 |
| Naive UI + 主题覆盖 | `web/src/styles/theme.ts` | P0 |
| 字体 Source Sans 3 + IBM Plex Mono | `index.html` | P0 |
| 布局组件 | Sidebar + Header + 响应式断点 | P0 |
| 路由 + 路由守卫 | `/login` `/` `/sites` ... | P0 |
| Pinia auth store | token 存取、登出 | P0 |
| API client | `web/src/api/client.ts` | P0 |
| 登录页 | 表单 + 错误提示 | P0 |
| 概览页 | CPU/内存/磁盘卡片 + ECharts 曲线 | P0 |
| Mobile 适配 | 底部 Tab + Drawer | P1 |

**设计约束：**
- 深侧栏 `#18181B`，内容区 `#F4F4F5`
- 圆角 4px，细边框，无 emoji，无渐变
- 断点：Desktop ≥1200 / Tablet 768~1199 / Mobile <768

**Git：** 完成后 `npm run build` 嵌入 Go → tag `v0.3.0`

---

### Phase C — LNMP 与站点（v0.4.0，第 3~5 周）

**目标：** 一键安装 LNMP，支持 HTML / PHP / Go 站点增删。

#### 第 3~4 周：LNMP + Nginx

| 任务 | 产出 | 优先级 |
|------|------|--------|
| `scripts/lnmp-install.sh` | apt/dnf 安装 nginx mysql php-fpm | P0 |
| LNMP 状态 API | 检测版本、运行状态 | P0 |
| LNMP 安装 API | 触发脚本、解析 JSON 结果 | P0 |
| Nginx 模板 | `templates/nginx/*.conf.tmpl` | P0 |
| systemd 模板 | `templates/systemd/go-site.service.tmpl` | P0 |
| 站点 model + store | SQLite + YAML 双写 | P0 |
| HTML 站点创建 | 目录 + nginx conf + reload | P0 |
| PHP 站点创建 | + php-fpm socket | P0 |
| Go 站点创建 | systemd unit + 反代 | P0 |
| 站点删除 | 清理配置、目录可选保留 | P0 |
| 前端：环境页 | LNMP 状态 + 安装按钮 | P0 |
| 前端：站点列表 + 创建向导 | 分类型表单 | P0 |

#### 第 4~5 周：站点完善

| 任务 | 产出 | 优先级 |
|------|------|--------|
| 站点启停 | Go 站点 systemd，PHP/HTML nginx | P0 |
| 域名多绑定 | server_name 更新 | P1 |
| Nginx 配置预览 | 只读展示 | P1 |
| 操作审计日志 | 写 SQLite audit 表 | P0 |
| 错误处理 | nginx -t 失败回滚 | P0 |

**Git：** tag `v0.4.0`

---

### Phase D — CLI 与部署（v0.5.0，第 5~6 周）

**目标：** `zpanel` 命令可用，install.sh 一键部署。

| 任务 | 产出 | 优先级 |
|------|------|--------|
| cobra CLI 框架 | `internal/cli/` | P0 |
| 交互式菜单 | 无参数执行 `zpanel` | P0 |
| default / start / stop / restart | systemd 控制面板 | P0 |
| user password / port set | 改密、改端口 | P0 |
| lnmp status / install | CLI 装环境 | P0 |
| site list / add / delete | CLI 管站点 | P1 |
| `scripts/install.sh` | 下载二进制、初始化、注册服务 | P0 |
| `scripts/uninstall.sh` | 卸载面板 | P1 |
| 交叉编译 | amd64 + arm64 Makefile | P0 |
| zp 软链 | `/usr/bin/zp` | P1 |

**Git：** tag `v0.5.0`

---

### Phase E — MVP 联调（v1.0.0，第 6~8 周）

**目标：** 完整 MVP 在 Ubuntu 22.04 干净环境跑通。

| 任务 | 产出 | 优先级 |
|------|------|--------|
| 端到端测试清单 | 安装→登录→装LNMP→建三类站点 | P0 |
| 安全加固 | 登录限流、路径沙箱、HTTPS 可选 | P0 |
| 文档完善 | README 安装说明 | P0 |
| Bug 修复 | 联调问题 | P0 |
| 性能验证 | 空闲内存 < 50MB | P1 |

**MVP 验收标准：**

1. `curl install.sh | bash` 安装成功
2. 浏览器登录面板
3. 一键安装 LNMP 成功
4. 分别创建 HTML、PHP、Go 站点并可访问
5. `zpanel site list` 显示站点
6. `zpanel stop/start` 控制面板
7. 每个版本有 Git tag

**Git：** tag `v1.0.0`

---

### Phase F — 日常运维（v1.1.0+，第 9~12 周，Phase 2）

按 [DEVELOPMENT.md §2.2](./DEVELOPMENT.md#22-版本分期) 排期：

| 版本 | 功能 |
|------|------|
| v1.1.0 | SSL 自动申请（lego ACME） |
| v1.2.0 | 文件管理（沙箱） |
| v1.3.0 | MySQL 库管理 + 备份 |
| v1.4.0 | 计划任务 + 日志 tail |

---

## 4. 任务依赖图

```
v0.1.0 文档
    │
    ▼
v0.2.0 Go 骨架 ──► Auth ──► Monitor API ──► systemd API
    │
    ▼
v0.3.0 Vue 骨架 ──► 登录页 ──► 布局 ──► 概览页
    │
    ▼
v0.4.0 LNMP 脚本 ──► Nginx 模板 ──► 站点 CRUD ──► 前端站点页
    │
    ▼
v0.5.0 CLI (cobra) ──► install.sh ──► 交叉编译
    │
    ▼
v1.0.0 联调 ──► 安全 ──► 发布
```

**可并行：**
- v0.2.0 后端 与 v0.3.0 前端（API 契约先定）
- 前端 Mock 数据在后端未完成时使用

---

## 5. 每周工作节奏

```
周一      确认本周任务，git pull
周二~周四  开发，每完成小功能 commit
周五      联调、修 bug、更新 CHANGELOG
里程碑    bump VERSION → tag → push
```

### Commit 规范

```
feat(auth): JWT 登录接口
feat(monitor): 系统概览 API
feat(web): 登录页与布局
fix(nginx): reload 前配置校验
docs: 更新开发计划
chore: release v0.2.0
```

---

## 6. 当前 Sprint（下一步）

**Sprint 1 目标：** 启动 v0.2.0，完成后端脚手架

| # | 任务 | 状态 |
|---|------|------|
| 1 | `go mod init`，创建目录结构 | 完成 |
| 2 | `cmd/zpanel/main.go` server/cli 入口 | 完成 |
| 3 | Gin HTTP + `/api/v1/health` | 完成 |
| 4 | config 读取 `configs/config.example.yaml` | 完成 |
| 5 | SQLite 初始化 + 用户表 | 完成 |
| 6 | JWT auth login API | 完成 |
| 7 | monitor overview API | 完成 |

**v0.2.0 已达成。** 下一步进入 **v0.3.0**：Vue 前端骨架 + 登录 + 概览页。

| # | 任务 | 状态 |
|---|------|------|
| 1 | Vite + Vue 3 + TS 初始化 | **下一步** |
| 2 | Naive UI + 主题 + 布局 | 待开始 |
| 3 | 登录页 + 路由守卫 | 待开始 |
| 4 | 概览页对接 monitor API | 待开始 |

**启动命令（开发环境）：**
```bash
cp configs/config.example.yaml configs/config.yaml
go run ./cmd/zpanel server
# 默认账号 admin / admin（首次启动日志会提示）
```

**API 测试：**
```bash
curl -X POST http://127.0.0.1:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

curl http://127.0.0.1:8888/api/v1/monitor/overview \
  -H "Authorization: Bearer <token>"
```

---

## 7. 风险与对策

| 风险 | 对策 |
|------|------|
| Windows 开发无法测 systemd/nginx | 用 Ubuntu VM 或 WSL2 做集成测试 |
| LNMP 各发行版差异 | install.sh 分 ubuntu/debian 和 centos 分支 |
| 安全漏洞（面板 = root 权限） | 命令白名单、路径沙箱、审计日志 |
| 前后端进度不一致 | 先定 OpenAPI 契约，前端 Mock |
|  scope 膨胀 | 严格按 Phase 分期，Phase 3 功能不进入 v1.0.0 |

---

## 8. 相关文档

- [DEVELOPMENT.md](./DEVELOPMENT.md) — 完整技术规范
- [PLAN.md](./PLAN.md) — 开发计划与 Sprint 排期
- [CHANGELOG.md](../CHANGELOG.md) — 版本变更
- [VERSION](../VERSION) — 当前版本号
- `.cursor/skills/zpanel-dev/SKILL.md` — Agent 开发技能

---

*本计划随开发进度更新，每个里程碑完成后同步修订。*
