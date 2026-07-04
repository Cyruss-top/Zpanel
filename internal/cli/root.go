package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "zpanel",
		Short:         "Zpanel Linux 面板管理工具",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       Version,
		Run: func(cmd *cobra.Command, args []string) {
			runInteractiveMenu()
		},
	}
	root.AddCommand(
		newServerCmd(),
		newDefaultCmd(),
		newVersionCmd(),
		newStartCmd(),
		newStopCmd(),
		newRestartCmd(),
		newStatusCmd(),
		newUserCmd(),
		newPortCmd(),
		newEntryCmd(),
		newLNMPCmd(),
		newSiteCmd(),
	)
	return root
}

// Execute 运行 CLI
func Execute(version string) {
	Version = version
	root := newRootCmd()
	root.Version = version
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "zpanel: %v\n", err)
		os.Exit(1)
	}
}
