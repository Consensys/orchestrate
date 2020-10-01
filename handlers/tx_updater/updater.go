package txupdater

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

const component = "handler.tx-updater"

// TransactionUpdater updates a transaction in the scheduler
func TransactionUpdater(txSchedulerClient txscheduler.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Next()

		if txctx.Envelope.OnlyWarnings() {
			return
		}

		// TODO: Improvement of the log message will be done when we move to clean architecture
		// TODO: because at the moment it is difficult to know what error messages need to be sent to users and which ones not.
		_, err := txSchedulerClient.UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), &txschedulertypes.UpdateJobRequest{
			Status:  utils.StatusFailed,
			Message: txctx.Envelope.Error(),
		})

		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("tx updater: could not update transaction status")
			return
		}
	}
}
