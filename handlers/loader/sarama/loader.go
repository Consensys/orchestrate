package sarama

import (
	broker "github.com/consensys/orchestrate/pkg/broker/sarama"
	encoding "github.com/consensys/orchestrate/pkg/encoding/sarama"
	"github.com/consensys/orchestrate/pkg/engine"
	"github.com/consensys/orchestrate/pkg/types/tx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	case viper.GetString(broker.TxDecodedViperKey), viper.GetString(broker.TxRecoverViperKey):
		loadTxResponse(txctx)
	default:
		loadTxRequest(txctx)
	}

	if txctx.Envelope != nil {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"id":       txctx.Envelope.GetID(),
			"job_uuid": txctx.Envelope.GetJobUUID(),
		})
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
		loadTxEnvelope(txctx)
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

	var envelope *tx.Envelope

	// If it's an external tx
	if txResponse.GetJobUUID() == "" {
		envelope = txResponse.ExternalTxEnvelope()
	} else {
		envelope, err = txResponse.Envelope()
	}
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("loader: failed to cast TxResponse to Envelope")
		return
	}

	txctx.Envelope = envelope
}
