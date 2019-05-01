package cmd

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
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
	// Create app
	a := app.New()

	// Process signals
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig := utils.NewSignalListener(func(signal os.Signal) { a.Close(ctx) })
	defer sig.Close()

	// Initialize  App
	a.Run()

	// Wait
	<-a.Done()
}
