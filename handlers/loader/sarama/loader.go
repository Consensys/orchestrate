package sarama

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"
)

// Loader is an handler that Load sarama.ConsumerGroup messages
func Loader(txctx *engine.TxContext) {
	// Cast message into sarama.ConsumerMessage
	msg, ok := txctx.In.(*broker.Msg)
	if !ok {
		txctx.Logger.Fatalf("loader: expected a sarama.ConsumerMessage")
	}

	// Enrich Logger
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"kafka.in.topic":     msg.Topic,
		"kafka.in.offset":    msg.Offset,
		"kafka.in.partition": msg.Partition,
	})

	switch txctx.In.Entrypoint() {
	case viper.GetString(broker.TxDecoderViperKey), viper.GetString(broker.TxDecodedViperKey), viper.GetString(broker.TxRecoverViperKey), viper.GetString(broker.WalletGeneratedViperKey):
		loadTxResponse(txctx)
	default:
		loadTxRequest(txctx)
	}
}

func loadTxEnvelope(txctx *engine.TxContext) {
	txEnvelope := &tx.TxEnvelope{}

	err := encoding.Unmarshal(txctx.In.(*broker.Msg), txEnvelope)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("loader: error unmarshalling")
		return
	}
	txctx.Logger.Tracef("loader: tx envelope loaded: %v", txEnvelope)

	envelope, err := txEnvelope.Envelope()
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("loader: invalid tx envelope")
		return
	}

	txctx.Envelope = envelope
}

func loadTxRequest(txctx *engine.TxContext) {
	txRequest := &tx.TxRequest{}
	err := encoding.Unmarshal(txctx.In.(*broker.Msg), txRequest)
	if err != nil {
		loadTxEnvelope(txctx)
		return
	}
	txctx.Logger.Tracef("loader: tx request loaded: %v", txRequest)

	envelope, err := txRequest.Envelope()
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("loader: invalid tx request")
		return
	}

	txctx.Envelope = envelope
}

func loadTxResponse(txctx *engine.TxContext) {
	txResponse := &tx.TxResponse{}
	err := encoding.Unmarshal(txctx.In.(*broker.Msg), txResponse)
	if err != nil {
		loadTxEnvelope(txctx)
		return
	}
	txctx.Logger.Tracef("loader: tx response loaded: %v", txResponse)

	envelope, err := txResponse.Envelope()
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("loader: invalid tx response")
		return
	}

	txctx.Envelope = envelope
}
