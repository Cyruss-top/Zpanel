package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newPortCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "port",
		Short: "面板端口管理",
	}
	cmd.AddCommand(newPortShowCmd(), newPortSetCmd())
	return cmd
}

func newPortShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "显示当前端口",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			fmt.Printf("端口: %d\n", cfg.Panel.Port)
			return nil
		},
	}
}

func newPortSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <port>",
		Short: "设置面板端口",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.Atoi(args[0])
			if err != nil || port < 1 || port > 65535 {
				return fmt.Errorf("无效端口: %s", args[0])
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			cfg.Panel.Port = port
			if err := saveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("端口已设置为 %d，请执行 zpanel restart 生效\n", port)
			return nil
		},
	}
}
