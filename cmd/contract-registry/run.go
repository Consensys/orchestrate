package contractregistry

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Set flags
	contractregistry.Flags(runCmd)

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { contractregistry.StopService(context.Background()) })
	defer sig.Close()

	// Start application
	contractregistry.StartService(context.Background())
}
