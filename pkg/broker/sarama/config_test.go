// +build unit

package sarama

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestKafkaUrl(t *testing.T) {
	name := "kafka.url"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	KafkaProducerFlags(flgs)

	expected := []string{
		"localhost:9092",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Default")

	_ = os.Setenv("KAFKA_URL", "localhost:9192")
	expected = []string{
		"localhost:9192",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "From Environment Variable")
	_ = os.Unsetenv("KAFKA_URL")

	args := []string{
		"--kafka-url=127.0.0.1:9091",
		"--kafka-url=127.0.0.2:9091",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = []string{
		"127.0.0.1:9091",
		"127.0.0.2:9091",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "From flag")
}

func TestKafkaConsumerMaxWaitTime(t *testing.T) {

	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	kafkaConsumerMaxWaitTime(f)

	assert.Equal(t, kafkaConsumerMaxWaitTimeDefault, viper.GetDuration(kafkaConsumerMaxWaitTimeViperKey), "Default")
}

func TestKafkaConsumerMaxProcessingTime(t *testing.T) {

	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	kafkaConsumerMaxProcessingTime(f)

	assert.Equal(t, kafkaConsumerMaxProcessingTimeDefault, viper.GetDuration(kafkaConsumerMaxProcessingTimeViperKey), "Default")
}

func TestKafkaConsumerGroupHeartbeatInterval(t *testing.T) {

	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	kafkaConsumerGroupHeartbeatInterval(f)

	assert.Equal(t, kafkaConsumerGroupHeartbeatIntervalDefault, viper.GetDuration(kafkaConsumerGroupHeartbeatIntervalViperKey), "Default")
}

func TestKafkaConsumerGroupRebalanceTimeout(t *testing.T) {

	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	kafkaConsumerGroupRebalanceTimeout(f)

	assert.Equal(t, kafkaConsumerGroupRebalanceTimeoutDefault, viper.GetDuration(kafkaConsumerGroupRebalanceTimeoutViperKey), "Default")
}

func TestTopics(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	KafkaTopicTxSender(flgs)
	assert.Equal(t, "topic-tx-sender", viper.GetString("topic.tx.sender"), "From default")

	KafkaTopicTxRecover(flgs)
	assert.Equal(t, "topic-tx-recover", viper.GetString("topic.tx.recover"), "From default")

	KafkaTopicTxDecoded(flgs)
	assert.Equal(t, "topic-tx-decoded", viper.GetString("topic.tx.decoded"), "From default")
}

func TestConsumerGroupName(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	KafkaConsumerFlags(flgs)
	assert.Equal(t, "group-sender", viper.GetString("kafka.consumer.group.name"), "From default")
}

func TestInitKafkaSASLFlags(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	KafkaProducerFlags(flgs)
	assert.Equal(t, false, viper.GetBool("kafka.sasl.enable"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.sasl.mechanism"), "From default")
	assert.Equal(t, true, viper.GetBool("kafka.sasl.handshake"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.sasl.user"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.sasl.password"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.sasl.scramauthzid"), "From default")
}

func TestInitKafkaTLSFlags(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	KafkaConsumerFlags(flgs)
	assert.Equal(t, false, viper.GetBool("kafka.tls.enabled"), "From default")
	assert.Equal(t, false, viper.GetBool("kafka.tls.insecure.skip.verify"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.tls.client.cert.file"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.tls.client.key.file"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.tls.ca.cert.file"), "From default")
}

func TestKafkaConsumerFlags(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	KafkaConsumerFlags(flgs)

	TestKafkaUrl(t)
	TestInitKafkaSASLFlags(t)
	TestInitKafkaTLSFlags(t)
	TestKafkaConsumerMaxWaitTime(t)
	TestKafkaConsumerGroupRebalanceTimeout(t)
	TestKafkaConsumerGroupHeartbeatInterval(t)
	TestKafkaConsumerMaxProcessingTime(t)
}
