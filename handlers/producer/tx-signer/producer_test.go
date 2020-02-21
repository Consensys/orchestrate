package txsigner

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type TestMsg string

func (msg TestMsg) Entrypoint() string    { return string(msg) }
func (msg TestMsg) Header() engine.Header { return &header{} }
func (msg TestMsg) Value() []byte         { return []byte{} }
func (msg TestMsg) Key() []byte           { return []byte{} }

type header struct{}

func (h *header) Add(key, value string) {}
func (h *header) Del(key string)        {}
func (h *header) Get(key string) string { return "" }
func (h *header) Set(key, value string) {}

func TestPrepareMsgSigner(t *testing.T) {
	// No error
	txctx := engine.NewTxContext()

	txctx.In = TestMsg("topic-tx-signer")
	msg := &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-sender", msg.Topic, "If no error out topic should be sender")

	// Classic error
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-recover", msg.Topic, "If error out topic should be recovery")
}

func TestPrepareMsgGenerateAccount(t *testing.T) {
	// No error
	txctx := engine.NewTxContext()

	txctx.In = TestMsg("topic-account-generator")
	msg := &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-account-generated", msg.Topic, "If no error out topic should be account-generated")

	// Classic error
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-account-generated", msg.Topic, "If error out topic should be recovery")
}
