package main

import (
	"context"
	"os"

	broker "github.com/consensys/orchestrate/pkg/broker/sarama"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/tests/service/stress"
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
