package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Loader creates an handler loading input
func Loader(u services.Unmarshaller) engine.HandlerFunc {
	return func(ctx *engine.TxContext) {
		// Unmarshal message
		err := u.Unmarshal(ctx.Msg, ctx.Envelope)

		if err != nil {
			// TODO: handle error
			ctx.Logger.Errorf("Error unmarshalling: %v", err)
			ctx.AbortWithError(err)
			return
		}

		ctx.Logger.Debugf("Message unmarshalled: %v", ctx.Envelope.String())
	}
}
