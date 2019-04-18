package loader

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/loader/domain"
)

// Loader creates an handler loading input
func Loader(u domain.Unmarshaller) engine.HandlerFunc {
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
