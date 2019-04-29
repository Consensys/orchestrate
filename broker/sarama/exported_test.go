package sarama

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TODO: implement Unittests but it requires to have a Kafka we can connect to on CI/CD

func TestInitClient(t *testing.T) {
	viper.Set(kafkaAddressViperKey, []string{
		"localhost:9092",
	})
	InitClient(context.Background())
	assert.NotNil(t, client, "Client should have been set")
}

func TestInitSyncProducer(t *testing.T) {
	InitSyncProducer(context.Background())
	assert.NotNil(t, producer, "Producer should have been set")
}

func TestInitConsumerGroup(t *testing.T) {
	InitConsumerGroup(context.Background())
	assert.NotNil(t, group, "ConsumerGroup should have been set")
}
