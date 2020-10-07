package identitymanager

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "identity-manager",
		Short: "Run identity-manager",
	}

	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newMigrateCmd())

	return rootCmd
}
