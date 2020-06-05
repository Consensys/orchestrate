package txupdater

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

const component = "handler.tx-updater"

// TransactionUpdater updates a transaction in the scheduler
func TransactionUpdater(txSchedulerClient txscheduler.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Next()

		// TODO: Remove this if when envelope store is removed
		if txctx.Envelope.ContextLabels["jobUUID"] != "" {
			if !txctx.Envelope.OnlyWarnings() {
				_, err := txSchedulerClient.UpdateJob(txctx.Context(), txctx.Envelope.GetID(), &types.UpdateJobRequest{
					Status: entities.JobStatusFailed,
				})

				if err != nil {
					e := txctx.AbortWithError(err).ExtendComponent(component)
					txctx.Logger.WithError(e).Errorf("tx updater: could not update transaction status")
					return
				}
			}
		}
	}
}
