package txdecoder

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// PrepareMsg prepare message to produce from TxContexts
func PrepareMsg(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	// Marshal Envelope into sarama Message
	err := encoding.Marshal(txctx.Envelope.TxResponse(), msg)
	if err != nil {
		return err
	}

	// Set Topic to Nonce topic
	msg.Topic = viper.GetString(broker.TxDecodedViperKey)

	// Set key
	// msg.Key = sarama.ByteEncoder(txctx.Envelope.ChainID.Bytes())

	return nil
}

// Producer creates a producer handler that filters in Orchestrate transaction
// NB: If the transaction
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}
