package storer

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/token"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
)

func RawTxStore(store evlpstore.EnvelopeStoreClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Extract JWT if present
		jwTokenGRPCOption := token.GetGRPCOptionJWTokenFromEnvelope(txctx)

		// Store envelope
		_, err := store.Store(
			txctx.Context(),
			&evlpstore.StoreRequest{
				Envelope: txctx.Envelope,
			},
			jwTokenGRPCOption)
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
				&evlpstore.SetStatusRequest{
					Id:     txctx.Envelope.GetMetadata().GetId(),
					Status: evlpstore.Status_ERROR,
				},
				jwTokenGRPCOption)
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
			&evlpstore.SetStatusRequest{
				Id:     txctx.Envelope.GetMetadata().GetId(),
				Status: evlpstore.Status_PENDING,
			},
			jwTokenGRPCOption)
		if err != nil {
			// Connection to store is broken
			e := errors.FromError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to set envelope status")
			return
		}
	}
}
