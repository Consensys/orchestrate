package sarama

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestKafkaAddresses(t *testing.T) {
	name := "kafka.addresses"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	KafkaAddresses(flgs)

	expected := []string{
		"localhost:9092",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Default")

	os.Setenv("KAFKA_ADDRESS", "localhost:9192")
	expected = []string{
		"localhost:9192",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "From Environment Variable")
	os.Unsetenv("KAFKA_ADDRESS")

	args := []string{
		"--kafka-address=127.0.0.1:9091",
		"--kafka-address=127.0.0.2:9091",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = []string{
		"127.0.0.1:9091",
		"127.0.0.2:9091",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "From flag")
}

func TestTopics(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	KafkaTopicTxCrafter(flgs)
	assert.Equal(t, "topic-tx-crafter", viper.GetString("kafka.topic.crafter"), "From default")

	KafkaTopicTxNonce(flgs)
	assert.Equal(t, "topic-tx-nonce", viper.GetString("kafka.topic.nonce"), "From default")

	KafkaTopicTxSigner(flgs)
	assert.Equal(t, "topic-tx-signer", viper.GetString("kafka.topic.signer"), "From default")

	KafkaTopicTxSender(flgs)
	assert.Equal(t, "topic-tx-sender", viper.GetString("kafka.topic.sender"), "From default")

	KafkaTopicTxDecoder(flgs)
	assert.Equal(t, "topic-tx-decoder", viper.GetString("kafka.topic.decoder"), "From default")

	KafkaTopicTxRecover(flgs)
	assert.Equal(t, "topic-tx-recover", viper.GetString("kafka.topic.recover"), "From default")

	KafkaTopicTxDecoded(flgs)
	assert.Equal(t, "topic-tx-decoded", viper.GetString("kafka.topic.decoded"), "From default")

	KafkaTopicWalletGenerator(flgs)
	assert.Equal(t, "topic-wallet-generator", viper.GetString("kafka.topic.wallet.generator"), "From default")

	KafkaTopicWalletGenerated(flgs)
	assert.Equal(t, "topic-wallet-generated", viper.GetString("kafka.topic.wallet.generated"), "From default")
}

func TestConsumerGroup(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	CrafterGroup(flgs)
	assert.Equal(t, "group-crafter", viper.GetString("kafka.group.crafter"), "From default")

	NonceGroup(flgs)
	assert.Equal(t, "group-nonce", viper.GetString("kafka.group.nonce"), "From default")

	SignerGroup(flgs)
	assert.Equal(t, "group-signer", viper.GetString("kafka.group.signer"), "From default")

	SenderGroup(flgs)
	assert.Equal(t, "group-sender", viper.GetString("kafka.group.sender"), "From default")

	DecoderGroup(flgs)
	assert.Equal(t, "group-decoder", viper.GetString("kafka.group.decoder"), "From default")

	BridgeGroup(flgs)
	assert.Equal(t, "group-bridge", viper.GetString("kafka.group.bridge"), "From default")
}
