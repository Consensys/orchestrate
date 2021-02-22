package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/api"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/key-manager"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/tx-listener"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/tx-sender"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/utils"
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
