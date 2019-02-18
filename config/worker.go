package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	workerInFlag     = "worker-in"
	workerInViperKey = "worker.in"
)

// WorkerInTopic register flag for kafka input topic
func WorkerInTopic(f *pflag.FlagSet, env string, defaultValue string) {
	desc := fmt.Sprintf(`Kafka topic to consume message from.
Environment variable: %q`, env)
	f.String(workerInFlag, defaultValue, desc)
	viper.BindPFlag(workerInViperKey, f.Lookup(workerInFlag))
	viper.BindEnv(workerInViperKey, env)
}

var (
	workerOutFlag     = "worker-out"
	workerOutViperKey = "worker.out"
)

// WorkerOutTopic register flag for kafka output topic
func WorkerOutTopic(f *pflag.FlagSet, env string, defaultValue string) {
	desc := fmt.Sprintf(`Kafka topic to send message to after processing.
Environment variable: %q`, env)
	f.String(workerOutFlag, defaultValue, desc)
	viper.BindPFlag(workerOutViperKey, f.Lookup(workerOutFlag))
	viper.BindEnv(workerOutViperKey, env)
}

var (
	txCrafterTopicEnv     = "KAFKA_TOPIC_TX_CRAFTER"
	txCrafterTopicDefault = "topic-tx-crafter"

	txNonceTopicEnv     = "KAFKA_TOPIC_TX_NONCE"
	txNonceTopicDefault = "topic-tx-nonce"

	txSignerTopicEnv     = "KAFKA_TOPIC_TX_SIGNER"
	txsignerTopicDefault = "topic-tx-signer"

	txSenderTopicEnvVar  = "KAFKA_TOPIC_TX_SENDER"
	txSenderTopicDefault = "topic-tx-sender"

	txDecoderTopicEnvVar  = "KAFKA_TOPIC_TX_DECODER"
	txDecoderTopicDefault = "topic-tx-decoder"

	txDecodedTopicEnvVar  = "KAFKA_TOPIC_TX_DECODED"
	txDecodedTopicDefault = "topic-tx-decoded"
)

// Kafka topic environemnt variables

// TxCrafterInTopic register flag for kafka input topic on tx crafter
func TxCrafterInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txCrafterTopicEnv, txCrafterTopicDefault)
}

// TxCrafterOutTopic register flag for kafka output topic on tx crafter
func TxCrafterOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txCrafterTopicEnv, txCrafterTopicDefault)
}

// TxNonceInTopic register flag for kafka input topic on tx nonce
func TxNonceInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txNonceTopicEnv, txNonceTopicDefault)
}

// TxNonceOutTopic register flag for kafka output topic on tx nonce
func TxNonceOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txNonceTopicEnv, txNonceTopicDefault)
}

// TxSignerInTopic register flag for kafka input topic on tx signer
func TxSignerInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txSignerTopicEnv, txsignerTopicDefault)
}

// TxSignerOutTopic register flag for kafka output topic on tx signer
func TxSignerOutTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txSignerTopicEnv, txsignerTopicDefault)
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

var (
	workerGroupFlag     = "worker-group"
	workerGroupViperKey = "worker.group"
)

// WorkerConsumerGroup register flag for kafka consumer group
func WorkerConsumerGroup(f *pflag.FlagSet, env string, defaultValue string) {
	desc := fmt.Sprintf(
		`Kafka consumer group. 
Environment variable: %q`, env)
	f.String(workerGroupFlag, defaultValue, desc)
	viper.BindPFlag(workerGroupViperKey, f.Lookup(workerGroupFlag))
	viper.BindEnv(workerGroupViperKey, env)
}

// Kafka Consumer group environment variables
var (
	crafterGroupEnv     = "KAFKA_CRAFTER_GROUP"
	crafterGroupDefault = "tx-crafter-group"

	nonceGroupEnv     = "KAFKA_NONCE_GROUP"
	nonceGroupDefault = "tx-nonce-group"

	signerGroupEnv     = "KAFKA_SIGNER_GROUP"
	signerGroupDefault = "tx-signer-group"

	senderGroupEnv     = "KAFKA_SENDER_GROUP"
	senderGroupDefault = "tx-sender-group"

	decoderGroupEnv     = "KAFKA_DECODER_GROUP"
	decoderGroupDefault = "tx-decoder-group"
)

// WorkerCrafterGroup register flag for kafka crafter group
func WorkerCrafterGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, crafterGroupEnv, crafterGroupDefault)
}

// WorkerNonceGroup register flag for kafka nonce group
func WorkerNonceGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, nonceGroupEnv, nonceGroupDefault)
}

// WorkerSignerGroup register flag for kafka signer group
func WorkerSignerGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, signerGroupEnv, signerGroupDefault)
}

// WorkerSenderGroup register flag for kafka sender group
func WorkerSenderGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, senderGroupEnv, senderGroupDefault)
}

// WorkerDecoderGroup register flag for kafka decoder group
func WorkerDecoderGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, decoderGroupEnv, decoderGroupDefault)
}

var (
	workerSlotsFlag     = "worker-slots"
	workerSlotsViperKey = "worker.slots"
	workerSlotsDefault  = uint(100)
	workerSlotsEnv      = "WORKER_SLOTS"
)

// WorkerSlots register flag for Kafka server addresses
func WorkerSlots(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of messages the worker can treat in parallel.
Environment variable: %q`, workerSlotsEnv)
	f.Uint(workerSlotsFlag, workerSlotsDefault, desc)
	viper.BindPFlag(workerSlotsViperKey, f.Lookup(workerSlotsFlag))
	viper.BindEnv(workerSlotsViperKey, workerSlotsEnv)
}
