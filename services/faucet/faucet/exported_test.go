// +build unit

package faucet

import (
	"context"
	"testing"

	"github.com/Shopify/sarama/mocks"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	broker.SetGlobalSyncProducer(mocks.NewSyncProducer(t, nil))
	viper.Set("faucet.type", "sarama")
	Init(context.Background())
	assert.NotNil(t, fct, "Faucet should have been set")

	var f Faucet
	SetGlobalFaucet(f)
	assert.Nil(t, GlobalFaucet(), "Global should be reset to nil")
}
