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
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-recover", msg.Topic, "If error out topic should be recovery")
}

func TestPrepareMsgGenerateAccount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	m := mock.NewMockMsg(mockCtrl)
	m.EXPECT().Key().Return([]byte(`test`)).AnyTimes()
	m.EXPECT().Entrypoint().Return("topic-account-generator").AnyTimes()

	// No error
	txctx := engine.NewTxContext()

	txctx.In = m
	msg := &sarama.ProducerMessage{}
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-account-generated", msg.Topic, "If no error out topic should be account-generated")

	// Classic error
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-account-generated", msg.Topic, "If error out topic should be recovery")
}
