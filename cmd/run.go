package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/contract-registry.git/app"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Opentracing flags
	jaeger.InitFlags(runCmd.Flags())

	// Register HTTP server flags
	http.Hostname(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, args []string) {
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { app.Close(context.Background()) })
	defer sig.Close()

	// Start application
	app.Start(context.Background())
}
