package chainregistry

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Hostname & port for servers
	rest.Flags(runCmd.Flags())

	// Register flags for chain registry
	chainregistry.Flags(runCmd.Flags())

	// Postgres flags
	postgres.PGFlags(runCmd.Flags())

	// Chain registry store flags
	store.Flags(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start application
	Start(ctx)
}
