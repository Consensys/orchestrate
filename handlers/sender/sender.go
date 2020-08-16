package sender

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/envelope/storer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

// Sender creates sender handler
func Sender(ec EthClient, txSchedulerClient client.TransactionSchedulerClient) engine.HandlerFunc {
	// Declare a set of handlers that will be forked by Sender handler
	rawTxStore := storer.RawTxStore(txSchedulerClient)

	rawTxSender := engine.CombineHandlers(
		rawTxStore,
		RawTxSender(ec),
	)

	tesseraPrivateTxSender := engine.CombineHandlers(
		rawTxStore,
		TesseraPrivateTxSender(ec),
	)

	tesseraRawMarkingTxSender := engine.CombineHandlers(
		rawTxStore,
		TesseraRawMarkingTxSender(ec),
	)

	eeaRawPrivateTxSender := engine.CombineHandlers(
		rawTxStore,
		EEAPrivateTxSender(ec),
	)

	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chainID": txctx.Envelope.GetChainIDString(),
			"raw":     txctx.Envelope.GetShortRaw(),
			"txHash":  txctx.Envelope.GetTxHashString(),
			"from":    txctx.Envelope.GetFromString(),
		})

		switch {
		// MERGE ALL IN ONE SINGLE TYPE, SendRaw
		case txctx.Envelope.IsEthSendRawTransaction() ||
			txctx.Envelope.IsEthSendTransaction() ||
			txctx.Envelope.IsEeaSendMarkingTransaction():
			rawTxSender(txctx)
		case txctx.Envelope.IsEthSendTesseraPrivateTransaction():
			tesseraPrivateTxSender(txctx)
		case txctx.Envelope.IsEthSendTesseraMarkingTransaction():
			tesseraRawMarkingTxSender(txctx)
		case txctx.Envelope.IsEeaSendPrivateTransaction():
			eeaRawPrivateTxSender(txctx)
		default:
			err := errors.DataError("invalid job type %q", txctx.Envelope.JobType.String()).SetComponent(component)
			txctx.Logger.WithError(err).Errorf("sender: could not send transaction")
			_ = txctx.AbortWithError(err)
		}
	}
}
