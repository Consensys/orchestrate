package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Loader creates an handler loading input
func Loader(u services.Unmarshaller) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		// Unmarshal message
		err := u.Unmarshal(ctx.Msg, ctx.T)

		if err != nil {
			// TODO: handle error
			ctx.Logger.Errorf("Error unmarshalling: %v", err)
			ctx.AbortWithError(err)
			return
		}

		ctx.Logger.Debugf("Message unmarshalled: %v", ctx.T.String())
	}
}
