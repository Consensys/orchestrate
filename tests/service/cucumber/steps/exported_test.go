package steps

import (
	"context"
	"testing"

	mockSarama "github.com/Shopify/sarama/mocks"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
)

func TestInit(t *testing.T) {
	producer := &mockSarama.SyncProducer{}
	broker.SetGlobalSyncProducer(producer)
	Init(context.Background())
}
