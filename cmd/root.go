package cmd

import (
	"github.com/spf13/cobra"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/chain-registry"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/contract-registry"
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/envelope-store"
	transactionscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/transaction-scheduler"
	txcrafter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-crafter"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-listener"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-sender"
	txsigner "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/tx-signer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/cmd/utils"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	pkglog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "orchestrate",
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// This is executed before each run (included on children command run)
			pkglog.InitLogger()
		},
	}

	// Set Persistent flags
	pkglog.LogLevel(rootCmd.PersistentFlags())
	pkglog.LogFormat(rootCmd.PersistentFlags())

	// Register OpenTracing flags
	jaeger.InitFlags(rootCmd.PersistentFlags())

	// Add Run command
	rootCmd.AddCommand(txcrafter.NewRootCommand())
	rootCmd.AddCommand(txsigner.NewRootCommand())
	rootCmd.AddCommand(txsender.NewRootCommand())
	rootCmd.AddCommand(txlistener.NewRootCommand())
	rootCmd.AddCommand(contractregistry.NewRootCommand())
	rootCmd.AddCommand(envelopestore.NewRootCommand())
	rootCmd.AddCommand(transactionscheduler.NewRootCommand())
	rootCmd.AddCommand(chainregistry.NewRootCommand())
	rootCmd.AddCommand(envelopestore.NewRootCommand())
	rootCmd.AddCommand(utils.NewRootCommand())

	// Register Multi-Tenancy flags
	multitenancy.Enabled(rootCmd.PersistentFlags())
	authjwt.Flags(rootCmd.PersistentFlags())
	authkey.APIKey(rootCmd.PersistentFlags())

	return rootCmd
}
