package sender

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/envelope/storer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

// Sender creates sender handler
func Sender(ec ethclient.TransactionSender, s svc.EnvelopeStoreClient, txSchedulerClient client.TransactionSchedulerClient) engine.HandlerFunc {
	// Declare a set of handlers that will be forked by Sender handler
	rawTxStore := storer.RawTxStore(s, txSchedulerClient)
	UnsignedTxStore := storer.UnsignedTxStore(s, txSchedulerClient)

	rawTxSender := engine.CombineHandlers(
		rawTxStore,
		RawTxSender(ec),
	)

	// Orion private tx
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
			"chainID": txctx.Envelope.GetChainIDString(),
			"raw":     txctx.Envelope.GetShortRaw(),
			"txHash":  txctx.Envelope.GetTxHashString(),
			"from":    txctx.Envelope.GetFromString(),
		})

		switch {
		case txctx.Envelope.IsEthSendRawTransaction() || txctx.Envelope.IsEthSendTransaction():
			rawTxSender(txctx)
		case txctx.Envelope.IsEthSendPrivateTransaction():
			unsignedTxSender(txctx)
		case txctx.Envelope.IsEthSendRawPrivateTransaction():
			tesseraRawPrivateTxSender(txctx)
		case txctx.Envelope.IsEeaSendPrivateTransaction():
			rawPrivateTxSender(txctx)
		default:
			var err error
			// @TODO Remove once envelope store is deleted
			if txctx.Envelope.BelongToEnvelopeStore() {
				err = errors.DataError(
					"invalid transaction protocol %q",
					txctx.Envelope.Method.String(),
				).SetComponent(component)
			} else {
				err = errors.DataError(
					"invalid job type %q",
					txctx.Envelope.JobType.String(),
				).SetComponent(component)
			}

			txctx.Logger.WithError(err).Errorf("sender: could not send transaction")
			_ = txctx.AbortWithError(err)
		}
	}
}
