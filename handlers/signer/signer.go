package signer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore"
)

// Signer creates a signer handler
func Signer(s keystore.KeyStore) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chain.id":  txctx.Envelope.GetChain().GetId(),
			"tx.sender": txctx.Envelope.GetFrom().Address(),
		})

		if txctx.Envelope.GetTx().GetRaw() != nil {
			// Tx has already been signed
			return
		}

		var t *ethtypes.Transaction
		if txctx.Envelope.GetArgs().GetCall().IsConstructor() {
			// Create contract deployment transaction
			t = ethtypes.NewContractCreation(
				txctx.Envelope.GetTx().GetTxData().GetNonce(),
				txctx.Envelope.GetTx().GetTxData().GetValueBig(),
				txctx.Envelope.GetTx().GetTxData().GetGas(),
				txctx.Envelope.GetTx().GetTxData().GetGasPriceBig(),
				txctx.Envelope.GetTx().GetTxData().GetDataBytes(),
			)
		} else {
			// Create transaction
			address := txctx.Envelope.GetTx().GetTxData().GetTo().Address()

			t = ethtypes.NewTransaction(
				txctx.Envelope.GetTx().GetTxData().GetNonce(),
				address,
				txctx.Envelope.GetTx().GetTxData().GetValueBig(),
				txctx.Envelope.GetTx().GetTxData().GetGas(),
				txctx.Envelope.GetTx().GetTxData().GetGasPriceBig(),
				txctx.Envelope.GetTx().GetTxData().GetDataBytes(),
			)
		}

		// Sign transaction
		sender := txctx.Envelope.GetFrom().Address()
		raw, h, err := s.SignTx(txctx.Envelope.GetChain(), sender, t)
		if err != nil {
			// TODO: handle error
			txctx.Logger.WithError(err).Warnf("signer: could not sign transaction")
			// We indicate that we got an error signing the transaction but we do not abort
			_ = txctx.Error(err)
			return
		}

		// Update trace information
		txctx.Envelope.Tx.SetRaw(raw)
		txctx.Envelope.Tx.SetHash(*h)
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.raw":  utils.ShortString(hexutil.Encode(raw), 10),
			"tx.hash": h.Hex(),
		})
		txctx.Logger.Debugf("signer: transaction signed")
	}
}
