package config

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
}

type AuthConfig struct {
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
}

type PathsConfig struct {
	WWW           string `yaml:"www"`
	Data          string `yaml:"data"`
	Logs          string `yaml:"logs"`
	NginxSites    string `yaml:"nginx_sites"`
	NginxEnabled  string `yaml:"nginx_enabled"`
}

type FilesConfig struct {
	AllowedPaths  []string `yaml:"allowed_paths"`
	MaxUploadSize int64    `yaml:"max_upload_size"`
}

type DatabaseConfig struct {
	SQLite string `yaml:"sqlite"`
}

// Load 从 YAML 文件加载配置（v0.2.0 下一步实现）
func Load(path string) (*Config, error) {
	_ = path
	return nil, nil
}
