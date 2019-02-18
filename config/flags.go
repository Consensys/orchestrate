package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// LogLevel register flag for LogLevel
func LogLevel(f *pflag.FlagSet) {
	flagName := "log-level"
	viperName := "log.level"
	envName := "LOG_LEVEL"
	defaultValue := "debug"
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}, envName)
	f.String(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}

// LogFormat register flag for Log Format
func LogFormat(f *pflag.FlagSet) {
	flagName := "log-format"
	viperName := "log.format"
	envName := "LOG_FORMAT"
	defaultValue := "text"
	desc := fmt.Sprintf(`Log formatter (one of %q).
Environment variable: %q`, []string{"text", "json"}, envName)
	f.String(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}

// EthClientURLs register flag for Ethereum client URLs
func EthClientURLs(f *pflag.FlagSet) {
	// TODO: add configure from Env Variable
	// It requires a PR on viper https://github.com/spf13/viper/blob/master/viper.go#L1008
	// So it parses automatically
	flagName := "eth-client"
	viperName := "eth.clients"
	envName := "ETH_CLIENT_URL"
	defaultValue := []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	}
	desc := fmt.Sprintf(`Ethereum client URLs.
	Environment variable: %q`, envName)
	f.StringSlice(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}

// KafkaAddresses register flag for Kafka server addresses
func KafkaAddresses(f *pflag.FlagSet) {
	flagName := "kafka-address"
	viperName := "kafka.addresses"
	envName := "KAFKA_ADDRESS"
	defaultValue := []string{
		"localhost:9092",
	}
	desc := fmt.Sprintf(`Address of Kafka server to connect to.
Environment variable: %q (but list is not yet supported)`, envName)
	f.StringSlice(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}

// WorkerInTopic register flag for kafka input topic
func WorkerInTopic(f *pflag.FlagSet, envName string, defaultValue string) {
	flagName := "worker-in"
	viperName := "worker.in"
	desc := fmt.Sprintf(`Kafka topic to consume message from.
Environment variable: %q`, envName)
	f.String(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}

// WorkerOutTopic register flag for kafka output topic
func WorkerOutTopic(f *pflag.FlagSet, envName string, defaultValue string) {
	flagName := "worker-out"
	viperName := "worker.out"
	desc := fmt.Sprintf(`Kafka topic to send message to after processing.
Environment variable: %q`, envName)
	f.String(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}

// Kafka topic environemnt variables
var (
	txCrafterTopicEnvVar, txCrafterTopicDefault = "KAFKA_TOPIC_TX_CRAFTER", "topic-tx-crafter"
	txNonceTopicEnvVar, txNonceTopicDefault     = "KAFKA_TOPIC_TX_NONCE", "topic-tx-nonce"
	txSignerTopicEnvVar, txsignerTopicDefault   = "KAFKA_TOPIC_TX_SIGNER", "topic-tx-signer"
	txSenderTopicEnvVar, txSenderTopicDefault   = "KAFKA_TOPIC_TX_SENDER", "topic-tx-sender"
	txDecoderTopicEnvVar, txDecoderTopicDefault = "KAFKA_TOPIC_TX_DECODER", "topic-tx-decoder"
	txDecodedTopicEnvVar, txDecodedTopicDefault = "KAFKA_TOPIC_TX_DECODED", "topic-tx-decoded"
)

// TxCrafterInTopic register flag for kafka input topic on tx crafter
func TxCrafterInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txCrafterTopicEnvVar, txCrafterTopicDefault)
}

// TxCrafterOutTopic register flag for kafka output topic on tx crafter
func TxCrafterOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txCrafterTopicEnvVar, txCrafterTopicDefault)
}

// TxNonceInTopic register flag for kafka input topic on tx nonce
func TxNonceInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txNonceTopicEnvVar, txNonceTopicDefault)
}

// TxNonceOutTopic register flag for kafka output topic on tx nonce
func TxNonceOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txNonceTopicEnvVar, txNonceTopicDefault)
}

// TxSignerInTopic register flag for kafka input topic on tx signer
func TxSignerInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txSignerTopicEnvVar, txsignerTopicDefault)
}

// TxSignerOutTopic register flag for kafka output topic on tx signer
func TxSignerOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txSignerTopicEnvVar, txsignerTopicDefault)
}

// TxSenderInTopic register flag for kafka input topic on tx sender
func TxSenderInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txSenderTopicEnvVar, txSenderTopicDefault)
}

// TxSenderOutTopic register flag for kafka output topic on tx sender
func TxSenderOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txSenderTopicEnvVar, txSenderTopicDefault)
}

// TxDecoderInTopic register flag for kafka input topic on tx decoder
func TxDecoderInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txDecoderTopicEnvVar, txDecoderTopicDefault)
}

// TxDecoderOutTopic register flag for kafka output topic on tx decoder
func TxDecoderOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txDecoderTopicEnvVar, txDecoderTopicDefault)
}

// TxDecodedInTopic register flag for kafka input topic on tx decoded
func TxDecodedInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txDecodedTopicEnvVar, txDecodedTopicDefault)
}

// TxDecodedOutTopic register flag for kafka output topic on tx decoded
func TxDecodedOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txDecodedTopicEnvVar, txDecodedTopicDefault)
}

// WorkerConsumerGroup register flag for kafka consumer group
func WorkerConsumerGroup(f *pflag.FlagSet, envName string, defaultValue string) {
	flagName := "worker-group"
	viperName := "worker.group"
	desc := fmt.Sprintf(
		`Kafka consumer group. 
Environment variable: %q`, envName)
	f.String(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}

// Kafka Consumer group environment variables
var (
	crafterGroupEnvVar, crafterGroupDefault = "KAFKA_CRAFTER_GROUP", "tx-crafter-group"
	nonceGroupEnvVar, nonceGroupDefault     = "KAFKA_NONCE_GROUP", "tx-nonce-group"
	signerGroupEnvVar, signerGroupDefault   = "KAFKA_SIGNER_GROUP", "tx-signer-group"
	senderGroupEnvVar, senderGroupDefault   = "KAFKA_SENDER_GROUP", "tx-sender-group"
	decoderGroupEnvVar, decoderGroupDefault = "KAFKA_DECODER_GROUP", "tx-decoder-group"
)

// WorkerCrafterGroup register flag for kafka crafter group
func WorkerCrafterGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, crafterGroupEnvVar, crafterGroupDefault)
}

// WorkerNonceGroup register flag for kafka nonce group
func WorkerNonceGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, nonceGroupEnvVar, nonceGroupDefault)
}

// WorkerSignerGroup register flag for kafka signer group
func WorkerSignerGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, signerGroupEnvVar, signerGroupDefault)
}

// WorkerSenderGroup register flag for kafka sender group
func WorkerSenderGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, senderGroupEnvVar, senderGroupDefault)
}

// WorkerDecoderGroup register flag for kafka decoder group
func WorkerDecoderGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, decoderGroupEnvVar, decoderGroupDefault)
}

// WorkerSlots register flag for Kafka server addresses
func WorkerSlots(f *pflag.FlagSet) {
	flagName := "worker-slots"
	viperName := "worker.slots"
	envName := "WORKER_SLOTS"
	defaultValue := uint(100)
	desc := fmt.Sprintf(`Maximum number of messages the worker can treat in parallel.
Environment variable: %q`, envName)
	f.Uint(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}

// RedisAddress register a flag for Redis server address
func RedisAddress(f *pflag.FlagSet) {
	flagName := "redis-address"
	viperName := "redis.address"
	envName := "REDIS_ADDRESS"
	defaultValue := "localhost:6379"
	desc := fmt.Sprintf(`Address of Redis server to connect to.
Environment variable: %q`, envName)
	f.String(flagName, defaultValue, desc)
	viper.BindPFlag(viperName, f.Lookup(flagName))
	viper.BindEnv(viperName, envName)
}
