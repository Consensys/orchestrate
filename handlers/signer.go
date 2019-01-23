package handlers

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Signer creates a signer handler
func Signer(s services.TxSigner) types.HandlerFunc {
	return func(ctx *types.Context) {
		if len(ctx.T.Tx().Raw()) > 0 {
			// Tx already signed
			return
		}

		t := ethtypes.NewTransaction(
			ctx.T.Tx().Nonce(),
			*ctx.T.Tx().To(),
			ctx.T.Tx().Value(),
			ctx.T.Tx().GasLimit(),
			ctx.T.Tx().GasPrice(),
			ctx.T.Tx().Data(),
		)

		// Sign transaction
		raw, h, err := s.Sign(ctx.T.Chain(), *ctx.T.Sender().Address, t)

		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}

		// Update trace information
		ctx.T.Tx().SetRaw(raw)
		ctx.T.Tx().SetHash(h)
	}
}
