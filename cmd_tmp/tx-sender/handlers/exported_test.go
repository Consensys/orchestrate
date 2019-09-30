package handlers

import (
	"context"
	"testing"

	mockSarama "github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers/sender"
)

// Init inialize handlers
func TestInit(t *testing.T) {
	producer := &mockSarama.SyncProducer{}
	broker.SetGlobalSyncProducer(producer)
	Init(context.Background())
	assert.NotNil(t, sender.GlobalHandler(), "Global store should have been set")
}
