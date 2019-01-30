package handlers

import (
	"context"

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

		err := sender.Send(context.Background(), ctx.T.Chain().ID, hexutil.Encode(ctx.T.Tx().Raw()))
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
