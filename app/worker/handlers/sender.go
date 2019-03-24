package handlers

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Sender creates a Sender handler
func Sender(sender ethclient.TxSender, store infra.TraceStore) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id": ctx.T.GetChain().GetId(),
		})

		if ctx.T.GetTx().GetRaw() == "" || ctx.T.GetTx().GetHash() == "" {
			// Transaction has not been signed externally
			// Send the transaction
			args := ethclient.Trace2SendTxArgs(ctx.T)
			txHash, err := sender.SendTransaction(ctx.Context(), ctx.T.GetChain().ID(), args)
			if err != nil {
				ctx.Logger.WithError(err).Errorf("sender: could not send transaction")
				ctx.AbortWithError(err)
				return
			}

			// Set transaction Hash on trace
			ctx.T.GetTx().SetHash(txHash)
			ctx.Logger = ctx.Logger.WithFields(log.Fields{
				"tx.hash": ctx.T.GetTx().GetHash(),
			})
			ctx.Logger.Debugf("sender: transaction sent")

			// Store trace
			// We can not store trace before sending transaction because we do not know the transaction hash
			// This is an issue for overall consistency of the system before/after transaction is mined
			ctx.Logger.Infof("%v %v %v", ctx.T.Chain.Id, ctx.T.Tx.Hash, ctx.T.Metadata.Id)
			_, _, err = store.Store(ctx.Context(), ctx.T)
			if err != nil {
				// Connexion to store is broken
				ctx.Logger.WithError(err).Errorf("sender: trace store failed to store trace")
				ctx.AbortWithError(err)
				return
			}

			// Transaction has been properly sent so we set status to `pending`
			err = store.SetStatus(ctx.Context(), ctx.T.GetMetadata().GetId(), "pending")
			if err != nil {
				// Connexion to store is broken
				ctx.Logger.WithError(err).Errorf("sender: piou trace store failed to set status")
				ctx.Error(err)
				return
			}

			return
		}

		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tx.raw":  ctx.T.GetTx().GetRaw(),
			"tx.hash": ctx.T.GetTx().GetHash(),
		})

		// Store trace
		status, _, err := store.Store(ctx.Context(), ctx.T)
		if err != nil {
			// Connexion to store is broken
			ctx.Logger.WithError(err).Errorf("sender: trace store failed to store trace")
			ctx.AbortWithError(err)
			return
		}

		if status == "pending" {
			// Tx has already been sent
			// TODO: Request TxHash from chain to make sure we do not miss a message
			ctx.Logger.Warnf("sender: transaction has already been sent")
			ctx.Abort()
			return
		}

		// Send raw transaction
		err = sender.SendRawTransaction(ctx.Context(), ctx.T.GetChain().ID(), ctx.T.GetTx().GetRaw())
		if err != nil {
			ctx.Logger.WithError(err).Errorf("sender: could not send transaction")

			// TODO: handle error
			ctx.Error(err)

			// We update status in storage
			err := store.SetStatus(ctx.Context(), ctx.T.GetMetadata().GetId(), "error")
			if err != nil {
				// Connexion to store is broken
				ctx.Logger.WithError(err).Errorf("sender: trace store failed to set status")
				ctx.Error(err)
			}
			ctx.Abort()
			return
		}
		ctx.Logger.Debugf("sender: raw transaction sent")

		// Transaction has been properly sent so we set status to `pending`
		err = store.SetStatus(ctx.Context(), ctx.T.GetMetadata().GetId(), "pending")
		if err != nil {
			// Connexion to store is broken
			ctx.Logger.WithError(err).Errorf("sender: trace store failed to set status")
			ctx.Error(err)
			return
		}

	}
}
