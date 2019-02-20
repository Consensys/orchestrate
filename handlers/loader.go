package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Loader creates an handler loading input
func Loader(u services.Unmarshaller) types.HandlerFunc {
	return func(ctx *types.Context) {
		// Unmarshal message
		err := u.Unmarshal(ctx.Msg, ctx.T)

		if err != nil {
			// TODO: handle error
			ctx.Logger.Errorf("error unmarshalling: %v", err)
			ctx.AbortWithError(err)
			return
		}

		ctx.Logger.Debugf("message unmarshalled: %v", ctx.T)
	}
}
