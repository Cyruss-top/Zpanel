package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zex/zpanel/internal/app"
)

func newServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "server",
		Short:  "启动面板 HTTP 服务",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.RunServer(app.Options{Version: Version})
		},
	}
}

func newDefaultCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "default",
		Short: "查看面板入口信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			printPanelURL(cfg)
			return nil
		},
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本号",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
}

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "启动面板服务",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireRoot(); err != nil {
				return err
			}
			return panelSystemctl("start", "zpanel")
		},
	}
}

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "停止面板服务",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireRoot(); err != nil {
				return err
			}
			return panelSystemctl("stop", "zpanel")
		},
	}
}

func newRestartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "重启面板服务",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireRoot(); err != nil {
				return err
			}
			return panelSystemctl("restart", "zpanel")
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "查看面板服务状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			return panelSystemctl("status", "zpanel")
		},
	}
}
