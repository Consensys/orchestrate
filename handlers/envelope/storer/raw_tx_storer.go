package storer

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx-scheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

func RawTxStore(txSchedulerClient client.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.Debug("transaction scheduler: updating transaction to PENDING")

		computedTxHash := txctx.Envelope.GetTxHashString()
		_, err := txSchedulerClient.UpdateJob(
			txctx.Context(),
			txctx.Envelope.GetID(),
			&txschedulertypes.UpdateJobRequest{
				Transaction: &entities.ETHTransaction{
					Hash:           computedTxHash,
					From:           txctx.Envelope.GetFromString(),
					To:             txctx.Envelope.GetToString(),
					Nonce:          txctx.Envelope.GetNonceString(),
					Value:          txctx.Envelope.GetValueString(),
					GasPrice:       txctx.Envelope.GetGasPriceString(),
					Gas:            txctx.Envelope.GetGasString(),
					Raw:            txctx.Envelope.GetRaw(),
					PrivateFrom:    txctx.Envelope.GetPrivateFrom(),
					PrivateFor:     txctx.Envelope.GetPrivateFor(),
					PrivacyGroupID: txctx.Envelope.GetPrivacyGroupID(),
				},
				Status: utils.StatusPending,
			})

		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("transaction scheduler: failed to update transaction")
			return
		}

		// Execute pending handlers (expected to send the transaction)
		txctx.Next()

		// If an error occurred when executing pending handlers
		if len(txctx.Envelope.GetErrors()) != 0 {
			txctx.Logger.Debug("transaction scheduler: updating transaction to RECOVERING")
			_, updateErr := txSchedulerClient.UpdateJob(
				txctx.Context(),
				txctx.Envelope.GetID(),
				&txschedulertypes.UpdateJobRequest{
					Status: utils.StatusRecovering,
					Message: fmt.Sprintf(
						"transaction attempt with nonce %v and sender %v failed with error: %v",
						txctx.Envelope.GetNonceString(),
						txctx.Envelope.GetFromString(),
						txctx.Envelope.Error(),
					),
				})

			if updateErr != nil {
				e := errors.FromError(updateErr).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("transaction scheduler: failed to set transaction status for recovering")
			}
			return
		}

		retrievedTxHash := txctx.Envelope.GetTxHashString()
		if computedTxHash != retrievedTxHash {
			errMessage := fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", computedTxHash, retrievedTxHash)
			txctx.Logger.Errorf("errMessage")

			_, updateErr := txSchedulerClient.UpdateJob(
				txctx.Context(),
				txctx.Envelope.GetID(),
				&txschedulertypes.UpdateJobRequest{
					Transaction: &entities.ETHTransaction{
						Hash: retrievedTxHash,
					},
					Status:  utils.StatusWarning,
					Message: errMessage,
				})

			if updateErr != nil {
				e := errors.FromError(updateErr).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("transaction scheduler: failed to set transaction status for recovering")
			}
		}

		txctx.Logger.Info("transaction successfully sent to the Blockchain node")
	}
}
