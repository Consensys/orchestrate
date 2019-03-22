package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/app"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register flags
	config.HTTPHostname(runCmd.Flags())
	config.KafkaAddresses(runCmd.Flags())
	config.TxSignerInTopic(runCmd.Flags())
	config.TxSenderOutTopic(runCmd.Flags())
	secretstore.InitFlags(runCmd.Flags())
	config.WorkerSignerGroup(runCmd.Flags())
	keystore.SecretPkeys(runCmd.Flags())
	worker.InitFlags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, args []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())
	a := app.New(ctx)

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Run App
	a.Run()

	// Wait for app to properly close
	<-a.Done()
}
