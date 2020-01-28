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
		if txctx.Envelope.GetChain().GetBigChainID() == nil {
			err := errors.DataError("cannot sign transaction without chainID").SetComponent(component)
			txctx.Logger.WithError(err).Errorf("failed to sign transaction")
			_ = txctx.AbortWithError(err)
			return
		}

		switch {
		default:
			// Default sign for public ethereum
			publicEthereumSigner(txctx)
		case txctx.Envelope.Protocol != nil && txctx.Envelope.Protocol.IsConstellation():
			// Do nothing as the ethereum node is going to perform the signature
		case txctx.Envelope.Protocol != nil && txctx.Envelope.Protocol.IsBesu():
			// Sign for EEA private transaction implementation
			eeaSigner(txctx)
		case txctx.Envelope.Protocol != nil && txctx.Envelope.Protocol.IsTessera():
			// Sign for Tessera
			tesseraSigner(txctx)
		}
	}
}
