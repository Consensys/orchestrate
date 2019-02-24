package handlers

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Sender creates a Sender handler
func Sender(sender services.TxSender) types.HandlerFunc {
	return func(ctx *types.Context) {
		if len(ctx.T.Tx().Raw()) == 0 {
			// Tx is not ready
			// TODO: handle case
			ctx.Abort()
			return
		}

		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id": ctx.T.Chain().ID.Text(16),
			"tx.raw": hexutil.Encode(ctx.T.Tx().Raw()),
		})

		err := sender.Send(context.Background(), ctx.T.Chain().ID, hexutil.Encode(ctx.T.Tx().Raw()))
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			ctx.Logger.WithError(err).Errorf("sender: could not send transaction")
			return
		}
		ctx.Logger.WithError(err).Errorf("sender: transaction sent")

	}
}
