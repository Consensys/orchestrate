package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Signer creates a signer handler
func Signer(s services.TxSigner) types.HandlerFunc {
	return func(ctx *types.Context) {
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id":  ctx.T.Chain().ID.Text(16),
			"tx.sender": ctx.T.Sender().Address.Hex(),
		})

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
			ctx.Logger.WithError(err).Infof("signer: could not sign transaction")
			ctx.AbortWithError(err)
			return
		}

		// Update trace information
		ctx.T.Tx().SetRaw(raw)
		ctx.T.Tx().SetHash(h)
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tw.raw":  hexutil.Encode(raw),
			"tx.hash": h.Hex(),
		})
		ctx.Logger.Debugf("signer: raw transaction set")
	}
}
