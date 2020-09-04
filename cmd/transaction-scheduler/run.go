package transactionscheduler

import (
	"os"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	transactionscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler"
)

var cmdErr error

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.PreRunBindFlags(viper.GetViper(), cmd.Flags(), "transaction-scheduler")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := errors.CombineErrors(cmdErr, cmd.Context().Err()); err != nil {
				os.Exit(1)
			}
		},
	}

	// Transaction scheduler flags
	transactionscheduler.Flags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()
	logger := log.FromContext(ctx)

	// Initialize and start service
	txScheduler, err := transactionscheduler.New(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to initialize transaction scheduler")
		cmdErr = errors.CombineErrors(cmdErr, err)
		return
	}

	err = txScheduler.Start(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to start transaction scheduler")
		cmdErr = errors.CombineErrors(cmdErr, err)
		return
	}

	// Process signals
	done := make(chan struct{})
	sig := utils.NewSignalListener(func(signal os.Signal) { close(done) })

	<-done

	sig.Close()
	err = txScheduler.Stop(ctx)
	if err != nil {
		cmdErr = errors.CombineErrors(cmdErr, err)
	}

	logger.Info("transaction scheduler stopped successfully")
}
