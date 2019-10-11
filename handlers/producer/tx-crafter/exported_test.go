package txcrafter

import (
	"context"
	"testing"

	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
)

func TestInit(t *testing.T) {
	producer := &mocks.SyncProducer{}
	broker.SetGlobalSyncProducer(producer)

	Init(context.Background())
	assert.NotNil(t, handler, "Global handler should have been set")

	var h engine.HandlerFunc
	SetGlobalHandler(h)
	assert.Nil(t, GlobalHandler(), "Global should be reset to nil")
}
