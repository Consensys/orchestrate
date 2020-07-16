package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

// RawTxSender creates an handler that send raw transactions to an Ethereum client
func RawTxSender(ec ethclient.TransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		if txctx.Envelope.GetRaw() == "" {
			err := errors.DataError("no raw filled")
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}

		txHash, err := ec.SendRawTransaction(
			txctx.Context(),
			url,
			txctx.Envelope.GetRaw(),
		)

		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to send raw transaction")
			return
		}

		if txHash.String() != txctx.Envelope.TxHash.String() {
			err := errors.InternalError("invalid transaction Hash. Expected %s, got %s",
				txctx.Envelope.TxHash.String(), txHash.String())
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}
	}
}
