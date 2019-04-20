package sarama

import (
	"context"
	"testing"

	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
)

func TestInit(t *testing.T) {
	broker.SetSyncProducer(mocks.NewSyncProducer(t, nil))
	Init(context.Background())
	assert.NotNil(t, fct, "Faucet should have been set")
}
