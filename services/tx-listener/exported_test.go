// +build unit

package txlistener

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
	assert.NotNil(t, GlobalListener(), "Global should have been set")

	SetGlobalListener(nil)
	assert.Nil(t, GlobalListener(), "Global should be reset to nil")
}
