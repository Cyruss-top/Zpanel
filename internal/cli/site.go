package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zex/zpanel/internal/service/site"
	"github.com/zex/zpanel/internal/store"
)

func newSiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "站点管理",
	}
	cmd.AddCommand(newSiteListCmd(), newSiteDeleteCmd())
	return cmd
}

func newSiteListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有站点",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			st, err := store.Open(cfg.SQLitePath())
			if err != nil {
				return err
			}
			defer st.Close()
			sites, err := site.NewService(cfg, st).List()
			if err != nil {
				return err
			}
			if len(sites) == 0 {
				fmt.Println("暂无站点")
				return nil
			}
			for _, s := range sites {
				fmt.Printf("%s  %-6s  %-8s  %s\n", s.ID[:8], s.Type, s.Status, s.Name)
			}
			return nil
		},
	}
}

func newSiteDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "删除站点",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			st, err := store.Open(cfg.SQLitePath())
			if err != nil {
				return err
			}
			defer st.Close()
			svc := site.NewService(cfg, st)
			sites, err := svc.List()
			if err != nil {
				return err
			}
			name := args[0]
			for _, s := range sites {
				if s.Name == name {
					if err := svc.Delete(s.ID); err != nil {
						return err
					}
					fmt.Printf("已删除站点: %s\n", name)
					return nil
				}
			}
			return fmt.Errorf("站点不存在: %s", name)
		},
	}
}
