package kafka

import (
	"context"
	"testing"

	"github.com/Shopify/sarama/mocks"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	producer := &mocks.SyncProducer{}
	broker.SetGlobalSyncProducer(producer)

	Init(context.Background())
	assert.NotNil(t, GlobalHook(), "Global should have been set")

	SetGlobalHook(nil)
	assert.Nil(t, GlobalHook(), "Global should be reset to nil")
}
