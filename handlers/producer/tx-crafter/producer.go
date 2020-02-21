package txcrafter

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

	// If an error occurred then we redirect to recovery with a tx response
	switch {
	case !txctx.Envelope.OnlyWarnings():
		msg.Topic = viper.GetString(broker.TxRecoverViperKey)
		p = txctx.Envelope.TxResponse()
	default:
		msg.Topic = viper.GetString(broker.TxSignerViperKey)
		p = txctx.Envelope.TxEnvelopeAsRequest()
	}

	// Marshal Envelope into sarama Message
	err := encoding.Marshal(p, msg)
	if err != nil {
		return err
	}

	msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(txctx.Envelope.ChainID, txctx.Envelope.MustGetFromAddress()))

	return nil
}

// Producer creates a producer handler
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}
