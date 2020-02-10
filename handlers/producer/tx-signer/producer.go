package txsigner

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

// PrepareMsg prepare message to produce from TxContexts
func PrepareMsg(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	var p proto.Message

	// Marshal Builder into sarama Message
	switch txctx.In.Entrypoint() {
	case viper.GetString(broker.TxSignerViperKey):
		switch {
		case !txctx.Builder.OnlyWarnings():
			msg.Topic = viper.GetString(broker.TxRecoverViperKey)
			p = txctx.Builder.TxResponse()
		default:
			msg.Topic = viper.GetString(broker.TxSenderViperKey)
			p = txctx.Builder.TxEnvelopeAsRequest()
		}

		// Set key for Kafka partitions
		msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(txctx.Builder.ChainID, txctx.Builder.MustGetFromAddress()))
	case viper.GetString(broker.WalletGeneratorViperKey):
		msg.Topic = viper.GetString(broker.WalletGeneratedViperKey)
		p = txctx.Builder.TxResponse()

		// Set key for Kafka partitions
		msg.Key = sarama.StringEncoder(txctx.Builder.GetFromString())
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
