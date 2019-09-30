package producer

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
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

	// Set Topic to Nonce topic
	msg.Topic = viper.GetString("kafka.topic.signer")

	// If an error occurred then we redirect to recovery
	for _, err := range txctx.Envelope.GetErrors() {
		if !errors.IsWarning(err) {
			msg.Topic = viper.GetString("kafka.topic.recover")
			break
		}
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
