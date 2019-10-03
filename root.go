package main

import (
	"github.com/spf13/cobra"

	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/contract-registry"
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/envelope-store"
	txcrafter "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-crafter"
	txdecoder "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-decoder"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-listener"
	txnonce "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-nonce"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-sender"
	txsigner "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-signer"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/tracing/opentracing/jaeger"
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

	// Register HTTP server flags
	http.Hostname(rootCmd.PersistentFlags())

	// Register OpenTracing flags
	jaeger.InitFlags(rootCmd.PersistentFlags())

	// Add Run command
	rootCmd.AddCommand(txcrafter.NewRootCommand())
	rootCmd.AddCommand(txnonce.NewRootCommand())
	rootCmd.AddCommand(txsigner.NewRootCommand())
	rootCmd.AddCommand(txsender.NewRootCommand())
	rootCmd.AddCommand(txlistener.NewRootCommand())
	rootCmd.AddCommand(txdecoder.NewRootCommand())
	rootCmd.AddCommand(contractregistry.NewRootCommand())
	rootCmd.AddCommand(envelopestore.NewRootCommand())

	return rootCmd
}
