package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zex/zpanel/internal/service/lnmp"
)

func newLNMPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lnmp",
		Short: "LNMP 环境管理",
	}
	cmd.AddCommand(newLNMPStatusCmd(), newLNMPInstallCmd())
	return cmd
}

func newLNMPStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "查看 LNMP 状态",
		Run: func(cmd *cobra.Command, args []string) {
			svc := lnmp.NewService()
			st := svc.Status()
			fmt.Printf("平台: %s  已安装: %v\n", st.Platform, st.Installed)
			for name, comp := range st.Components {
				state := "未安装"
				if comp.Installed {
					if comp.Running {
						state = "运行中"
					} else {
						state = "已停止"
					}
				}
				fmt.Printf("  %-8s %s  %s\n", name+":", comp.Version, state)
			}
		},
	}
}

func newLNMPInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "一键安装 LNMP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireRoot(); err != nil {
				return err
			}
			svc := lnmp.NewService()
			res, err := svc.Install()
			if err != nil {
				return err
			}
			fmt.Printf("LNMP 安装成功\n  nginx: %s\n  php: %s\n  mysql: %s\n", res.Nginx, res.PHP, res.MySQL)
			return nil
		},
	}
}
