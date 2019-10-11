package sarama

import (
	"context"
	"testing"

	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
)

func TestInit(t *testing.T) {
	broker.SetGlobalSyncProducer(mocks.NewSyncProducer(t, nil))
	Init(context.Background())
	assert.NotNil(t, fct, "Faucet should have been set")

	var f *Faucet
	SetGlobalFaucet(f)
	assert.Nil(t, GlobalFaucet(), "Global should be reset to nil")
}
