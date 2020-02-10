package receipt

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

// Loader is a Middleware engine.HandlerFunc that Load sarama.ConsumerGroup messages
func Loader(txctx *engine.TxContext) {
	// Cast message into sarama.ConsumerMessage
	receipt, ok := txctx.In.(*types.TxListenerReceipt)
	if !ok {
		_ = txctx.AbortWithError(errors.InternalError("invalid input message format")).
			SetComponent(component)
		return
	}

	// Set receipt
	txctx.Builder.Receipt = ethereum.FromGethReceipt(&receipt.Receipt).
		SetBlockHash(receipt.BlockHash).
		SetBlockNumber(uint64(receipt.BlockNumber)).
		SetTxIndex(receipt.TxIndex)
	txctx.Builder.ChainID = receipt.ChainID

	// Enrich Logger
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"chain.chainID": receipt.ChainID.Text(10),
		"tx.hash":       receipt.TxHash.Hex(),
		"block.hash":    receipt.BlockHash.Hex(),
	})

	txctx.Logger.Tracef("loader: message loaded: %v", txctx.Builder)
}
