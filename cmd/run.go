package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/tracing/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/contract-registry.git/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/rpc"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register OpenTracing flags
	jaeger.InitFlags(runCmd.Flags())

	// Register HTTP server flags
	http.Hostname(runCmd.Flags())

	// EthClient flag
	ethclient.URLs(runCmd.Flags())

	// ContractRegistry flag
	registry.ContractRegistryType(runCmd.Flags())
	registry.ABIs(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, args []string) {
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { app.Close(context.Background()) })
	defer sig.Close()

	// Start application
	app.Start(context.Background())
}
