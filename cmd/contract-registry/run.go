package contractregistry

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Hostname & port for servers
	grpcserver.Flags(runCmd.Flags())
	rest.Flags(runCmd.Flags())

	// ContractRegistry flag
	contractregistry.Type(runCmd.Flags())
	contractregistry.ABIs(runCmd.Flags())

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
