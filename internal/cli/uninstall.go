package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	uninstallKeepWWW = "keep-www"
	uninstallAll     = "all"
)

func newUninstallCmd() *cobra.Command {
	var yes, keepWWW, purgeAll, purgeLegacy bool
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "卸载 Zpanel 面板",
		Long: `卸载 Zpanel 面板，支持两种模式：

  保留站点数据 (--keep-www)：删除面板程序、配置、数据库，保留 /var/www
  彻底删除干净 (--all)：删除面板及 /var/www 全部数据，不可恢复

无参数时将进入交互选择。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireRoot(); err != nil {
				return err
			}
			mode := ""
			if keepWWW {
				mode = uninstallKeepWWW
			} else if purgeAll || purgeLegacy {
				mode = uninstallAll
			}
			if mode == "" {
				selected, err := promptUninstallMode(yes)
				if err != nil {
					return err
				}
				if selected == "" {
					fmt.Println("已取消")
					return nil
				}
				mode = selected
			} else if !yes {
				selected, err := confirmUninstallMode(mode)
				if err != nil {
					return err
				}
				if selected == "" {
					fmt.Println("已取消")
					return nil
				}
			}
			return runUninstall(mode, yes)
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "跳过确认")
	cmd.Flags().BoolVar(&keepWWW, "keep-www", false, "保留 /var/www 网站数据")
	cmd.Flags().BoolVar(&purgeAll, "all", false, "彻底删除干净（含 /var/www）")
	cmd.Flags().BoolVar(&purgeLegacy, "purge", false, "同 --all（兼容）")
	return cmd
}

func promptUninstallMode(skipConfirm bool) (string, error) {
	fmt.Println("")
	fmt.Println("============================================")
	fmt.Println("  Zpanel 卸载")
	fmt.Println("============================================")
	fmt.Println("  1. 保留站点数据")
	fmt.Println("     删除面板程序、配置、数据库、日志")
	fmt.Println("     保留网站目录: /var/www")
	fmt.Println("")
	fmt.Println("  2. 彻底删除干净")
	fmt.Println("     删除面板及 /var/www 全部网站文件")
	fmt.Println("     此操作不可恢复！")
	fmt.Println("============================================")
	fmt.Print("请选择 [1/2]: ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	switch strings.TrimSpace(line) {
	case "1":
		if skipConfirm {
			return uninstallKeepWWW, nil
		}
		return confirmUninstallMode(uninstallKeepWWW)
	case "2":
		if skipConfirm {
			return uninstallAll, nil
		}
		return confirmUninstallMode(uninstallAll)
	default:
		return "", nil
	}
}

func confirmUninstallMode(mode string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	if mode == uninstallAll {
		fmt.Println("")
		fmt.Println("即将彻底删除 Zpanel 及 /var/www 全部数据")
		fmt.Print("确认彻底删除? 输入 yes: ")
		line, _ := reader.ReadString('\n')
		if strings.TrimSpace(line) != "yes" {
			return "", nil
		}
		return uninstallAll, nil
	}
	fmt.Println("")
	fmt.Println("即将卸载 Zpanel，保留 /var/www 网站数据")
	fmt.Print("确认卸载? [y/N] ")
	line, _ := reader.ReadString('\n')
	ans := strings.TrimSpace(strings.ToLower(line))
	if ans != "y" && ans != "yes" {
		return "", nil
	}
	return uninstallKeepWWW, nil
}

func runUninstall(mode string, skipConfirm bool) error {
	script := resolveUninstallScript()
	if script != "" {
		args := []string{script}
		if mode == uninstallAll {
			args = append(args, "--all")
		} else {
			args = append(args, "--keep-www")
		}
		if skipConfirm {
			args = append(args, "--yes")
		}
		cmd := exec.Command("bash", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return uninstallInline(mode)
}

func resolveUninstallScript() string {
	for _, p := range []string{
		"/usr/local/zpanel/scripts/uninstall.sh",
		"scripts/uninstall.sh",
	} {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}

func uninstallInline(mode string) error {
	_ = panelSystemctl("stop", "zpanel")
	_ = panelSystemctl("disable", "zpanel")

	for _, p := range []string{
		"/etc/systemd/system/zpanel.service",
		"/usr/local/bin/zpanel",
		"/usr/bin/zp",
	} {
		_ = os.Remove(p)
	}
	_ = os.RemoveAll("/usr/local/zpanel")
	for _, p := range []string{"/etc/zpanel", "/var/lib/zpanel", "/var/log/zpanel"} {
		_ = os.RemoveAll(p)
	}
	if mode == uninstallAll {
		_ = os.RemoveAll("/var/www")
	}
	_ = panelSystemctl("daemon-reload")

	if mode == uninstallAll {
		fmt.Println("Zpanel 已彻底卸载，所有数据已删除")
	} else {
		fmt.Println("Zpanel 已卸载完成")
		fmt.Println("网站数据保留于: /var/www")
	}
	return nil
}

// RunUninstallInteractive 交互菜单卸载入口
func RunUninstallInteractive() error {
	if err := requireRoot(); err != nil {
		return err
	}
	mode, err := promptUninstallMode(false)
	if err != nil {
		return err
	}
	if mode == "" {
		fmt.Println("已取消")
		return nil
	}
	return runUninstall(mode, true)
}
