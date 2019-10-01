package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
)

// RawPrivateTxSender creates an handler that send raw private transactions to an Ethereum client
func RawPrivateTxSender(ec ethclient.TransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txHash, err := ec.SendRawPrivateTransaction(
			txctx.Context(),
			txctx.Envelope.GetChain().ID(),
			txctx.Envelope.GetTx().GetRaw().GetRaw(),
			types.Call2PrivateArgs(txctx.Envelope.GetArgs()),
		)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to send raw private transaction")
			return
		}

		// Transaction has been properly sent so we set tx hash on Envelope
		txctx.Envelope.GetTx().SetHash(txHash)
	}
}
