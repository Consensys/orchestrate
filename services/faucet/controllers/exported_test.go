// +build unit

package controllers

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
	assert.NotNil(t, ctrl, "Control should have been set")
}
