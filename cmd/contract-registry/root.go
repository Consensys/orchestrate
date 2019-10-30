package contractregistry

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "contract-registry",
		Short: "Run contract-registry",
	}

	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newMigrateCmd())

	return rootCmd
}
