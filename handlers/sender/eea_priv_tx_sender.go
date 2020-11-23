package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/proxy"
)

func EEAPrivateTxSender(ec ethclient.EEATransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		url, err := proxy.GetURL(txctx)
		if err != nil {
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(err).Errorf("sender: failed to get chain URL")
			return
		}

		if txctx.Envelope.Raw == "" {
			err := errors.DataError("sender: empty raw field filled")
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}

		txHash, err := ec.PrivDistributeRawTransaction(txctx.Context(), url, txctx.Envelope.Raw)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to send eea raw private transaction")
			return
		}

		// Set TxHash to be the newly returned one instead of the computed one
		_ = txctx.Envelope.SetTxHash(txHash)

		txctx.Logger.WithField("txHash", txHash.String()).
			Warnf("sender: private tx was sent")
	}
}
