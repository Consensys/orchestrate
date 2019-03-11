package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"
)

// Loader creates an handler loading input
func Loader(u services.Unmarshaller) types.HandlerFunc {
	return func(ctx *types.Context) {
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
