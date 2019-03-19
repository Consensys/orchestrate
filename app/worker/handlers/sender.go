package handlers

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Sender creates a Sender handler
func Sender(sender services.TxSender, store infra.TraceStore) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		if ctx.T.GetTx().GetRaw() == "" || ctx.T.GetTx().GetHash() == "" {
			// Tx is not ready
			// TODO: handle case
			ctx.Abort()
			return
		}

		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id": ctx.T.GetChain().GetId(),
			"tx.raw":   ctx.T.GetTx().GetRaw(),
		})

		// Store trace
		status, _, err := store.Store(context.Background(), ctx.T)
		if err != nil {
			// Connexion to store is broken
			ctx.AbortWithError(err)
			return
		}

		if status == "pending" {
			// Tx has already been sent
			// TODO: Still request Tx from chain to make sure we do not miss a message
			ctx.Abort()
			return
		}

		err = sender.SendRawTransaction(context.Background(), ctx.T.GetChain().ID(), ctx.T.GetTx().GetRaw())
		if err != nil {
			ctx.Logger.WithError(err).Errorf("sender: could not send transaction")
			// TODO: handle error
			ctx.Error(err)

			// We update status in storage
			err := store.SetStatus(context.Background(), ctx.T.GetMetadata().GetId(), "error")
			if err != nil {
				// Connexion to store is broken
				ctx.Error(err)
			}

			ctx.Abort()
			return
		}

		// Transaction has been properly sent so we set status to `pending`
		err = store.SetStatus(context.Background(), ctx.T.GetMetadata().GetId(), "pending")
		if err != nil {
			// Connexion to store is broken
			ctx.AbortWithError(err)
			return
		}
		ctx.Logger.WithError(err).Errorf("sender: transaction sent")
	}
}
