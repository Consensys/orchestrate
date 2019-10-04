package envelopestore

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Hostname & port for servers
	grpcserver.Hostname(runCmd.Flags())
	grpcserver.Port(runCmd.Flags())
	rest.Hostname(runCmd.Flags())
	rest.Port(runCmd.Flags())

	// EnvelopeStore flag
	envelopestore.Type(runCmd.Flags())

	// Postgres flags
	postgres.PGFlags(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { Close(context.Background()) })
	defer sig.Close()

	// Start application
	Start(context.Background())
}
