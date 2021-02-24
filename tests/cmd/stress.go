package main

import (
	"context"
	"os"

	broker "github.com/ConsenSys/orchestrate/pkg/broker/sarama"

	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/ConsenSys/orchestrate/tests/service/stress"
	"github.com/spf13/cobra"
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
