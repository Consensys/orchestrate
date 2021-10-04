package api

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "api",
		Short: "Run api",
	}

	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newMigrateCmd())
	rootCmd.AddCommand(newAccountCmd())

	return rootCmd
}
