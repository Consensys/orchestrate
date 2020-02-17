package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

// UnsignedTxSender creates an handler that send transaction to be signed by ethereum client
func UnsignedTxSender(ec ethclient.TransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		txArgs, err := types.Envelope2SendTxArgs(txctx.Envelope)
		if err != nil {
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}

		txHash, err := ec.SendTransaction(
			txctx.Context(),
			url,
			txArgs,
		)
		if err != nil {
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}

		// Transaction has been properly sent so we set tx hash on Envelope
		_ = txctx.Envelope.SetTxHash(txHash)
	}
}
