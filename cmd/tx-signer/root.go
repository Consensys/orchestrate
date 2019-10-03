package txsigner

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tx-signer",
		Short: "Run tx-signer",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
