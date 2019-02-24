package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/common.git/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/common.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/app/infra"
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
	config.WorkerSignerGroup(runCmd.Flags())
	config.WorkerSlots(runCmd.Flags())
	infra.VaultAccounts(runCmd.Flags())

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
