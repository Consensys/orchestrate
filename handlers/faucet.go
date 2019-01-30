package handlers

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Faucet creates a Faucet handler
func Faucet(faucet services.Faucet, creditAmount *big.Int) types.HandlerFunc {
	return func(ctx *types.Context) {
		faucetRequest := &services.FaucetRequest{
			ChainID: ctx.T.Chain().ID,
			Address: *ctx.T.Sender().Address,
			Value:   creditAmount,
		}
		_, _, err := faucet.Credit(context.Background(), faucetRequest)
		if err != nil {
			// TODO: handle error
			ctx.Error(err)
		}
	}
}
