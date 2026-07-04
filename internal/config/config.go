package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultPort = 8888
	DefaultBind = "127.0.0.1"
)

// Config 面板主配置，对应 configs/config.yaml
type Config struct {
	Panel    PanelConfig    `yaml:"panel"`
	Auth     AuthConfig     `yaml:"auth"`
	Paths    PathsConfig    `yaml:"paths"`
	Files    FilesConfig    `yaml:"files"`
	Database DatabaseConfig `yaml:"database"`
}

type PanelConfig struct {
	Port int    `yaml:"port"`
	Bind string `yaml:"bind"`
	SSL  bool   `yaml:"ssl"`
	Entry string `yaml:"entry"` // 安全入口后缀，如 abc123 → /abc123/
}

type AuthConfig struct {
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
}

type PathsConfig struct {
	WWW          string `yaml:"www"`
	Data         string `yaml:"data"`
	Logs         string `yaml:"logs"`
	NginxSites   string `yaml:"nginx_sites"`
	NginxEnabled string `yaml:"nginx_enabled"`
}

type FilesConfig struct {
	AllowedPaths  []string `yaml:"allowed_paths"`
	MaxUploadSize int64    `yaml:"max_upload_size"`
}

type DatabaseConfig struct {
	SQLite string `yaml:"sqlite"`
}

// ResolvePath 解析配置文件路径：环境变量 > 开发路径 > 生产路径
func ResolvePath() string {
	if p := os.Getenv("ZPANEL_CONFIG"); p != "" {
		return p
	}
	if _, err := os.Stat("configs/config.yaml"); err == nil {
		return "configs/config.yaml"
	}
	return "/etc/zpanel/config.yaml"
}

// Load 从 YAML 文件加载配置并填充默认值
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if c.Panel.Port == 0 {
		c.Panel.Port = DefaultPort
	}
	if c.Panel.Bind == "" {
		c.Panel.Bind = DefaultBind
	}
	if c.Auth.Username == "" {
		c.Auth.Username = "admin"
	}
	if c.Paths.WWW == "" {
		c.Paths.WWW = "./data/www"
	}
	if c.Paths.Data == "" {
		c.Paths.Data = "./data/lib"
	}
	if c.Paths.Logs == "" {
		c.Paths.Logs = "./data/logs"
	}
	if c.Paths.NginxSites == "" {
		c.Paths.NginxSites = "./data/nginx/sites-available"
	}
	if c.Paths.NginxEnabled == "" {
		c.Paths.NginxEnabled = "./data/nginx/sites-enabled"
	}
	if c.Database.SQLite == "" {
		c.Database.SQLite = "zpanel.db"
	}
	if c.Files.MaxUploadSize == 0 {
		c.Files.MaxUploadSize = 50 * 1024 * 1024
	}
}

// EnsureDirs 创建运行所需目录
func (c *Config) EnsureDirs() error {
	dirs := []string{
		c.Paths.WWW, c.Paths.Data, c.Paths.Logs,
		c.Paths.NginxSites, c.Paths.NginxEnabled,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}
	return nil
}

// SQLitePath 返回 SQLite 数据库绝对路径
func (c *Config) SQLitePath() string {
	if filepath.IsAbs(c.Database.SQLite) {
		return c.Database.SQLite
	}
	return filepath.Join(c.Paths.Data, c.Database.SQLite)
}

// ListenAddr 返回 HTTP 监听地址
func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Panel.Bind, c.Panel.Port)
}

// EntryPrefix 返回安全入口路径前缀，如 /abc123；未设置返回空字符串
func (c *Config) EntryPrefix() string {
	e := strings.TrimSpace(c.Panel.Entry)
	e = strings.Trim(e, "/")
	if e == "" {
		return ""
	}
	return "/" + e
}

// NormalizeEntry 校验并规范化入口后缀
func NormalizeEntry(entry string) (string, error) {
	entry = strings.TrimSpace(entry)
	entry = strings.Trim(entry, "/")
	if entry == "" {
		return "", nil
	}
	if len(entry) < 3 || len(entry) > 32 {
		return "", fmt.Errorf("入口后缀长度须为 3~32 个字符")
	}
	for _, r := range entry {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return "", fmt.Errorf("入口后缀仅允许字母、数字、-、_")
	}
	return entry, nil
}

// Save 写入配置文件
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
