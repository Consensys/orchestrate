package txlistener

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

// PrepareMsg prepare message to produce from TxContexts
func PrepareMsg(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	// Marshal Envelope into sarama Message
	err := encoding.Marshal(txctx.Envelope, msg)
	if err != nil {
		return err
	}

	// Set Topic to Nonce topic
	msg.Topic = utils.KafkaChainTopic(viper.GetString(broker.TxDecoderViperKey), txctx.Envelope.GetChain().GetBigChainID())

	// Set key
	Sender := txctx.Envelope.Sender()
	msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(txctx.Envelope.GetChain().GetBigChainID(), Sender))

	return nil
}

// Producer creates a producer handler
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {

		// TODO: Make it possible to filter at the last moment. So that we can produce in multiple topics
		if ExternalTxDisabled() && (txctx.Envelope.Metadata == nil || txctx.Envelope.Tx == nil) {

			// For robustness make sure that receipt and tx hash are set even though tx-listener already guarantee it
			if txctx.Envelope.Receipt != nil && txctx.Envelope.Receipt.TxHash != "" {
				txctx.Logger.WithFields(log.Fields{
					"tx.hash": txctx.Envelope.Receipt.GetTxHash(),
				}).Debugf("External tx disabled, skipping transaction with tx hash")
			}

			// Ignore the following handlers
			txctx.Abort()
			return
		}

		// Else normally call the producer as we would normally do
		producer.Producer(p, PrepareMsg)(txctx)
	}
}
