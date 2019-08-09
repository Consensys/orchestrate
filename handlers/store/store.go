package store

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
)

// EnvelopeLoader creates and handler that load envelopes
func EnvelopeLoader(s evlpstore.EnvelopeStoreClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		resp, err := s.LoadByTxHash(
			txctx.Context(),
			&evlpstore.LoadByTxHashRequest{
				Chain:  txctx.Envelope.GetChain(),
				TxHash: txctx.Envelope.GetReceipt().GetTxHash(),
			},
		)
		if err != nil {
			// We got an error, possibly due to timeout Connection to database or something else
			// TODO: what should we do in case of error?
			_ = txctx.Error(err)
			txctx.Logger.WithError(err).Debugf("envelope-loader: no envelopes stored")
			return
		}

		// Set context envelope
		txctx.Envelope = resp.GetEnvelope()

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		// Transaction has been mined so we set status to `mined`
		_, err = s.SetStatus(
			txctx.Context(),
			&evlpstore.SetStatusRequest{
				Id:     txctx.Envelope.GetMetadata().GetId(),
				Status: evlpstore.Status_MINED,
			},
		)
		if err != nil {
			// Connection to store is broken
			txctx.Logger.WithError(err).Errorf("envelope-loader: envelope store failed to set status")
			_ = txctx.Error(err)
		}

		txctx.Logger.Debugf("envelope-loader: envelope re-constituted")
	}
}
