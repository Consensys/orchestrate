// +build unit

package txsigner

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

func TestPrepareMsgSigner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	m := mock.NewMockMsg(mockCtrl)
	m.EXPECT().Key().Return([]byte(`test`)).AnyTimes()
	m.EXPECT().Entrypoint().Return("topic-tx-signer").AnyTimes()

	// No error
	txctx := engine.NewTxContext()

	txctx.In = m
	msg := &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-sender", msg.Topic, "If no error out topic should be sender")

	// Classic error
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	msg = &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-recover", msg.Topic, "If error out topic should be recovery")
	
	// Skip child job error
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	_ = txctx.Envelope.SetContextLabelsValue(tx.ParentJobUUIDLabel, "parentJobUUID")
	msg = &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Empty(t, msg.Topic, "If error on children job don't send message")
}
