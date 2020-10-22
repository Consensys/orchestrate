package txsigner

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// PrepareMsg prepare message to produce from TxContexts
func PrepareMsg(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	var p proto.Message

	// Marshal Envelope into sarama Message
	if txctx.In.Entrypoint() == viper.GetString(broker.TxSignerViperKey) {
		switch {
		case txctx.Envelope.OnlyWarnings():
			msg.Topic = viper.GetString(broker.TxSenderViperKey)
			p = txctx.Envelope.TxEnvelopeAsRequest()
		case !txctx.Envelope.OnlyWarnings() && txctx.Envelope.IsParentJob():
			msg.Topic = viper.GetString(broker.TxRecoverViperKey)
			p = txctx.Envelope.TxResponse()
		default:
			return nil
		}

		// Set key for Kafka partitions
		if partitionKey := txctx.Envelope.PartitionKey(); partitionKey != "" {
			msg.Key = sarama.StringEncoder(partitionKey)
		}
	}

	err := encoding.Marshal(p, msg)
	if err != nil {
		return err
	}

	return nil
}

// Producer creates a producer handler### Version 0.5.2
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}
