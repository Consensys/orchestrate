package storer

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

func UnsignedTxStore(store svc.EnvelopeStoreClient, txSchedulerClient client.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// TODO: Remove if statement when envelope store is removed
		if txctx.Envelope.BelongToEnvelopeStore() {
			unsignedTxStoreInEnvelopeStore(txctx, store)
		} else {
			unsignedTxStoreInTxScheduler(txctx, txSchedulerClient)
		}
	}
}

func unsignedTxStoreInTxScheduler(txctx *engine.TxContext, txSchedulerClient client.TransactionSchedulerClient) {
	// Execute pending handlers (expected to send the transaction)
	txctx.Next()

	// We do not retry on unsigned txs
	if len(txctx.Envelope.GetErrors()) != 0 {
		return
	}

	// Transaction hash is generated after the transaction is sent
	txctx.Logger.Debug("transaction scheduler: updating transaction to SENT")
	_, err := txSchedulerClient.UpdateJob(
		txctx.Context(),
		txctx.Envelope.GetID(),
		&types.UpdateJobRequest{
			Transaction: &types.ETHTransaction{
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
			Status: types.StatusSent,
		},
	)

	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("transaction scheduler: failed to update transaction")
		return
	}
}

func unsignedTxStoreInEnvelopeStore(txctx *engine.TxContext, store svc.EnvelopeStoreClient) {
	// Execute pending handlers (expected to send the transaction)
	txctx.Next()

	// If no error occurred while executing pending handlers
	if len(txctx.Envelope.GetErrors()) == 0 {
		// Store envelope
		// We can not store envelope before sending transaction because we do not know the transaction hash
		// This is an issue for overall consistency of the system before/after transaction is mined
		_, err := store.Store(txctx.Context(),
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

		// Transaction has been properly sent so we set status to `pending`
		_, err = store.SetStatus(txctx.Context(),
			&svc.SetStatusRequest{
				Id:     txctx.Envelope.GetID(),
				Status: svc.Status_PENDING,
			},
		)
		if err != nil {
			// Connection to store is broken
			e := errors.FromError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Warnf("store: failed to set status")
			return
		}
	}
}
