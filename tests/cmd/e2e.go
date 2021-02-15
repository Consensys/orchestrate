package main

import (
	"context"
	"os"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/cucumber"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/cucumber/steps"
)

func NewRunE2ECommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:          "e2e",
		Short:        "Run e2e test",
		RunE:         runE2E,
		SilenceUsage: true,
	}

	// Register Cucumber flag
	cucumber.InitFlags(runCmd.Flags())
	steps.InitFlags(runCmd.Flags())
	e2e.InitFlags(runCmd.Flags())
	broker.ConsumerGroupName(runCmd.Flags())

	return runCmd
}

func runE2E(cmd *cobra.Command, _ []string) error {
	logger := log.NewLogger().SetComponent("e2e.cucumber")
	ctx, cancel := context.WithCancel(cmd.Context())
	ctx = log.With(ctx, logger)

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	var gerr error
	if err := e2e.Start(ctx); err != nil {
		logger.WithError(err).Error("did not complete successfully")
		gerr = errors.CombineErrors(gerr, err)
	}

	if err := e2e.Stop(ctx); err != nil {
		logger.WithError(err).Error("did not shutdown properly")
		gerr = errors.CombineErrors(gerr, err)
	}

	return gerr
}
