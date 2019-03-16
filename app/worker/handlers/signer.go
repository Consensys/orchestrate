package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Signer creates a signer handler
func Signer(s keystore.KeyStore) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id":  ctx.T.Chain.GetId(),
			"tx.sender": ctx.T.Sender.GetAddr(),
		})

		if ctx.T.Tx.Raw != "" {
			// Tx already signed
			return
		}

		t := ethtypes.NewTransaction(
			ctx.T.Tx.TxData.GetNonce(),
			ctx.T.Tx.TxData.ToAddress(),
			ctx.T.Tx.TxData.ValueBig(),
			ctx.T.Tx.TxData.GetGas(),
			ctx.T.Tx.TxData.GasPriceBig(),
			ctx.T.Tx.TxData.DataBytes(),
		)

		// Sign transaction
		raw, h, err := s.SignTx(ctx.T.Chain, ctx.T.Sender.Address(), t)
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
