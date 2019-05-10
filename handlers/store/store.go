package store

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// EnvelopeLoader creates and handler that load traces
func EnvelopeLoader(s store.EnvelopeStore) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		_, _, err := s.LoadByTxHash(txctx.Context(), txctx.Envelope.GetChain().GetId(), txctx.Envelope.GetReceipt().GetTxHash(), txctx.Envelope)

		if err != nil {
			// We got an error, possibly due to timeout Connection to database or something else
			// TODO: what should we do in case of error?
			_ = txctx.Error(err)
			txctx.Logger.WithError(err).Debugf("envelope-loader: no trace stored")
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		// Transaction has been mined so we set status to `mined`
		err = s.SetStatus(txctx.Context(), txctx.Envelope.GetMetadata().GetId(), "mined")
		if err != nil {
			// Connection to store is broken
			txctx.Logger.WithError(err).Errorf("envelope-loader: trace store failed to set status")
			_ = txctx.Error(err)
		}

		txctx.Logger.Debugf("envelope-loader: envelope re-constituted")
	}
}
