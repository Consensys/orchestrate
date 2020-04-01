package contractregistry

import (
	"context"
	"os"

	contract_registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Set flags
	contract_registry.Flags(runCmd)

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { contract_registry.StopService(context.Background()) })
	defer sig.Close()

	// Start application
	contract_registry.StartService(context.Background(), viper.GetString(contract_registry.TypeViperKey))
}
