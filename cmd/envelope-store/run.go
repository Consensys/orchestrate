package envelopestore

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/tracing/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/envelope-store"
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

	// EnvelopeStore flag
	envelopestore.Type(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, args []string) {
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { Close(context.Background()) })
	defer sig.Close()

	// Start application
	Start(context.Background())
}
