package main

import (
	"context"
	"os"

	traefiklog "github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/stress"
)

func NewRunStressTestCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "stress",
		Short: "Run stress test",
		Run:   runStress,
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := errors.CombineErrors(cmdErr, cmd.Context().Err()); err != nil {
				os.Exit(1)
			}
		},
	}

	// Register Stress flag
	stress.InitFlags(runCmd.Flags())

	return runCmd
}

func runStress(cmd *cobra.Command, _ []string) {
	ctx, cancel := context.WithCancel(cmd.Context())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	if err := stress.Start(ctx); err != nil {
		cmdErr = errors.CombineErrors(cmdErr, err)
		traefiklog.WithoutContext().WithError(err).Errorf("test execution did not complete successfully")
	}

	if err := stress.Stop(ctx); err != nil {
		traefiklog.WithoutContext().WithError(err).Errorf("test execution did not shutdown properly")
	} else {
		traefiklog.WithoutContext().Info("test execution gracefully closed")
	}
}
