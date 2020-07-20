package storer

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

func RawTxStore(store svc.EnvelopeStoreClient, txSchedulerClient client.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// TODO: Remove if statement when envelope store is removed
		if txctx.Envelope.BelongToEnvelopeStore() {
			rawTxStoreInEnvelopeStore(txctx, store)
		} else {
			rawTxStoreInTxScheduler(txctx, txSchedulerClient)
		}
	}
}

func rawTxStoreInTxScheduler(txctx *engine.TxContext, txSchedulerClient client.TransactionSchedulerClient) {
	txctx.Logger.Debug("transaction scheduler: updating transaction to PENDING")

	computedTxHash := txctx.Envelope.GetTxHashString()
	_, err := txSchedulerClient.UpdateJob(
		txctx.Context(),
		txctx.Envelope.GetID(),
		&types.UpdateJobRequest{
			Transaction: &types.ETHTransaction{
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
		},
	)

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
			&types.UpdateJobRequest{
				Status: utils.StatusRecovering,
				Message: fmt.Sprintf(
					"transaction attempt with nonce %v and sender %v failed with error: %v",
					txctx.Envelope.GetNonceString(),
					txctx.Envelope.GetFromString(),
					txctx.Envelope.Error(),
				),
			},
		)
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
			&types.UpdateJobRequest{
				Transaction: &types.ETHTransaction{
					Hash: retrievedTxHash,
				},
				Status:  utils.StatusWarning,
				Message: errMessage,
			},
		)
		if updateErr != nil {
			e := errors.FromError(updateErr).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("transaction scheduler: failed to set transaction status for recovering")
		}
	}

	txctx.Logger.Info("transaction successfully sent to the Blockchain node")
}

func rawTxStoreInEnvelopeStore(txctx *engine.TxContext, store svc.EnvelopeStoreClient) {
	// Store envelope
	_, err := store.Store(
		txctx.Context(),
		&svc.StoreRequest{
			Envelope: txctx.Envelope.TxEnvelopeAsRequest(),
		},
	)
	if err != nil {
		// Connection to store is broken
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("store: failed to store envelope")
		return
	}

	// Execute pending handlers (expected to send the transaction)
	txctx.Next()

	// If an error occurred when executing pending handlers
	if len(txctx.Envelope.GetErrors()) != 0 {
		// We update status in storage
		txctx.Logger.Debug("store: updating transaction to ERROR")
		_, storeErr := store.SetStatus(
			txctx.Context(),
			&svc.SetStatusRequest{
				Id:     txctx.Envelope.GetID(),
				Status: svc.Status_ERROR,
			},
		)
		if storeErr != nil {
			// Connection to store is broken
			e := errors.FromError(storeErr).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("store: failed to set envelope status")
		}
		return
	}

	// Transaction has been properly sent so we set status to `pending`
	txctx.Logger.Debug("store: updating transaction to PENDING")
	_, err = store.SetStatus(
		txctx.Context(),
		&svc.SetStatusRequest{
			Id:     txctx.Envelope.GetID(),
			Status: svc.Status_PENDING,
		},
	)
	if err != nil {
		// Connection to store is broken
		e := errors.FromError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("store: failed to set envelope status")
		return
	}
}
