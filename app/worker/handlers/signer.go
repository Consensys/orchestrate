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
			// Tx has already been signed
			return
		}

		var t *ethtypes.Transaction
		if ctx.T.GetCall().GetMethod().GetName() == "constructor" {
			// Create contract deployment transaction
			t = ethtypes.NewContractCreation(
				ctx.T.GetTx().GetTxData().GetNonce(),
				ctx.T.GetTx().GetTxData().ValueBig(),
				ctx.T.GetTx().GetTxData().GetGas(),
				ctx.T.GetTx().GetTxData().GasPriceBig(),
				ctx.T.GetTx().GetTxData().DataBytes(),
			)
		} else {
			// Create transaction
			t = ethtypes.NewTransaction(
				ctx.T.GetTx().GetTxData().GetNonce(),
				ctx.T.GetTx().GetTxData().ToAddress(),
				ctx.T.GetTx().GetTxData().ValueBig(),
				ctx.T.GetTx().GetTxData().GetGas(),
				ctx.T.GetTx().GetTxData().GasPriceBig(),
				ctx.T.GetTx().GetTxData().DataBytes(),
			)
		}

		// Sign transaction
		raw, h, err := s.SignTx(ctx.T.GetChain(), ctx.T.GetSender().Address(), t)
		if err != nil {
			// TODO: handle error
			ctx.Logger.WithError(err).Warnf("signer: could not sign transaction")
			// We indicate that we got an error signing the transaction but we do not abort
			ctx.Error(err)
			return
		}

		// Update trace information
		enc := hexutil.Encode(raw)
		ctx.T.Tx.SetRaw(enc)
		ctx.T.Tx.SetHash(*h)
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tx.raw":  enc,
			"tx.hash": h.Hex(),
		})
		ctx.Logger.Debugf("signer: transaction signed")
	}
}
