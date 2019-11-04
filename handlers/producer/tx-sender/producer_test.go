package txsender

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
)

func TestPrepareMsg(t *testing.T) {
	// No error
	txctx := engine.NewTxContext()
	txctx.Envelope = &envelope.Envelope{}
	msg := &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "", msg.Topic, "If no error there should be no out topic")

	// Error
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-recover", msg.Topic, "If error out topic should be recovery")

	txctx = engine.NewTxContext()
	txctx.Envelope = &envelope.Envelope{}
	txctx.Set("invalid.nonce", true)
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-nonce", msg.Topic, "If invalid nonce out topic should be tx-nonce")

	txctx = engine.NewTxContext()
	txctx.Envelope = &envelope.Envelope{}
	txctx.Set("invalid.nonce", true)
	_ = txctx.Error(errors.ConnectionError("nonce too low"))
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-recover", msg.Topic, "If invalid nonce and error topic should be tx-recover")
}
