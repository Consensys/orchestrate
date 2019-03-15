package handlers

import (
	"context"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"

)

// Sender creates a Sender handler
func Sender(sender services.TxSender) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		if ctx.T.Tx.GetRaw() == "" {
			// Tx is not ready
			// TODO: handle case
			ctx.Abort()
			return
		}

		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id": ctx.T.Chain.GetId(),
			"tx.raw": ctx.T.Tx.GetRaw(),
		})

		err := sender.SendRawTransaction(context.Background(), ctx.T.Chain.ID(), ctx.T.Tx.GetRaw())
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			ctx.Logger.WithError(err).Errorf("sender: could not send transaction")
			return
		}
		ctx.Logger.WithError(err).Errorf("sender: transaction sent")

	}
}
