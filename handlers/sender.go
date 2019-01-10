package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Sender creates a Sender handler
func Sender(sender infra.TxSender) types.HandlerFunc {
	return func(ctx *types.Context) {
		if len(ctx.T.Tx().Raw()) == 0 {
			// Tx is not ready
			// TODO: handle case
			ctx.Abort()
			return
		}

		err := sender.Send(ctx.T.Chain().ID, hexutil.Encode(ctx.T.Tx().Raw()))
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
