package signer

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

// TransactionSignerFunc is a generic function interface that support signature with EEA, Tessera, and Ethereum
type TransactionSignerFunc func(keystore.KeyStore, *engine.TxContext, ethcommon.Address, *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error)

// GenerateSignerHandler creates a signer handler
func GenerateSignerHandler(signerFunc TransactionSignerFunc, vks, onetime keystore.KeyStore, successMsg, errorMsg string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chainID": txctx.Envelope.GetChainIDString(),
			"from":    txctx.Envelope.GetFromString(),
			"id":      txctx.Envelope.GetID(),
		})

		if txctx.Envelope.GetRaw() != "" {
			// Tx has already been signed
			return
		}

		transaction, err := txctx.Envelope.GetTransaction()
		if err != nil {
			txctx.Logger.WithError(err).Errorf(errorMsg)
			_ = txctx.AbortWithError(err)
			return
		}

		var backend keystore.KeyStore
		var from ethcommon.Address
		if txctx.Envelope.IsOneTimeKeySignature() {
			backend = onetime
			from, err = onetime.GenerateAccount(txctx.Context())
			_ = txctx.Envelope.SetFrom(from)
		} else {
			backend = vks
			from, err = txctx.Envelope.GetFromAddress()
		}

		if err != nil {
			txctx.Logger.WithError(err).Errorf(errorMsg)
			_ = txctx.AbortWithError(err)
			return
		}

		// Sign transaction
		raw, h, err := signerFunc(backend, txctx, from, transaction)
		if err != nil {
			txctx.Logger.WithError(err).Warnf(errorMsg)
			_ = txctx.AbortWithError(err)
			return
		}

		// Update trace information
		_ = txctx.Envelope.SetRaw(raw).SetTxHash(*h)

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"raw":    utils.ShortString(hexutil.Encode(raw), 10),
			"txHash": h.Hex(),
		})
		txctx.Logger.Debugf(successMsg)
	}
}
