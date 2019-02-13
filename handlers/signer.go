package handlers

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
)

// HandlerSignature implements the interface TxSigner
func HandlerSignature(s services.KeyStore) types.HandlerFunc {

	return func(ctx *types.Context) {

		if len(ctx.T.Tx().Raw()) > 0 {
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
		raw, h, err := s.SignTx(ctx.T.Chain(), *ctx.T.Sender().Address, t)

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