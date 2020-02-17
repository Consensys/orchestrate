package txdecoder

import (
	"math/big"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

type MockSyncProducer struct {
	Produced chan string
}

const justProducedAMsg = "Just produced a message"

func NewMockSyncProducer() *MockSyncProducer {
	return &MockSyncProducer{
		// Buffering 32 slots in the buffer to avoid deadlocks
		Produced: make(chan string, 32),
	}
}

// SendMessage produces a given message, and returns only when it either has
// succeeded or failed to produce. It will return the partition and the offset
// of the produced message, or an error if the message failed to produce.
func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	m.Produced <- justProducedAMsg
	return 0, 0, nil
}

// SendMessages produces a given set of messages, and returns only when all
// messages in the set have either succeeded or failed. Note that messages
// can succeed and fail individually; if some succeed and some fail,
// SendMessages will return an error.
func (m *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	for _, msg := range msgs {
		_, _, _ = m.SendMessage(msg)
	}
	return nil
}

// Close shuts down the producer and waits for any buffered messages to be
// flushed. You must call this function before a producer object passes out of
// scope, as it may otherwise leak memory. You must call this before calling
// Close on the underlying client.
func (m *MockSyncProducer) Close() error {
	close(m.Produced)
	return nil
}

func makeProducerContext() *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())

	_ = txctx.Envelope.SetChainID(big.NewInt(3))
	txctx.Set("produced", true)

	return txctx
}

func TestTxDisabling(t *testing.T) {

	mock := NewMockSyncProducer()
	handler := Producer(mock)

	txctx := makeProducerContext()
	expected := txctx.Get("produced").(bool)
	handler(txctx)
	// Depending on the case we detect if the message was produced or not
	var actual bool
	select {
	case <-time.After(time.Duration(1) * time.Second):
		actual = false
	case <-mock.Produced:
		actual = true
	}

	assert.Equalf(t, expected, actual, "Error tx filter failed")

}
