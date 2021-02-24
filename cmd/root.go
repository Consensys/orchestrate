package cmd

import (
	"github.com/ConsenSys/orchestrate/cmd/api"
	keymanager "github.com/ConsenSys/orchestrate/cmd/key-manager"
	txlistener "github.com/ConsenSys/orchestrate/cmd/tx-listener"
	txsender "github.com/ConsenSys/orchestrate/cmd/tx-sender"
	"github.com/ConsenSys/orchestrate/cmd/utils"
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
	rootCmd.AddCommand(keymanager.NewRootCommand())
	rootCmd.AddCommand(utils.NewRootCommand())

	return rootCmd
}
