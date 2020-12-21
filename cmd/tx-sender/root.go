package txsender

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tx-sender",
		Short: "Run tx-sender",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
