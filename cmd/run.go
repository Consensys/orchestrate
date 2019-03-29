package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/app"
	ethabi "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/abi"

)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register flags
	ethabi.ABIs(runCmd.Flags())
	config.HTTPHostname(runCmd.Flags())
	config.EthClientURLs(runCmd.Flags())
	config.KafkaAddresses(runCmd.Flags())
	config.TxDecoderInTopic(runCmd.Flags())
	config.TxDecodedOutTopic(runCmd.Flags())
	config.WorkerDecoderGroup(runCmd.Flags())
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
