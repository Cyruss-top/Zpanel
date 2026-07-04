package lnmp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zex/zpanel/internal/model"
)

// Service LNMP 环境管理
type Service struct {
	scriptPath string
}

func NewService() *Service {
	return &Service{scriptPath: resolveScript("lnmp-install.sh")}
}

func resolveScript(name string) string {
	candidates := []string{
		filepath.Join("scripts", name),
		filepath.Join("/usr/local/zpanel/scripts", name),
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "scripts", name))
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return filepath.Join("scripts", name)
}

// Status 检测 LNMP 状态
func (s *Service) Status() *model.LNMPStatus {
	st := &model.LNMPStatus{
		Platform:   runtime.GOOS,
		Components: map[string]model.ComponentStatus{},
	}
	if runtime.GOOS != "linux" {
		st.Components["nginx"] = model.ComponentStatus{}
		st.Components["mysql"] = model.ComponentStatus{}
		st.Components["php"] = model.ComponentStatus{}
		return st
	}

	nginx := probeComponent("nginx", []string{"nginx", "-v"}, "nginx")
	mysql := probeComponent("mysql", []string{"mysql", "--version"}, "mysql", "mysqld")
	php := probeComponent("php", []string{"php", "-v"}, "php-fpm")
	phpVer := ""
	if out, err := exec.Command("php", "-r", `echo PHP_MAJOR_VERSION.".".PHP_MINOR_VERSION;`).Output(); err == nil {
		phpVer = strings.TrimSpace(string(out))
		php.Version = phpVer
	}

	st.Components["nginx"] = nginx
	st.Components["mysql"] = mysql
	st.Components["php"] = php
	st.Installed = nginx.Installed && mysql.Installed && php.Installed
	return st
}

func probeComponent(name string, versionCmd []string, units ...string) model.ComponentStatus {
	st := model.ComponentStatus{}
	if out, err := exec.Command(versionCmd[0], versionCmd[1:]...).CombinedOutput(); err == nil {
		st.Installed = true
		st.Version = parseVersion(name, string(out))
	}
	for _, unit := range units {
		if active, _ := isActive(unit); active {
			st.Running = true
			break
		}
	}
	return st
}

func parseVersion(name, out string) string {
	out = strings.TrimSpace(out)
	switch name {
	case "nginx":
		if i := strings.Index(out, "/"); i >= 0 {
			return strings.TrimSpace(out[i+1:])
		}
	case "mysql":
		fields := strings.Fields(out)
		if len(fields) >= 3 {
			return strings.Trim(fields[2], ",")
		}
	case "php":
		if len(out) > 0 {
			return strings.Fields(out)[1]
		}
	}
	return out
}

func isActive(unit string) (bool, error) {
	out, err := exec.Command("systemctl", "is-active", unit).Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "active", nil
}

type installResult struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Nginx   string `json:"nginx"`
	PHP     string `json:"php"`
	MySQL   string `json:"mysql"`
}

// Install 执行 LNMP 安装脚本
func (s *Service) Install() (*installResult, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("LNMP install requires Linux")
	}
	if _, err := os.Stat(s.scriptPath); err != nil {
		return nil, fmt.Errorf("install script not found: %s", s.scriptPath)
	}
	cmd := exec.Command("bash", s.scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("install failed: %s: %w", string(out), err)
	}
	line := strings.TrimSpace(lastLine(string(out)))
	var res installResult
	if err := json.Unmarshal([]byte(line), &res); err != nil {
		return nil, fmt.Errorf("parse install result: %w (output: %s)", err, line)
	}
	if !res.OK {
		return &res, fmt.Errorf(res.Message)
	}
	return &res, nil
}

func lastLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	return lines[len(lines)-1]
}
