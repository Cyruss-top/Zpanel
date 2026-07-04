package model

import "time"

// SiteType 站点类型
type SiteType string

const (
	SiteHTML SiteType = "html"
	SitePHP  SiteType = "php"
	SiteGo   SiteType = "go"
)

// SiteStatus 站点状态
type SiteStatus string

const (
	SiteRunning SiteStatus = "running"
	SiteStopped SiteStatus = "stopped"
	SiteError   SiteStatus = "error"
)

// Site 站点模型
type Site struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Type            SiteType   `json:"type"`
	Status          SiteStatus `json:"status"`
	Domains         []string   `json:"domains"`
	Root            string     `json:"root"`
	PHPVersion      string     `json:"php_version,omitempty"`
	GoPort          int        `json:"go_port,omitempty"`
	GoBinary        string     `json:"go_binary,omitempty"`
	SystemdUnit     string     `json:"systemd_unit,omitempty"`
	NginxConfigPath string     `json:"nginx_config_path"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateSiteRequest 创建站点请求
type CreateSiteRequest struct {
	Name       string   `json:"name" binding:"required"`
	Type       SiteType `json:"type" binding:"required"`
	Domains    []string `json:"domains" binding:"required"`
	PHPVersion string   `json:"php_version"`
	GoPort     int      `json:"go_port"`
	GoBinary   string   `json:"go_binary"`
}

// LNMPStatus LNMP 组件状态
type LNMPStatus struct {
	Installed bool            `json:"installed"`
	Platform  string          `json:"platform"`
	Components map[string]ComponentStatus `json:"components"`
}

type ComponentStatus struct {
	Installed bool   `json:"installed"`
	Version   string `json:"version"`
	Running   bool   `json:"running"`
}
