//go:build linux

package site

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zex/zpanel/internal/config"
	"github.com/zex/zpanel/internal/model"
	nx "github.com/zex/zpanel/internal/service/nginx"
	"github.com/zex/zpanel/internal/store"
)

// Service 站点管理
type Service struct {
	cfg   *config.Config
	store *store.Store
}

func NewService(cfg *config.Config, st *store.Store) *Service {
	return &Service{cfg: cfg, store: st}
}

func (s *Service) List() ([]model.Site, error) {
	return s.store.ListSites()
}

func (s *Service) Get(id string) (*model.Site, error) {
	return s.store.GetSite(id)
}

func (s *Service) Create(req model.CreateSiteRequest) (*model.Site, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	id := newID()
	name := sanitizeName(req.Name)
	root := filepath.Join(s.cfg.Paths.WWW, name)
	confName := name + ".conf"
	confPath := filepath.Join(s.cfg.Paths.NginxSites, confName)

	site := &model.Site{
		ID:              id,
		Name:            name,
		Type:            req.Type,
		Status:          model.SiteRunning,
		Domains:         req.Domains,
		Root:            root,
		PHPVersion:      req.PHPVersion,
		GoPort:          req.GoPort,
		GoBinary:        req.GoBinary,
		NginxConfigPath: confPath,
	}

	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}

	tmplData := nx.TemplateData{
		Name:    name,
		Domains: req.Domains,
		Root:    root,
	}

	switch req.Type {
	case model.SiteHTML:
		if err := writeDefaultHTML(root); err != nil {
			return nil, err
		}
		content, err := nx.Render(resolveTemplate("nginx/html.conf.tmpl"), tmplData)
		if err != nil {
			return nil, err
		}
		if err := nx.WriteConfig(confPath, content); err != nil {
			return nil, err
		}
	case model.SitePHP:
		if site.PHPVersion == "" {
			site.PHPVersion = detectPHPVersion()
		}
		tmplData.PHPVersion = site.PHPVersion
		if err := writeDefaultPHP(root); err != nil {
			return nil, err
		}
		content, err := nx.Render(resolveTemplate("nginx/php.conf.tmpl"), tmplData)
		if err != nil {
			return nil, err
		}
		if err := nx.WriteConfig(confPath, content); err != nil {
			return nil, err
		}
	case model.SiteGo:
		if site.GoPort == 0 {
			site.GoPort = 8080
		}
		if site.GoBinary == "" {
			site.GoBinary = filepath.Join(root, "app")
		}
		site.SystemdUnit = "zpanel-site-" + name + ".service"
		if err := s.createGoUnit(site); err != nil {
			return nil, err
		}
		tmplData.GoPort = site.GoPort
		content, err := nx.Render(resolveTemplate("nginx/go-proxy.conf.tmpl"), tmplData)
		if err != nil {
			return nil, err
		}
		if err := nx.WriteConfig(confPath, content); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported site type")
	}

	if err := nx.EnableSite(s.cfg.Paths.NginxSites, s.cfg.Paths.NginxEnabled, name); err != nil {
		return nil, err
	}
	if err := nx.Reload(); err != nil {
		return nil, err
	}

	if err := s.store.InsertSite(site); err != nil {
		return nil, err
	}
	return site, nil
}

func (s *Service) Delete(id string) error {
	site, err := s.store.GetSite(id)
	if err != nil || site == nil {
		return fmt.Errorf("site not found")
	}

	_ = os.Remove(filepath.Join(s.cfg.Paths.NginxEnabled, filepath.Base(site.NginxConfigPath)))
	_ = os.Remove(site.NginxConfigPath)

	if site.SystemdUnit != "" {
		_ = runSystemctl("stop", site.SystemdUnit)
		_ = os.Remove("/etc/systemd/system/" + site.SystemdUnit)
		_ = runSystemctl("daemon-reload")
	}

	_ = nx.Reload()
	return s.store.DeleteSite(id)
}

func (s *Service) createGoUnit(site *model.Site) error {
	content, err := nx.Render(resolveTemplate("systemd/go-site.service.tmpl"), nx.TemplateData{
		Name:   site.Name,
		Root:   site.Root,
		Binary: site.GoBinary,
		GoPort: site.GoPort,
	})
	if err != nil {
		return err
	}
	unitPath := "/etc/systemd/system/" + site.SystemdUnit
	if err := os.WriteFile(unitPath, []byte(content), 0o644); err != nil {
		return err
	}
	if err := runSystemctl("daemon-reload"); err != nil {
		return err
	}
	if err := runSystemctl("enable", "--now", site.SystemdUnit); err != nil {
		return err
	}
	return nil
}

func runSystemctl(args ...string) error {
	return nx.RunSystemctl(args...)
}

func validateCreate(req model.CreateSiteRequest) error {
	if len(req.Domains) == 0 {
		return fmt.Errorf("domains required")
	}
	switch req.Type {
	case model.SiteHTML, model.SitePHP, model.SiteGo:
		return nil
	default:
		return fmt.Errorf("invalid site type")
	}
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func sanitizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "/", "")
	return name
}

func resolveTemplate(rel string) string {
	candidates := []string{
		filepath.Join("templates", rel),
		filepath.Join("/usr/local/zpanel/templates", rel),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return filepath.Join("templates", rel)
}

func detectPHPVersion() string {
	// 默认 8.2，安装脚本会装系统默认 PHP
	return "8.2"
}

func writeDefaultHTML(root string) error {
	p := filepath.Join(root, "index.html")
	if _, err := os.Stat(p); err == nil {
		return nil
	}
	return os.WriteFile(p, []byte("<!DOCTYPE html><html><body><h1>Zpanel Site</h1></body></html>"), 0o644)
}

func writeDefaultPHP(root string) error {
	p := filepath.Join(root, "index.php")
	if _, err := os.Stat(p); err == nil {
		return nil
	}
	return os.WriteFile(p, []byte("<?php phpinfo();"), 0o644)
}
