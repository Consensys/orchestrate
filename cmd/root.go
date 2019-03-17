package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/config"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "app",
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// This is executed before each run (included on children command run)
			config.ConfigureLogger()
		},
	}

	// Set Persistent flags
	config.LogLevel(rootCmd.PersistentFlags())
	config.LogFormat(rootCmd.PersistentFlags())
	config.PGFlags(rootCmd.PersistentFlags())

	// Add Run command
	rootCmd.AddCommand(newRunCommand())

	// Add Migrate command
	rootCmd.AddCommand(mewMigrateCmd())

	return rootCmd
}
