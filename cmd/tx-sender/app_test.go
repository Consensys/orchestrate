// +build unit

package txsender

import (
	"context"
	"testing"

	mockSarama "github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/sender"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
)

func TestInitHandlers(t *testing.T) {
	producer := &mockSarama.SyncProducer{}
	broker.SetGlobalSyncProducer(producer)
	initHandlers(context.Background())
	assert.NotNil(t, sender.GlobalHandler(), "Global store should have been set")
}
