// +build unit

package txcrafter

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
	assert.Equal(t, "topic-tx-signer", msg.Topic, "If no error out topic should be nonce")
	
	// Faucet warning
	msg = &sarama.ProducerMessage{}
	_ = txctx.Error(errors.FaucetWarning("no credit"))
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-signer", msg.Topic, "If Faucet warning out topic should be nonce")
	
	// Invalid Auth error
	msg = &sarama.ProducerMessage{}
	_ = txctx.Error(errors.InvalidAuthenticationError("Connection error"))
	_ = PrepareMsg(txctx, msg)
	assert.Equal(t, "topic-tx-recover", msg.Topic, "If error out topic should be recovery")
	
	// Communication error
	msg = &sarama.ProducerMessage{}
	_ = txctx.Error(errors.ConnectionError("Connection error"))
	_ = PrepareMsg(txctx, msg)
	assert.Empty(t, msg.Topic, "If error on children job don't send message")
	assert.NotNil(t, txctx.HasRetryMsgErr())
	
	// Skip child job error
	msg = &sarama.ProducerMessage{}
	_ = txctx.Error(errors.InvalidAuthenticationError("Connection error"))
	_ = txctx.Envelope.SetContextLabelsValue(tx.ParentJobUUIDLabel, "parentJobUUID")
	_ = PrepareMsg(txctx, msg)
	assert.Empty(t, msg.Topic, "If error on children job don't send message")
}
