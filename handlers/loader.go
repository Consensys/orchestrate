package handlers

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Loader creates an handler loading input
func Loader(u services.Unmarshaller) types.HandlerFunc {
	return func(ctx *types.Context) {
		// Unmarshal message
		err := u.Unmarshal(ctx.Msg, ctx.T)
		ctx.Logger.WithFields(log.Fields{
			"ctx.T": ctx.T,
		}).Debug("Unmarshal message")

		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
