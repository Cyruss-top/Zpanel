package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zex/zpanel/internal/auth"
	"github.com/zex/zpanel/internal/store"
)

func newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "管理员账号管理",
	}
	cmd.AddCommand(newUserShowCmd(), newUserPasswordCmd())
	return cmd
}

func newUserShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "显示当前管理员",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			fmt.Printf("用户名: %s\n", cfg.Auth.Username)
			return nil
		},
	}
}

func newUserPasswordCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "password [新密码]",
		Short: "修改管理员密码",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			password := ""
			if len(args) > 0 {
				password = args[0]
			}
			if password == "" {
				return fmt.Errorf("请提供新密码: zpanel user password <密码>")
			}
			hash, err := auth.HashPassword(password)
			if err != nil {
				return err
			}
			st, err := store.Open(cfg.SQLitePath())
			if err != nil {
				return err
			}
			defer st.Close()
			if err := st.UpsertUser(cfg.Auth.Username, hash); err != nil {
				return err
			}
			fmt.Println("密码已更新")
			return nil
		},
	}
}
