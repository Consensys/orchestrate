package main

import (
	"context"
	"os"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"

	traefiklog "github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
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
	broker.ConsumerGroupName(runCmd.Flags())

	return runCmd
}

func runStress(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	if err := stress.Start(ctx); err != nil {
		traefiklog.WithoutContext().WithError(err).Errorf("test execution did not complete successfully")
		return err
	}

	if err := stress.Stop(ctx); err != nil {
		traefiklog.WithoutContext().WithError(err).Errorf("test execution did not shutdown properly")
	} else {
		traefiklog.WithoutContext().Info("test execution gracefully closed")
	}

	return nil
}
