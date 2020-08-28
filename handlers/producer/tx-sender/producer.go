package txsender

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

	b, ok := txctx.Get("invalid.nonce").(bool)
	switch {
	case len(txctx.Envelope.GetErrors()) == 0 && ok && b:
		// If nonce is invalid we redirect envelope to tx-crafter
		msg.Topic = viper.GetString(broker.TxCrafterViperKey)
		p = txctx.Envelope.TxEnvelopeAsRequest()
	case !txctx.Envelope.OnlyWarnings():
		msg.Topic = viper.GetString(broker.TxRecoverViperKey)
		p = txctx.Envelope.TxResponse()
	default:
		// Not sending msg in kafka if no error or warnings
		return nil
	}

	// Marshal Envelope into sarama Message
	err := encoding.Marshal(p, msg)
	if err != nil {
		return err
	}

	// Set message Key
	if partitionKey := txctx.Envelope.PartitionKey(); partitionKey != "" {
		msg.Key = sarama.StringEncoder(partitionKey)
	}
	return nil
}

// Producer creates a producer handler
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}
