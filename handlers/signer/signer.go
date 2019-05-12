package signer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore"
)

// Signer creates a signer handler
func Signer(s keystore.KeyStore) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chain.id":  txctx.Envelope.GetChain().GetId(),
			"tx.sender": txctx.Envelope.GetSender().GetAddr(),
		})

		if txctx.Envelope.GetTx().GetRaw() != "" {
			// Tx has already been signed
			return
		}

		var t *ethtypes.Transaction
		if txctx.Envelope.GetCall().GetMethod().GetName() == "constructor" {
			// Create contract deployment transaction
			t = ethtypes.NewContractCreation(
				txctx.Envelope.GetTx().GetTxData().GetNonce(),
				txctx.Envelope.GetTx().GetTxData().ValueBig(),
				txctx.Envelope.GetTx().GetTxData().GetGas(),
				txctx.Envelope.GetTx().GetTxData().GasPriceBig(),
				txctx.Envelope.GetTx().GetTxData().DataBytes(),
			)
		} else {
			// Create transaction
			address, err := txctx.Envelope.GetTx().GetTxData().ToAddress()
			if err != nil {
				// TODO: handle error
				txctx.Logger.WithError(err).Warnf("signer: could not get 'to' address from envelope")
				// We indicate that we got an error signing the transaction but we do not abort
				_ = txctx.Error(err)
				return
			}

			t = ethtypes.NewTransaction(
				txctx.Envelope.GetTx().GetTxData().GetNonce(),
				address,
				txctx.Envelope.GetTx().GetTxData().ValueBig(),
				txctx.Envelope.GetTx().GetTxData().GetGas(),
				txctx.Envelope.GetTx().GetTxData().GasPriceBig(),
				txctx.Envelope.GetTx().GetTxData().DataBytes(),
			)
		}

		// Sign transaction
		sender, err := txctx.Envelope.GetSender().Address()
		if err != nil {
			// TODO: handle error
			txctx.Logger.WithError(err).Warnf("signer: could not get sender address from envelope")
			// We indicate that we got an error signing the transaction but we do not abort
			_ = txctx.Error(err)
			return
		}
		raw, h, err := s.SignTx(txctx.Envelope.GetChain(), sender, t)
		if err != nil {
			// TODO: handle error
			txctx.Logger.WithError(err).Warnf("signer: could not sign transaction")
			// We indicate that we got an error signing the transaction but we do not abort
			_ = txctx.Error(err)
			return
		}

		// Update trace information
		enc := hexutil.Encode(raw)
		txctx.Envelope.Tx.SetRaw(enc)
		txctx.Envelope.Tx.SetHash(*h)
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.raw":  enc,
			"tx.hash": h.Hex(),
		})
		txctx.Logger.Debugf("signer: transaction signed")
	}
}
