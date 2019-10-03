package txnonce

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tx-nonce",
		Short: "Run tx-nonce",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
