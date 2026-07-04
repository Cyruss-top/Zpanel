# Zpanel Quick Reference

## Path Defaults

| Item | Path |
|------|------|
| Panel port | 8888 |
| Config | /etc/zpanel/config.yaml |
| Data | /var/lib/zpanel/ |
| WWW root | /var/www/ |
| Nginx sites | /etc/nginx/sites-available/ |
| Panel logs | /var/log/zpanel/ |
| Binary | /usr/local/bin/zpanel |

## Milestone Checklist

| Tag | Deliverable |
|-----|-------------|
| v0.2.0 | Go server, JWT, monitor API, systemd API |
| v0.3.0 | Vue layout, login, dashboard, responsive |
| v0.4.0 | LNMP script, site CRUD (HTML/PHP/Go), nginx templates |
| v0.5.0 | cobra CLI, install.sh, cross-compile |
| v1.0.0 | E2E on Ubuntu 22.04, security hardening |

## Release Steps

```bash
# Update CHANGELOG [Unreleased] -> [X.Y.Z]
echo "0.2.0" > VERSION
git add -A && git commit -m "chore: release v0.2.0"
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin main --tags
```

## Naive UI Theme Keys

```typescript
// web/src/styles/theme.ts
common: { primaryColor: '#2563EB', borderRadius: '4px', bodyColor: '#F4F4F5' }
Layout: { siderColor: '#18181B', headerColor: '#FFFFFF' }
Menu: { itemColorActive: '#27272A', itemTextColorActive: '#FAFAFA' }
```

## Allowed Shell Commands (Whitelist Pattern)

```go
"nginx_test":   {"nginx", "-t"}
"nginx_reload": {"systemctl", "reload", "nginx"}
// Never: exec.Command("sh", "-c", userInput)
```

## File Sandbox Paths

```
/var/www
/etc/nginx/sites-available
/etc/nginx/sites-enabled
/var/log/nginx
```

## Frontend Routes

```
/login          /               /sites          /sites/create
/sites/:id      /environment    /database       /files
/cron           /logs           /settings
```

## Phase 2 (Post-MVP, Not v1.0.0)

SSL (lego), file manager, MySQL backup, cron UI, log tail WebSocket, PHP multi-version.

Do not implement Phase 2 features during MVP unless user explicitly requests.
