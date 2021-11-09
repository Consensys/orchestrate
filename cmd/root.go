package cmd

import (
	"github.com/consensys/orchestrate/cmd/api"
	txlistener "github.com/consensys/orchestrate/cmd/tx-listener"
	txsender "github.com/consensys/orchestrate/cmd/tx-sender"
	"github.com/spf13/cobra"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "orchestrate",
		TraverseChildren: true,
		SilenceUsage:     true,
	}

	// Add Run command
	rootCmd.AddCommand(txsender.NewRootCommand())
	rootCmd.AddCommand(txlistener.NewRootCommand())
	rootCmd.AddCommand(api.NewRootCommand())

	return rootCmd
}
