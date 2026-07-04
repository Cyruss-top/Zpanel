# Gitee 自动同步指南

主仓库（GitHub）：[Cyruss-top/Zpanel](https://github.com/Cyruss-top/Zpanel)  
国内镜像（Gitee）：[Ressss2023/Zpanel](https://gitee.com/Ressss2023/Zpanel)

---

## 方式一：Gitee 官方镜像（最简单）

适合：**代码 + 标签** 从 GitHub 同步到 Gitee。

### 步骤

1. 登录 Gitee，打开仓库 **Ressss2023/Zpanel**
2. 进入 **管理** → **仓库镜像管理** → **添加镜像**
3. 镜像方向选择 **Pull**（从 GitHub 拉取到 Gitee）
4. 填写 GitHub 仓库地址：
   ```
   https://github.com/Cyruss-top/Zpanel
   ```
5. 勾选 **「自动从 GitHub 同步仓库」**
6. 按提示配置 **Gitee 私人令牌**（需 `projects` 和 `admin:repo_hook` 权限）

### 触发同步

| 方式 | 说明 |
|------|------|
| 自动 | GitHub 有新提交时，Gitee 通过 Webhook 自动拉取 |
| 手动 | 镜像管理页点击 **「更新镜像」** |

### 限制

- **不同步 GitHub Releases 附件**（安装包需单独上传到 Gitee Releases）
- 官方镜像偶尔不稳定，若失效请用方式二

官方文档：[Gitee ↔ GitHub 仓库镜像](https://help.gitee.com/repository/settings/sync-between-gitee-github)

---

## 方式二：GitHub Actions 自动同步（推荐）

项目已内置工作流：`.github/workflows/sync-gitee.yml`

推送 `main` 分支或 `v*` 标签到 GitHub 后，自动同步到 Gitee。

### 一次性配置

#### 1. 生成 Gitee 私人令牌

1. 打开 [Gitee 私人令牌](https://gitee.com/profile/personal_access_tokens)
2. 新建令牌，勾选 **`projects`** 权限
3. 复制令牌（只显示一次）

#### 2. 在 GitHub 配置 Secret

1. 打开 [GitHub 仓库 Settings → Secrets](https://github.com/Cyruss-top/Zpanel/settings/secrets/actions)
2. 点击 **New repository secret**
3. Name：`GITEE_TOKEN`
4. Value：粘贴 Gitee 令牌

#### 3. 确认 Gitee 仓库已存在

若尚未创建，先在 Gitee 新建空仓库 `Ressss2023/Zpanel`（不要初始化 README）。

### 验证

配置完成后，向 GitHub 推送任意 commit，在 **Actions → Sync to Gitee** 查看是否成功。

也可手动触发：**Actions → Sync to Gitee → Run workflow**

---

## Releases 安装包同步

**代码同步不会自动带上 Release 二进制包**，国内一键安装需要 Gitee 也有 Release：

### 选项 A：GitHub Actions 发布（GitHub 侧）

1. [Actions → Release](https://github.com/Cyruss-top/Zpanel/actions/workflows/release.yml) → Run workflow
2. 填写 `v0.6.0`，生成 `zpanel-linux-amd64.tar.gz` 等

### 选项 B：手动上传到 Gitee

1. 打开 [Gitee Releases](https://gitee.com/Ressss2023/Zpanel/releases)
2. 新建版本（与 GitHub tag 一致，如 `v0.6.0`）
3. 上传 `zpanel-linux-amd64.tar.gz` 和 `zpanel-linux-arm64.tar.gz`

---

## 常见问题

### Gitee 代码不是最新？

1. Gitee 镜像管理 → 点击 **「更新镜像」**
2. 或检查 GitHub Actions `Sync to Gitee` 是否失败
3. 确认 `GITEE_TOKEN` 未过期

### 安装脚本下载失败？

Gitee 需有对应版本的 **Release 附件**，仅同步代码不够。执行：

```bash
wget -qO- https://gitee.com/Ressss2023/Zpanel/raw/main/scripts/install.sh | bash -s -- --mirror gitee --interactive
```

若 Gitee 无 Release，脚本会自动回退尝试 GitHub。

### 两个账号名不同有影响吗？

无影响。GitHub `Cyruss-top` 与 Gitee `Ressss2023` 是不同账号，通过镜像 URL 或 Actions 指定目标仓库即可。
