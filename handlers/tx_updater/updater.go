package txupdater

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

const component = "handler.tx-updater"

// TransactionUpdater updates a transaction
func TransactionUpdater(client orchestrateclient.OrchestrateClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.OnlyWarnings() {
			return
		}

		// In case we are retrying on same MSG we exit
		if txctx.HasRetryMsgErr() != nil {
			return
		}

		// Don't update to FAILED if we are going to send message to tx-crafter
		if txctx.HasInvalidNonceErr() {
			txctx.Logger.Debug("transaction scheduler: updating transaction to RECOVERING")
			_, err := client.UpdateJob(
				txctx.Context(),
				txctx.Envelope.GetJobUUID(),
				&types.UpdateJobRequest{
					Status:  utils.StatusRecovering,
					Message: txctx.Envelope.Error(),
				})

			if err != nil {
				e := txctx.Error(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("tx updater: could not update transaction status")
			}
			return
		}

		// TODO: Improvement of the log message will be done when we move to clean architecture
		// TODO: because at the moment it is difficult to know what error messages need to be sent to users and which ones not.
		_, err := client.UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), &types.UpdateJobRequest{
			Status:  utils.StatusFailed,
			Message: txctx.Envelope.Error(),
		})

		if err != nil {
			e := txctx.Error(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("tx updater: could not update transaction status")
			return
		}
	}
}
