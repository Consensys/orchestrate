package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "app",
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// This is executed before each run (included on children command run)
			logger.InitLogger()
		},
	}

	// Set Persistent flags
	logger.LogLevel(rootCmd.PersistentFlags())
	logger.LogFormat(rootCmd.PersistentFlags())

	// Register Opentracing flags
	jaeger.InitFlags(rootCmd.PersistentFlags())

	// Register HTTP server flags
	http.Hostname(rootCmd.PersistentFlags())

	// Add Run command
	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
