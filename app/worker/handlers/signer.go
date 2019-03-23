package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Signer creates a signer handler
func Signer(s keystore.KeyStore) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id":  ctx.T.GetChain().GetId(),
			"tx.sender": ctx.T.GetSender().GetAddr(),
		})

		if ctx.T.GetTx().GetRaw() != "" {
			// Tx already signed
			return
		}

		// Create transaction
		t := ethtypes.NewTransaction(
			ctx.T.GetTx().GetTxData().GetNonce(),
			ctx.T.GetTx().GetTxData().ToAddress(),
			ctx.T.GetTx().GetTxData().ValueBig(),
			ctx.T.GetTx().GetTxData().GetGas(),
			ctx.T.GetTx().GetTxData().GasPriceBig(),
			ctx.T.GetTx().GetTxData().DataBytes(),
		)

		// Sign transaction
		raw, h, err := s.SignTx(ctx.T.GetChain(), ctx.T.GetSender().Address(), t)
		EncodedRaw := hexutil.Encode(raw)

		if err != nil {
			// TODO: handle error
			ctx.Logger.WithError(err).Infof("signer: could not sign transaction")
			ctx.AbortWithError(err)
			return
		}

		// Update trace information
		ctx.T.Tx.SetRaw(EncodedRaw)
		ctx.T.Tx.SetHash(*h)
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tx.raw":  EncodedRaw,
			"tx.hash": h.Hex(),
		})
		ctx.Logger.Debugf("signer: raw transaction set")
	}
}
