package chainregistry

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "chain-registry",
		Short: "Run chain-registry",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
