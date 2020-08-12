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
	transactionscheduler.TxSchedulerFlags(runCmd.Flags())
	transactionscheduler.TxSentryFlags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()
	logger := log.FromContext(ctx)

	txScheduler, err := transactionscheduler.NewTxScheduler(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to initialize tx-scheduler")
	}

	done := make(chan struct{})
	sentryErrorChan := txScheduler.StartSentry(ctx)
	err = txScheduler.StartScheduler(ctx)
	if err != nil {
		cmdErr = errors.CombineErrors(cmdErr, err)
		logger.WithError(err).Error("could not start transaction-scheduler API")
		close(done)
	}

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { close(done) })

	select {
	case sentryErr := <-sentryErrorChan:
		cmdErr = errors.CombineErrors(cmdErr, sentryErr)
	case <-done:
	}

	sig.Close()

	err = txScheduler.StopScheduler(ctx)
	if err != nil {
		cmdErr = errors.CombineErrors(cmdErr, err)
		logger.WithError(err).Error("could not stop transaction-scheduler API")
	}
	txScheduler.StopSentry(ctx)

	logger.Info("transaction scheduler and all its services successfully stopped")
}
