package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

// TesseraRawPrivateTxSender creates an handler that send raw private transactions to a Quorum Tessera client
func TesseraRawPrivateTxSender(ec ethclient.TransactionSender) engine.HandlerFunc {
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

		if txHash.String() != txctx.Envelope.TxHash.String() {
			err := errors.DataError("invalid generate txHash. Expected %s, got %s",
				txctx.Envelope.TxHash.String(), txHash.String())
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}
	}
}
