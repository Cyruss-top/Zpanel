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

func requireRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("此操作需要 root 权限")
	}
	return nil
}
