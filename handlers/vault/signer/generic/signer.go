package generic

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
)

// TransactionSignerFunc is a generic function interface that support signature with EEA, Tessera, and Ethereum
type TransactionSignerFunc func(keystore.KeyStore, *engine.TxContext, common.Address, *ethtypes.Transaction) ([]byte, *common.Hash, error)

// GenerateSignerHandler creates a signer handler
func GenerateSignerHandler(signerFunc TransactionSignerFunc, backend keystore.KeyStore, successMsg, errorMsg string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chainID": txctx.Builder.GetChainIDString(),
			"from":    txctx.Builder.GetFromString(),
			"id":      txctx.Builder.GetID(),
		})

		if txctx.Builder.GetRaw() != "" {
			// Tx has already been signed
			return
		}

		transaction, err := txctx.Builder.GetTransaction()
		if err != nil {
			txctx.Logger.WithError(err).Errorf(errorMsg)
			_ = txctx.AbortWithError(err)
			return
		}

		from, err := txctx.Builder.GetFromAddress()
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
		_ = txctx.Builder.SetRaw(raw).SetTxHash(*h)

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"raw":    utils.ShortString(hexutil.Encode(raw), 10),
			"txHash": h.Hex(),
		})
		txctx.Logger.Debugf(successMsg)
	}
}
