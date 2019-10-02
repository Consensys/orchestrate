package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/tests/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/tests/service/cucumber"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/tests/service/cucumber/steps"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Engine flags
	engine.InitFlags(runCmd.Flags())

	// Register HTTP server flags
	http.Hostname(runCmd.Flags())

	// Register Kafka flags
	broker.KafkaAddresses(runCmd.Flags())
	broker.KafkaGroup(runCmd.Flags())
	broker.KafkaTopicTxCrafter(runCmd.Flags())
	broker.KafkaTopicTxNonce(runCmd.Flags())
	broker.KafkaTopicTxSigner(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicTxDecoder(runCmd.Flags())
	broker.KafkaTopicTxDecoded(runCmd.Flags())
	broker.KafkaTopicWalletGenerator(runCmd.Flags())
	broker.KafkaTopicWalletGenerated(runCmd.Flags())
	broker.InitKafkaSASLTLSFlags(runCmd.Flags())

	// Register Cucumber flag
	cucumber.InitFlags(runCmd.Flags())
	steps.InitFlags(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start application
	app.Start(ctx)
}
