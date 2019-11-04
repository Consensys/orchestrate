package txsigner

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

// PrepareMsg prepare message to produce from TxContexts
func PrepareMsg(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	// Marshal Envelope into sarama Message
	err := encoding.Marshal(txctx.Envelope, msg)
	if err != nil {
		return err
	}

	switch txctx.In.Entrypoint() {
	case viper.GetString("topic.tx.signer"):
		// Set Topic at sender by default
		msg.Topic = viper.GetString("topic.tx.sender")

		// If an error occurred then we redirect to recovery
		for _, err := range txctx.Envelope.GetErrors() {
			if !errors.IsWarning(err) {
				msg.Topic = viper.GetString("topic.tx.recover")
				break
			}
		}

		// Set key for Kafka partitions
		Sender := txctx.Envelope.GetFrom().Address()
		msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(txctx.Envelope.GetChain().ID(), Sender))

		return nil
	case viper.GetString("topic.wallet.generator"):
		msg.Topic = viper.GetString("topic.wallet.generated")

		// Set key for Kafka partitions
		msg.Key = sarama.StringEncoder(txctx.Envelope.GetFrom().Address().Hex())

		return nil
	}

	return nil
}

// Producer creates a producer handler### Version 0.5.2
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}
