package contractregistry

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// ContractRegistry flag
	contractregistry.Type(runCmd.Flags())
	contractregistry.ABIs(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { Close(context.Background()) })
	defer sig.Close()

	// Start application
	Start(context.Background())
}
