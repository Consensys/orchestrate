package txcrafter

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tx-crafter",
		Short: "Run tx-crafter",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
