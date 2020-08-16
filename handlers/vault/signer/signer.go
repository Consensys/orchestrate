package signer

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// TxSigner creates a signer handler
//
// It is a fork handler that allow signing for either eea, tessera or public ethereum
func TxSigner(publicEthereumSigner, eeaSigner, tesseraSigner engine.HandlerFunc) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetChainID() == nil {
			err := errors.DataError("cannot sign transaction without chainID").SetComponent(component)
			txctx.Logger.WithError(err).Errorf("failed to sign transaction")
			_ = txctx.AbortWithError(err)
			return
		}

		switch {
		case txctx.Envelope.IsEthSendTesseraPrivateTransaction():
			// StoreRaw transaction does not require to be signed
			return
		case txctx.Envelope.IsEthSendTesseraMarkingTransaction():
			// Sign for Tessera
			tesseraSigner(txctx)
		case txctx.Envelope.IsEeaSendPrivateTransaction():
			// Sign for EEA private transaction implementation
			eeaSigner(txctx)
		default:
			// Default sign for public ethereum
			publicEthereumSigner(txctx)
		}
	}
}
