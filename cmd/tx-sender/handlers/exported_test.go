package handlers

import (
	"context"
	"testing"

	mockSarama "github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-sender/handlers/sender"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
)

// Init inialize handlers
func TestInit(t *testing.T) {
	producer := &mockSarama.SyncProducer{}
	broker.SetGlobalSyncProducer(producer)
	Init(context.Background())
	assert.NotNil(t, sender.GlobalHandler(), "Global store should have been set")
}
