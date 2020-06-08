package txupdater

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	types2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
)

const component = "handler.tx-updater"

// TransactionUpdater updates a transaction in the scheduler
func TransactionUpdater(txSchedulerClient txscheduler.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Next()

		// TODO: Remove next statement once envelope store is removed
		if txctx.Envelope.ContextLabels["jobUUID"] == "" {
			return
		}

		if txctx.Envelope.OnlyWarnings() {
			return
		}

		_, err := txSchedulerClient.UpdateJob(txctx.Context(), txctx.Envelope.GetID(), &types.UpdateJobRequest{
			Status: types2.StatusFailed,
		})

		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("tx updater: could not update transaction status")
			return
		}
	}
}
