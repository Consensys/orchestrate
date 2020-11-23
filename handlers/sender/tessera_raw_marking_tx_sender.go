package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/proxy"
)

// TesseraRawMarkingTxSender creates an handler that send raw private transactions to a Quorum Tessera client
func TesseraRawMarkingTxSender(ec ethclient.QuorumTransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		if txctx.Envelope.Raw == "" || len(txctx.Envelope.PrivateFor) == 0 {
			err := errors.DataError("no raw or privateFor filled")
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}

		txHash, err := ec.SendQuorumRawPrivateTransaction(
			txctx.Context(),
			url,
			txctx.Envelope.Raw,
			types.Call2PrivateArgs(txctx.Envelope).PrivateFor,
		)

		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to send quorum raw private transaction")
			return
		}

		// Set TxHash to be the newly returned one instead of the computed one
		_ = txctx.Envelope.SetTxHash(txHash)
	}
}
