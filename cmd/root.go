package cmd

import (
	"github.com/spf13/cobra"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/chain-registry"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/contract-registry"
	identitymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/identity-manager"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/key-manager"
	transactionscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/transaction-scheduler"
	txcrafter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/tx-crafter"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/tx-listener"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/tx-sender"
	txsigner "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/tx-signer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/utils"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tracing/opentracing/jaeger"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "orchestrate",
		TraverseChildren: true,
		SilenceUsage:     true,
	}

	// Set Persistent flags
	log.Level(rootCmd.PersistentFlags())
	log.Format(rootCmd.PersistentFlags())

	// Register OpenTracing flags
	jaeger.InitFlags(rootCmd.PersistentFlags())

	// Add Run command
	rootCmd.AddCommand(txcrafter.NewRootCommand())
	rootCmd.AddCommand(txsigner.NewRootCommand())
	rootCmd.AddCommand(txsender.NewRootCommand())
	rootCmd.AddCommand(txlistener.NewRootCommand())
	rootCmd.AddCommand(contractregistry.NewRootCommand())
	rootCmd.AddCommand(transactionscheduler.NewRootCommand())
	rootCmd.AddCommand(chainregistry.NewRootCommand())
	rootCmd.AddCommand(identitymanager.NewRootCommand())
	rootCmd.AddCommand(keymanager.NewRootCommand())
	rootCmd.AddCommand(utils.NewRootCommand())

	// Register Multi-Tenancy flags
	multitenancy.Enabled(rootCmd.PersistentFlags())
	authjwt.Flags(rootCmd.PersistentFlags())
	authkey.APIKey(rootCmd.PersistentFlags())

	http.MetricFlags(rootCmd.PersistentFlags())

	return rootCmd
}
