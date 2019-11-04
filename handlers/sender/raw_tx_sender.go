package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// RawTxSender creates an handler that send raw transactions to an Ethereum client
func RawTxSender(ec ethclient.TransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		err := ec.SendRawTransaction(
			txctx.Context(),
			txctx.Envelope.GetChain().ID(),
			txctx.Envelope.GetTx().GetRaw().Hex(),
		)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to send raw transaction")
			return
		}
	}
}
