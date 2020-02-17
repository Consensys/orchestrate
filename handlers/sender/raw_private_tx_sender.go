package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

// RawPrivateTxSender creates an handler that send raw private transactions to an Ethereum client
func RawPrivateTxSender(ec ethclient.TransactionSender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		if txctx.Envelope.Raw == "" {
			err := errors.DataError("no raw filled")
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}

		txHash, err := ec.SendRawPrivateTransaction(
			txctx.Context(),
			url,
			txctx.Envelope.Raw,
			types.Call2PrivateArgs(txctx.Envelope),
		)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: failed to send raw private transaction")
			return
		}

		// Transaction has been properly sent so we set tx hash on Envelope
		_ = txctx.Envelope.SetTxHash(txHash)
	}
}
