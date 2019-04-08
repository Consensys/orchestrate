package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Loader creates an handler loading input
func Loader(u services.Unmarshaller) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Unmarshal message
		err := u.Unmarshal(txctx.Msg, txctx.Envelope)

		if err != nil {
			// TODO: handle error
			txctx.Logger.Errorf("Error unmarshalling: %v", err)
			txctx.AbortWithError(err)
			return
		}

		txctx.Logger.Debugf("Message unmarshalled: %v", txctx.Envelope.String())
	}
}
