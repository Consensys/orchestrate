// +build unit

package sarama

import (
	"reflect"
	"testing"

	broker "github.com/consensys/orchestrate/pkg/broker/sarama"
	"github.com/consensys/orchestrate/pkg/engine"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/tx"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestLoader(t *testing.T) {
	testSet := []struct {
		name          string
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"Loader without error",
			func(txctx *engine.TxContext) *engine.TxContext {
				b := tx.NewEnvelope().SetID("dce80ed3-8b0e-4045-9a91-832ba0391c44")
				msg := &broker.Msg{}
				msg.ConsumerMessage.Value, _ = proto.Marshal(b.TxRequest())
				txctx.In = msg
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.ID = "dce80ed3-8b0e-4045-9a91-832ba0391c44"
				return txctx
			},
		},
		{
			"Loader with error when unmarshalling envelope",
			func(txctx *engine.TxContext) *engine.TxContext {
				msg := &broker.Msg{ConsumerMessage: sarama.ConsumerMessage{Value: []byte{1}}}
				txctx.In = msg
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.EncodingError("proto: envelope.Envelope: illegal tag 0 (wire type 1)").ExtendComponent("handler.loader.encoding.sarama")
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewLogger()

			Loader(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx.Envelope.InternalLabels, expectedTxctx.Envelope.InternalLabels), "Expected same input")
		})
	}
}
