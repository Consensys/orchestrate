package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-crafter"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-sender"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-signer"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/logger"
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

	// Add Run command
	rootCmd.AddCommand(txcrafter.NewRunCommand())
	rootCmd.AddCommand(txnonce.NewRunCommand())
	rootCmd.AddCommand(txsigner.NewRunCommand())
	rootCmd.AddCommand(txsender.NewRunCommand())
	rootCmd.AddCommand(txlistener.NewRunCommand())
	rootCmd.AddCommand(txdecoder.NewRunCommand())
	rootCmd.AddCommand(contractregistry.NewRunCommand())
	rootCmd.AddCommand(envelopestore.NewRunCommand())

	return rootCmd
}
