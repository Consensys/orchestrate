package cmd

import (
	"github.com/spf13/cobra"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/chain-registry"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/contract-registry"
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/envelope-store"
	txcrafter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-crafter"
	txdecoder "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-decoder"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-listener"
	txnonce "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-nonce"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-sender"
	txsigner "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-signer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "orchestrate",
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
	metrics.Flags(rootCmd.PersistentFlags())

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
	rootCmd.AddCommand(chainregistry.NewRootCommand())
	rootCmd.AddCommand(utils.NewRootCommand())

	// Register Multi-Tenancy flags
	multitenancy.Enabled(rootCmd.PersistentFlags())
	authjwt.Flags(rootCmd.PersistentFlags())
	authkey.APIKey(rootCmd.PersistentFlags())

	return rootCmd
}
