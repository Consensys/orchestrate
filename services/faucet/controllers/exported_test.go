// +build unit

package controllers

import (
	"context"
	"testing"

	"github.com/Shopify/sarama/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
)

func TestInit(t *testing.T) {
	broker.SetGlobalSyncProducer(mocks.NewSyncProducer(t, nil))
	viper.Set("faucet.type", "sarama")
	Init(context.Background())
	assert.NotNil(t, ctrl, "Control should have been set")
}
