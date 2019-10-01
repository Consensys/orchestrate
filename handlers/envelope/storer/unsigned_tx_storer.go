package storer

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/envelope-store"
)

func UnsignedTxStore(store evlpstore.EnvelopeStoreClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Execute pending handlers (expected to send the transaction)
		txctx.Next()

		// If no error occurred while executing pending handlers
		if len(txctx.Envelope.GetErrors()) == 0 {
			// Store envelope
			// We can not store envelope before sending transaction because we do not know the transaction hash
			// This is an issue for overall consistency of the system before/after transaction is mined
			_, err := store.Store(txctx.Context(), &evlpstore.StoreRequest{
				Envelope: txctx.Envelope,
			})
			if err != nil {
				// Connection to store is broken
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("store: failed to store envelope")
				return
			}

			// Transaction has been properly sent so we set status to `pending`
			_, err = store.SetStatus(txctx.Context(), &evlpstore.SetStatusRequest{
				Id:     txctx.Envelope.GetMetadata().GetId(),
				Status: evlpstore.Status_PENDING,
			})
			if err != nil {
				// Connection to store is broken
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Warnf("store: failed to set status")
				return
			}
		}
	}
}
