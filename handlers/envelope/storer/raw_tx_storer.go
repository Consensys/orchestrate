package storer

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	types2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
)

func RawTxStore(store svc.EnvelopeStoreClient, txSchedulerClient client.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// TODO: Remove if statement when envelope store is removed
		if txctx.Envelope.ContextLabels["jobUUID"] != "" {
			updateTxInScheduler(txctx, txSchedulerClient)
		} else {
			updateTxInStore(txctx, store)
		}
	}
}

func updateTxInScheduler(txctx *engine.TxContext, txSchedulerClient client.TransactionSchedulerClient) {
	txctx.Logger.Debug("updating transaction")

	_, err := txSchedulerClient.UpdateJob(
		txctx.Context(),
		txctx.Envelope.GetID(),
		&types.UpdateJobRequest{
			Transaction: &types2.ETHTransaction{
				Hash:           txctx.Envelope.GetTxHashString(),
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
			Status: types2.StatusPending,
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
		_, storeErr := txSchedulerClient.UpdateJob(
			txctx.Context(),
			txctx.Envelope.GetID(),
			&types.UpdateJobRequest{
				Status: types2.StatusRecovering,
				Message: fmt.Sprintf(
					"transaction attempt with nonce %v and sender %v failed with error: %v",
					txctx.Envelope.GetNonceString(),
					txctx.Envelope.GetFromString(),
					txctx.Envelope.Error(),
				),
			},
		)
		if storeErr != nil {
			e := errors.FromError(storeErr).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("transaction scheduler: failed to set transaction status for recovering")
		}
		return
	}

	// Transaction has been properly sent so we set status to `sent`
	_, err = txSchedulerClient.UpdateJob(
		txctx.Context(),
		txctx.Envelope.GetID(),
		&types.UpdateJobRequest{
			Status: types2.StatusSent,
		},
	)
	if err != nil {
		e := errors.FromError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("transaction scheduler: failed to set transaction status")
		return
	}
}

func updateTxInStore(txctx *engine.TxContext, store svc.EnvelopeStoreClient) {
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
		txctx.Logger.WithError(e).Errorf("sender: failed to set envelope status")
		return
	}
}
