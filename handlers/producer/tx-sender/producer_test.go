// +build unit

package txsender

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
)

func TestPrepareMsg(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	m := mock.NewMockMsg(mockCtrl)
	m.EXPECT().Key().Return([]byte(`test`)).AnyTimes()

	// No error
	txctx := engine.NewTxContext()
	txctx.In = m
	msg := &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "", msg.Topic, "If no error there should be no out topic")

	// Error
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	msg = &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-recover", msg.Topic, "If error out topic should be recovery")

	txctx = engine.NewTxContext()
	txctx.In = m
	txctx.Set("invalid.nonce", true)
	msg = &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-crafter", msg.Topic, "If invalid nonce out topic should be tx-crafter")

	txctx = engine.NewTxContext()
	txctx.In = m
	txctx.Set("invalid.nonce", true)
	_ = txctx.Error(errors.ConnectionError("nonce too low"))
	msg = &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-recover", msg.Topic, "If invalid nonce and error topic should be tx-recover")
	
	// Skip child job error
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	_ = txctx.Envelope.SetContextLabelsValue(tx.ParentJobUUIDLabel, "parentJobUUID")
	msg = &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Empty(t, msg.Topic, "If error on children job don't send message")
}
