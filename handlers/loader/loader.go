package loader

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

// Loader is a Middleware enginer.HandlerFunc that Load sarama.ConsumerGroup messages
func Loader(txctx *engine.TxContext) {
	// Cast message into sarama.ConsumerMessage
	receipt, ok := txctx.Msg.(*types.TxListenerReceipt)
	if !ok {
		txctx.Logger.Errorf("loader: expected a types.TxListenerReceipt")
		_ = txctx.AbortWithError(fmt.Errorf("invalid input message format"))
		return
	}

	// Set receipt
	txctx.Envelope.Receipt = ethereum.FromGethReceipt(&receipt.Receipt).
		SetBlockHash(receipt.BlockHash).
		SetBlockNumber(uint64(receipt.BlockNumber)).
		SetTxIndex(receipt.TxIndex)
	txctx.Envelope.Chain = (&common.Chain{}).SetID(receipt.ChainID)

	// Enrich Logger
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"chain.id":   receipt.ChainID.Text(10),
		"tx.hash":    receipt.TxHash.Hex(),
		"block.hash": receipt.BlockHash.Hex(),
	})

	txctx.Logger.Tracef("loader: message loaded: %v", txctx.Envelope.String())
}
