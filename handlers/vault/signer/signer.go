package signer

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// TxSigner creates a signer handler
//
// It is a fork handler that allow signing for either eea, tessera or public ethereum
func TxSigner(eeaSigner, publicEthereumSigner, tesseraSigner engine.HandlerFunc) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Builder.GetChainID() == nil {
			err := errors.DataError("cannot sign transaction without chainID").SetComponent(component)
			txctx.Logger.WithError(err).Errorf("failed to sign transaction")
			_ = txctx.AbortWithError(err)
			return
		}

		switch {
		case txctx.Builder.IsEthSendPrivateTransaction():
			// Do nothing as the ethereum node is going to perform the signature
		case txctx.Builder.IsEthSendRawPrivateTransaction():
			// Sign for Tessera
			tesseraSigner(txctx)
		case txctx.Builder.IsEeaSendPrivateTransaction():
			// Sign for EEA private transaction implementation
			eeaSigner(txctx)
		default:
			// Default sign for public ethereum
			publicEthereumSigner(txctx)
		}
	}
}
