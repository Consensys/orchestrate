package main

import (
	"context"
	"os"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress"
)

func NewRunStressTestCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:          "stress",
		Short:        "Run stress test",
		RunE:         runStress,
		SilenceUsage: true,
	}

	// Register Stress flag
	stress.InitFlags(runCmd.Flags())
	broker.KafkaConsumerFlags(runCmd.Flags())

	return runCmd
}

func runStress(cmd *cobra.Command, _ []string) error {
	logger := log.NewLogger().SetComponent("cmd-stress-test")
	ctx, cancel := context.WithCancel(cmd.Context())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	if err := stress.Start(ctx); err != nil {
		logger.WithError(err).Error("failed to complete")
		return err
	}

	if err := stress.Stop(ctx); err != nil {
		logger.WithError(err).Errorf("execution did not shutdown properly")
	} else {
		logger.Info("execution gracefully closed")
	}

	return nil
}
