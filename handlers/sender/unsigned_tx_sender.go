package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
)

// UnsignedTxSender creates an handler that send transaction to be signed by ethereum client
func UnsignedTxSender(ec ethclient.TransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txHash, err := ec.SendTransaction(
			txctx.Context(),
			txctx.Envelope.GetChain().ID(),
			types.Envelope2SendTxArgs(txctx.Envelope),
		)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to send unsigned transaction")
			return
		}

		// Transaction has been properly sent so we set tx hash on Envelope
		txctx.Envelope.GetTx().SetHash(txHash)
	}
}
