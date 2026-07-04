package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func newUninstallCmd() *cobra.Command {
	var yes, purge bool
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "卸载 Zpanel 面板",
		Long: `卸载面板服务与二进制。默认保留配置、数据库与 /var/www 网站数据。
使用 --purge 可同时删除 /etc/zpanel 与 /var/lib/zpanel。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireRoot(); err != nil {
				return err
			}
			if !yes {
				fmt.Println("即将卸载 Zpanel 面板")
				if purge {
					fmt.Println("将同时删除配置与数据库")
				} else {
					fmt.Println("将保留配置与数据库，仅删除服务与二进制")
				}
				fmt.Println("网站目录 /var/www 不会删除")
				fmt.Print("确认卸载? [y/N] ")
				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				ans := strings.TrimSpace(strings.ToLower(line))
				if ans != "y" && ans != "yes" {
					fmt.Println("已取消")
					return nil
				}
			}
			return runUninstall(purge)
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "跳过确认")
	cmd.Flags().BoolVar(&purge, "purge", false, "同时删除配置与数据")
	return cmd
}

func runUninstall(purge bool) error {
	script := resolveUninstallScript()
	if script != "" {
		args := []string{script}
		if purge {
			args = append(args, "--yes", "--purge")
		} else {
			args = append(args, "--yes")
		}
		cmd := exec.Command("bash", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return uninstallInline(purge)
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

func uninstallInline(purge bool) error {
	_ = panelSystemctl("stop", "zpanel")
	_ = panelSystemctl("disable", "zpanel")

	paths := []string{
		"/etc/systemd/system/zpanel.service",
		"/usr/local/bin/zpanel",
		"/usr/bin/zp",
	}
	for _, p := range paths {
		_ = os.Remove(p)
	}
	_ = os.RemoveAll("/usr/local/zpanel")

	if purge {
		for _, p := range []string{"/etc/zpanel", "/var/lib/zpanel", "/var/log/zpanel"} {
			_ = os.RemoveAll(p)
		}
	}

	_ = panelSystemctl("daemon-reload")
	fmt.Println("Zpanel 已卸载完成")
	if !purge {
		fmt.Println("配置保留于: /etc/zpanel")
	}
	fmt.Println("网站数据保留于: /var/www")
	return nil
}
