package sender

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/envelope/storer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
)

// Sender creates sender handler
func Sender(ec ethclient.TransactionSender, s evlpstore.EnvelopeStoreClient) engine.HandlerFunc {
	// Declare a set of handlers that will be forked by Sender handler
	rawTxStore := storer.RawTxStore(s)
	UnsignedTxStore := storer.UnsignedTxStore(s)

	rawTxSender := engine.CombineHandlers(
		rawTxStore,
		RawTxSender(ec),
	)

	rawPrivateTxSender := engine.CombineHandlers(
		UnsignedTxStore,
		RawPrivateTxSender(ec),
	)

	tesseraRawPrivateTxSender := engine.CombineHandlers(
		rawTxStore,
		TesseraRawPrivateTxSender(ec),
	)

	unsignedTxSender := engine.CombineHandlers(
		UnsignedTxStore,
		UnsignedTxSender(ec),
	)

	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chainID": txctx.Builder.GetChainIDString(),
			"id":      txctx.Builder.GetID(),
			"tx.raw":  txctx.Builder.GetShortRaw(),
			"tx.hash": txctx.Builder.GetTxHashString(),
			"from":    txctx.Builder.GetFromString(),
		})

		switch {
		case txctx.Builder.IsEthSendRawTransaction():
			rawTxSender(txctx)
		case txctx.Builder.IsEthSendPrivateTransaction():
			unsignedTxSender(txctx)
		case txctx.Builder.IsEthSendRawPrivateTransaction():
			tesseraRawPrivateTxSender(txctx)
		case txctx.Builder.IsEeaSendPrivateTransaction():
			rawPrivateTxSender(txctx)
		default:
			err := errors.DataError(
				"invalid private protocol %q",
				txctx.Builder.Method.String(),
			).SetComponent(component)
			txctx.Logger.WithError(err).Errorf("sender: could not send private transaction")
			_ = txctx.AbortWithError(err)
		}
	}
}
