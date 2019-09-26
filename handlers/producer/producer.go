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
	// Set message Key
	msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(txctx.Envelope.GetChain().ID(), txctx.Envelope.Sender()))

	// We loop over errors to
	// - determine to which topic to send transaction
	// - remove invalid nonce warnings
loop:
	for _, err := range txctx.Envelope.GetErrors() {
		switch {
		case errors.IsWarning(err):
			continue
		default:
			// If an error occurred we redirect to recovery
			msg.Topic = viper.GetString("kafka.topic.recover")
			break loop
		}
	}

	// If no error and nonce is invalid we redirect envelope to tx-nonce
	if b, ok := txctx.Get("invalid.nonce").(bool); len(txctx.Envelope.GetErrors()) == 0 && ok && b {
		msg.Topic = viper.GetString("kafka.topic.nonce")
	}

	// Marshal Envelope into sarama Message
	err := encoding.Marshal(txctx.Envelope, msg)
	if err != nil {
		return err
	}

	return nil
}

// Producer creates a producer handler
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}
