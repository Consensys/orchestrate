package generic

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore"
)

// TransactionSignerFunc is a generic function interface that support signature with EEA, Tessera, and Ethereum
type TransactionSignerFunc func(keystore.KeyStore, *engine.TxContext, common.Address, *ethtypes.Transaction) ([]byte, *common.Hash, error)

// GenerateSignerHandler creates a signer handler
func GenerateSignerHandler(signerFunc TransactionSignerFunc, backend keystore.KeyStore, successMsg, errorMsg string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chain.id":    txctx.Envelope.GetChain().ID().String(),
			"tx.sender":   txctx.Envelope.GetFrom().Address().Hex(),
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		if txctx.Envelope.GetTx().GetRaw() != nil {
			// Tx has already been signed
			return
		}

		var t = TransactionFromTxContext(txctx)

		// Sign transaction
		sender := txctx.Envelope.GetFrom().Address()
		raw, h, err := signerFunc(backend, txctx, sender, t)
		if err != nil {
			txctx.Logger.WithError(err).Warnf(errorMsg)
			_ = txctx.AbortWithError(err)
			return
		}

		// Update trace information
		txctx.Envelope.Tx.SetRaw(raw)
		txctx.Envelope.Tx.SetHash(*h)
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.raw":  utils.ShortString(hexutil.Encode(raw), 10),
			"tx.hash": h.Hex(),
		})
		txctx.Logger.Debugf(successMsg)
	}
}
