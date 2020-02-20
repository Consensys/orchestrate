package main

import (
	"context"
	"os"

	authentication "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt/generator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/steps"
)

func NewRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
		PreRun: func(cmd *cobra.Command, args []string) {
			// This is executed before each run (included on children command run)
			logger.InitLogger()
		},
	}

	// Set logger flags
	logger.LogLevel(runCmd.Flags())
	logger.LogFormat(runCmd.Flags())

	// Register HTTP server flags
	metrics.Hostname(runCmd.Flags())

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxCrafter(runCmd.Flags())
	broker.KafkaTopicTxNonce(runCmd.Flags())
	broker.KafkaTopicTxSigner(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicTxDecoded(runCmd.Flags())
	broker.KafkaTopicWalletGenerator(runCmd.Flags())
	broker.KafkaTopicWalletGenerated(runCmd.Flags())

	// Register Cucumber flag
	cucumber.InitFlags(runCmd.Flags())
	steps.InitFlags(runCmd.Flags())

	// Register Multi-Tenancy flags
	multitenancy.Enabled(runCmd.Flags())
	authentication.Flags(runCmd.Flags())
	generator.PrivateKey(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start application
	Start(ctx)
}
