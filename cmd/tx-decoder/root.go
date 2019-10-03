package txdecoder

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tx-decoder",
		Short: "Run tx-decoder",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
