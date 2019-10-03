package txlistener

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tx-listener",
		Short: "Run tx-listener",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
