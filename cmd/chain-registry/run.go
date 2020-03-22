package chainregistry

import (
	"context"
	"os"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
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

	// Chain registry store flags
	store.Flags(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	_ = chainregistry.Start(context.Background())

	done := make(chan struct{})

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		err := chainregistry.Stop(context.Background())
		if err != nil {
			log.WithoutContext().WithError(err).Errorf("Application did not shutdown properly")
		} else {
			log.WithoutContext().WithError(err).Infof("Application gracefully closed")
		}
		close(done)
	})
	<-done

	sig.Close()
}
