package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/zex/zpanel/internal/config"
)

// Version CLI 版本号（由 main 注入）
var Version = "dev"

func configPath() string {
	return config.ResolvePath()
}

func loadConfig() (*config.Config, error) {
	return config.Load(configPath())
}

func saveConfig(cfg *config.Config) error {
	return config.Save(configPath(), cfg)
}

func panelSystemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func serverIP() string {
	if out, err := exec.Command("curl", "-s", "--connect-timeout", "3", "ip.sb").Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("hostname", "-I").Output(); err == nil {
		fields := strings.Fields(string(out))
		if len(fields) > 0 {
			return fields[0]
		}
	}
	return "127.0.0.1"
}

func printPanelURL(cfg *config.Config) {
	ip := serverIP()
	fmt.Printf("面板地址: http://%s:%d\n", ip, cfg.Panel.Port)
	fmt.Printf("用户名:   %s\n", cfg.Auth.Username)
	fmt.Printf("配置文件: %s\n", configPath())
}

func requireRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("此操作需要 root 权限")
	}
	return nil
}
