---
name: zpanel-dev
description: >-
  Develop Zpanel, a lightweight Linux admin panel (Go + Vue 3 + Naive UI).
  Use when working on Zpanel code, features, API, CLI, LNMP, site management,
  install scripts, UI, or when the user mentions Zpanel, zpanel, zp, or this repo.
---

# Zpanel Development

## Project Summary

Super-lightweight Linux visual ops panel (Baota-like). Go single binary, embedded Vue 3 SPA.

**Core capabilities:** LNMP install, HTML/PHP/Go site management, systemd, CLI (`zpanel`/`zp`).

**Docs (read when needed):**
- [docs/DEVELOPMENT.md](../../../docs/DEVELOPMENT.md) — full spec
- [docs/PLAN.md](../../../docs/PLAN.md) — sprint plan & milestones

**Current version:** read [VERSION](../../../VERSION). Milestone tags: v0.1.0 docs → v0.2.0 Go → v0.3.0 Vue → v0.4.0 sites → v0.5.0 CLI → v1.0.0 MVP.

## Tech Stack (Do Not Change Without User Approval)

| Layer | Choice |
|-------|--------|
| Backend | Go 1.22+, Gin, SQLite, cobra |
| Frontend | Vue 3, Vite, TypeScript, Naive UI, Pinia |
| Icons | `@vicons/ionicons5` (outline only) |
| Charts | ECharts, low-saturation palette |
| Deploy | go:embed, systemd, scripts/install.sh |
| SSL (Phase 2) | lego ACME client |

## Directory Layout

```
cmd/zpanel/main.go          # server | cli entry
internal/
  app/ config/ auth/ cli/
  handler/ service/ model/ store/
  web/embed.go              # embed web/dist
web/src/                    # Vue SPA
templates/nginx/ systemd/   # config templates
scripts/                    # install.sh, lnmp-install.sh, release.sh
configs/config.example.yaml
```

## API Conventions

- Base: `/api/v1`
- Response: `{ "ok": true|false, "message": "", "data": {} }`
- Auth: JWT Bearer (except `POST /auth/login`)
- All panel actions must have audit log entries

Key endpoints: `/health`, `/auth/login`, `/monitor/overview`, `/monitor/processes`, `/services`, `/lnmp/status`, `/lnmp/install`, `/sites` CRUD.

## UI Design (Ops Console — Mandatory)

- **Layout:** dark sidebar `#18181B` + light content `#F4F4F5`
- **Primary:** `#2563EB` only — no purple/cyan gradients
- **Radius:** 4px everywhere
- **Fonts:** Source Sans 3 (UI), IBM Plex Mono (logs/paths)
- **Cards:** white + 1px border `#E4E4E7`, minimal shadow
- **Responsive:** Desktop ≥1200 / Tablet 768~1199 / Mobile <768 (bottom Tab + Drawer)
- **Forbidden:** emoji, AI gradient backgrounds, glassmorphism, large rounded corners, marketing copy

Theme file: `web/src/styles/theme.ts` — see DEVELOPMENT.md §8.7.

Site type tags: HTML `#52525B`, PHP `#2563EB`, Go `#16A34A` (bordered, low saturation).

## Security Rules (Non-Negotiable)

1. Never pass user input to shell directly — command whitelist only
2. File API: sandbox paths only (`/var/www`, nginx configs, logs) — block `../`
3. Passwords: bcrypt cost ≥ 12
4. Login rate limit: 5 failures / 15 min
5. Never commit `.env`, keys, `configs/config.yaml`, secrets
6. nginx reload only after `nginx -t` passes

## Site Types

| Type | Backend action |
|------|----------------|
| HTML | Create dir → render html.conf.tmpl → nginx reload |
| PHP | + php-fpm socket in nginx conf |
| Go | systemd unit + go-proxy.conf.tmpl reverse proxy |

Site config: SQLite + YAML in `/var/lib/zpanel/`. Model in DEVELOPMENT.md §2.3.

## CLI (cobra)

Binary: `zpanel`, symlink `zp`. No args → interactive menu (like Baota `bt`).

Essential commands: `default`, `start/stop/restart`, `user password`, `port set`, `lnmp install/status`, `site list/add/delete`, `version`.

## Version & Git Workflow

Every milestone: update CHANGELOG + VERSION → commit → tag `vX.Y.Z` → push --tags.

Commit format: Conventional Commits (`feat(scope):`, `fix:`, `docs:`, `chore: release vX.Y.Z`).

Scripts: `scripts/bump-version.sh`, `scripts/release.sh`.

**Commit frequently.** One feature = one commit. User: zex / zex0128@163.com (local repo config).

## Development Commands

```bash
# Backend
cp configs/config.example.yaml configs/config.yaml
go run ./cmd/zpanel server

# Frontend (separate terminal)
cd web && npm run dev    # proxies /api to :8888

# Build
cd web && npm run build
make build               # embed + compile binary
go test ./...
```

## Implementation Guidelines

1. **Minimize scope** — smallest correct diff; match existing patterns
2. **API First** — backend endpoint before UI, or define contract first
3. **No LNMP compile** — use package manager via scripts/lnmp-install.sh
4. **Templates** — Nginx/systemd configs via Go text/template, not string concat
5. **Error handling** — return structured API errors; roll back on nginx/systemd failure
6. **No tests** unless requested or meaningful behavior coverage needed

## Current Sprint

See [docs/PLAN.md §6](../../../docs/PLAN.md#6-当前-sprint下一步). Start v0.2.0: Go scaffold → health → config → SQLite → JWT → monitor API.

## Additional Reference

For path defaults, nginx templates, CLI full command list, and release checklist, see [reference.md](reference.md).
