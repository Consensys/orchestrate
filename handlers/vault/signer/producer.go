package signer

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

// PrepareMsg prepare message to produce from TxContexts
func PrepareMsg(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	// Marshal Envelope into sarama Message
	err := encoding.Marshal(txctx.Envelope, msg)
	if err != nil {
		return err
	}

	// Set Topic at sender by default
	msg.Topic = viper.GetString("kafka.topic.sender")

	if _, ok := txctx.Envelope.GetMetadata().GetExtra()["key-gen-request"]; ok {

		// In case this is a keygen request, we set the output topic to keygen.out
		msg.Topic = viper.GetString("kafka.topic.keygen.out")
	}

	// Set key
	Sender := txctx.Envelope.GetFrom().Address()
	msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(txctx.Envelope.GetChain().ID(), Sender))

	return nil
}

// Producer creates a producer handler
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}
