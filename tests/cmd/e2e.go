package main

import (
	"context"
	"os"

	traefiklog "github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/e2e"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/e2e/cucumber"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/e2e/cucumber/alias"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/e2e/cucumber/steps"
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
	alias.InitFlags(runCmd.Flags())

	return runCmd
}

func runE2E(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	if err := e2e.Start(ctx); err != nil {
		traefiklog.WithoutContext().WithError(err).Errorf("Cucumber did not complete successfully")
		return err
	}

	if err := e2e.Stop(ctx); err != nil {
		traefiklog.WithoutContext().WithError(err).Errorf("Cucumber did not shutdown properly")
	} else {
		traefiklog.WithoutContext().Info("Cucumber gracefully closed")
	}

	return nil
}
