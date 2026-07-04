package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zex/zpanel/internal/config"
)

func newEntryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entry",
		Short: "安全入口后缀管理",
	}
	cmd.AddCommand(newEntryShowCmd(), newEntrySetCmd())
	return cmd
}

func newEntryShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "显示当前安全入口后缀",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			if cfg.Panel.Entry == "" {
				fmt.Println("未设置安全入口（直接通过根路径访问）")
			} else {
				fmt.Printf("安全入口: /%s/\n", cfg.Panel.Entry)
			}
			return nil
		},
	}
}

func newEntrySetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <后缀>",
		Short: "设置安全入口后缀（留空用 clear 清除）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			val := args[0]
			if val == "clear" || val == "none" || val == "-" {
				val = ""
			}
			normalized, err := config.NormalizeEntry(val)
			if err != nil {
				return err
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			cfg.Panel.Entry = normalized
			if err := saveConfig(cfg); err != nil {
				return err
			}
			if normalized == "" {
				fmt.Println("安全入口已清除，请执行 zpanel restart 生效")
			} else {
				fmt.Printf("安全入口已设置为 /%s/ ，请执行 zpanel restart 生效\n", normalized)
			}
			return nil
		},
	}
}
