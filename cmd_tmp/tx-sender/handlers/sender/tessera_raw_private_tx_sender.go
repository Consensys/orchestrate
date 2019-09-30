package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

// TesseraRawPrivateTxSender creates an handler that send raw private transactions to a Quorum Tessera client
func TesseraRawPrivateTxSender(ec ethclient.TransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		_, err := ec.SendQuorumRawPrivateTransaction(
			txctx.Context(),
			txctx.Envelope.GetChain().ID(),
			txctx.Envelope.GetTx().GetRaw().GetRaw(),
			types.Call2PrivateArgs(txctx.Envelope.GetArgs()).PrivateFor,
		)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to send quorum raw private transaction")
			return
		}
	}
}
