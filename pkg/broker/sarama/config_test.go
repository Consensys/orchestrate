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
	KafkaURL(flgs)

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
	KafkaConsumerMaxWaitTime(f)

	assert.Equal(t, kafkaConsumerMaxWaitTimeDefault, viper.GetInt(kafkaConsumerMaxWaitTimeViperKey), "Default")
}

func TestTopics(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	KafkaTopicTxCrafter(flgs)
	assert.Equal(t, "topic-tx-crafter", viper.GetString("topic.tx.crafter"), "From default")

	KafkaTopicTxSigner(flgs)
	assert.Equal(t, "topic-tx-signer", viper.GetString("topic.tx.signer"), "From default")

	KafkaTopicTxSender(flgs)
	assert.Equal(t, "topic-tx-sender", viper.GetString("topic.tx.sender"), "From default")

	KafkaTopicTxRecover(flgs)
	assert.Equal(t, "topic-tx-recover", viper.GetString("topic.tx.recover"), "From default")

	KafkaTopicTxDecoded(flgs)
	assert.Equal(t, "topic-tx-decoded", viper.GetString("topic.tx.decoded"), "From default")

	KafkaTopicAccountGenerator(flgs)
	assert.Equal(t, "topic-account-generator", viper.GetString("topic.account.generator"), "From default")

	KafkaTopicAccountGenerated(flgs)
	assert.Equal(t, "topic-account-generated", viper.GetString("topic.account.generated"), "From default")
}

func TestConsumerGroup(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	CrafterGroup(flgs)
	assert.Equal(t, "group-crafter", viper.GetString("kafka.group.crafter"), "From default")

	SignerGroup(flgs)
	assert.Equal(t, "group-signer", viper.GetString("kafka.group.signer"), "From default")

	SenderGroup(flgs)
	assert.Equal(t, "group-sender", viper.GetString("kafka.group.sender"), "From default")

	DecoderGroup(flgs)
	assert.Equal(t, "group-decoder", viper.GetString("kafka.group.decoder"), "From default")
}

func TestInitKafkaSASLFlags(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	InitKafkaSASLFlags(flgs)
	assert.Equal(t, false, viper.GetBool("kafka.sasl.enable"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.sasl.mechanism"), "From default")
	assert.Equal(t, true, viper.GetBool("kafka.sasl.handshake"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.sasl.user"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.sasl.password"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.sasl.scramauthzid"), "From default")
}

func TestInitKafkaTLSFlags(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	InitKafkaTLSFlags(flgs)
	assert.Equal(t, false, viper.GetBool("kafka.tls.enabled"), "From default")
	assert.Equal(t, false, viper.GetBool("kafka.tls.insecure.skip.verify"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.tls.client.cert.file"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.tls.client.key.file"), "From default")
	assert.Equal(t, "", viper.GetString("kafka.tls.ca.cert.file"), "From default")
}

func TestInitKafkaFlags(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	InitKafkaFlags(flgs)

	TestKafkaUrl(t)
	TestInitKafkaSASLFlags(t)
	TestInitKafkaTLSFlags(t)
	TestKafkaConsumerMaxWaitTime(t)
}
