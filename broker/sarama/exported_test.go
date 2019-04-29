package sarama

// TODO: remove comments once we can connect to a Kafka broker in CI

// import (
// 	"context"
// 	"testing"

// 	"github.com/spf13/viper"
// 	"github.com/stretchr/testify/assert"
// )

// func TestInitClient(t *testing.T) {
// 	viper.Set(kafkaAddressViperKey, []string{
// 		"localhost:9092",
// 	})
// 	InitClient(context.Background())
// 	assert.NotNil(t, client, "Client should have been set")
// }

// func TestInitSyncProducer(t *testing.T) {
// 	InitSyncProducer(context.Background())
// 	assert.NotNil(t, producer, "Producer should have been set")
// }

// func TestInitConsumerGroup(t *testing.T) {
// 	InitConsumerGroup(context.Background())
// 	assert.NotNil(t, group, "ConsumerGroup should have been set")
// }
