package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(txCrafterViperKey, txCrafterTopicDefault)
	viper.BindEnv(txCrafterViperKey, txCrafterTopicEnv)
	viper.SetDefault(txNonceViperKey, txNonceTopicDefault)
	viper.BindEnv(txNonceViperKey, txNonceTopicEnv)
	viper.SetDefault(generateWalletViperKey, generateWalletDefault)
	viper.BindEnv(generateWalletViperKey, generateWalletTopicEnv)
	viper.SetDefault(txSignerViperKey, txSignerTopicDefault)
	viper.BindEnv(txSignerViperKey, txSignerTopicEnv)
	viper.SetDefault(txSenderViperKey, txSenderTopicDefault)
	viper.BindEnv(txSenderViperKey, txSenderTopicEnv)
	viper.SetDefault(txDecoderViperKey, txDecoderTopicDefault)
	viper.BindEnv(txDecoderViperKey, txDecoderTopicEnv)
	viper.SetDefault(txRecoverViperKey, txRecoverTopicDefault)
	viper.BindEnv(txRecoverViperKey, txRecoverTopicEnv)
}

var (
	txCrafterFlag         = "topic-crafter"
	txCrafterViperKey     = "kafka.topic.crafter"
	txCrafterTopicEnv     = "KAFKA_TOPIC_TX_CRAFTER"
	txCrafterTopicDefault = "topic-tx-crafter"

	txNonceFlag         = "topic-nonce"
	txNonceViperKey     = "kafka.topic.nonce"
	txNonceTopicEnv     = "KAFKA_TOPIC_TX_NONCE"
	txNonceTopicDefault = "topic-tx-nonce"

	generateWalletFlag     = "topic-wallet"
	generateWalletViperKey = "kafka.topic.wallet.generator"
	generateWalletTopicEnv = "KAFKA_TOPIC_WALLET_GENERATOR"
	generateWalletDefault  = "topic-wallet-generator"

	txSignerFlag         = "topic-signer"
	txSignerViperKey     = "kafka.topic.signer"
	txSignerTopicEnv     = "KAFKA_TOPIC_TX_SIGNER"
	txSignerTopicDefault = "topic-tx-signer"

	txSenderFlag         = "topic-sender"
	txSenderViperKey     = "kafka.topic.sender"
	txSenderTopicEnv     = "KAFKA_TOPIC_TX_SENDER"
	txSenderTopicDefault = "topic-tx-sender"

	txDecoderFlag         = "topic-decoder"
	txDecoderViperKey     = "kafka.topic.decoder"
	txDecoderTopicEnv     = "KAFKA_TOPIC_TX_DECODER"
	txDecoderTopicDefault = "topic-tx-decoder"

	txRecoverFlag         = "topic-recover"
	txRecoverViperKey     = "kafka.topic.recover"
	txRecoverTopicEnv     = "KAFKA_TOPIC_TX_RECOVER"
	txRecoverTopicDefault = "topic-tx-recover"
)

// TODO: implement test for all flags

// KafkaTopicTxCrafter register flag for Kafka topic
func KafkaTopicTxCrafter(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages waiting to have transaction payload crafted
Environment variable: %q`, txCrafterTopicEnv)
	f.String(txCrafterFlag, txCrafterTopicDefault, desc)
	viper.BindPFlag(txCrafterViperKey, f.Lookup(txCrafterFlag))
}

// KafkaTopicTxNonce register flag for Kafka topic
func KafkaTopicTxNonce(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages waiting to have transaction nonce set
Environment variable: %q`, txNonceViperKey)
	f.String(txNonceFlag, txNonceTopicDefault, desc)
	viper.BindPFlag(txNonceViperKey, f.Lookup(txNonceFlag))
}

// KafkaTopicTxSigner register flag for Kafka topic
func KafkaTopicTxSigner(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages waiting to have transaction signed
Environment variable: %q`, txSignerViperKey)
	f.String(txSignerFlag, txSignerTopicDefault, desc)
	viper.BindPFlag(txSignerViperKey, f.Lookup(txSignerFlag))
}

// KafkaTopicWalletGenerator register flag for Kafka topic
func KafkaTopicWalletGenerator(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages waiting to generate a new wallet
Environment variable: %q`, generateWalletViperKey)
	f.String(generateWalletFlag, generateWalletDefault, desc)
	viper.BindPFlag(generateWalletViperKey, f.Lookup(generateWalletFlag))
}

// KafkaTopicTxSender register flag for Kafka topic
func KafkaTopicTxSender(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages waiting to have transaction sent
Environment variable: %q`, txSenderViperKey)
	f.String(txSenderFlag, txSenderTopicDefault, desc)
	viper.BindPFlag(txSenderViperKey, f.Lookup(txSenderFlag))
}

// KafkaTopicTxDecoder register flag for Kafka topic
func KafkaTopicTxDecoder(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages waiting to have receipt decoded
Environment variable: %q`, txDecoderViperKey)
	f.String(txDecoderFlag, txDecoderTopicDefault, desc)
	viper.BindPFlag(txDecoderViperKey, f.Lookup(txDecoderFlag))
}

// KafkaTopicTxRecover register flag for Kafka topic
func KafkaTopicTxRecover(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages waiting to have transaction recovered
Environment variable: %q`, txRecoverViperKey)
	f.String(txRecoverFlag, txRecoverTopicDefault, desc)
	viper.BindPFlag(txRecoverViperKey, f.Lookup(txRecoverFlag))
}

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

	bridgeGroupEnv     = "KAFKA_BRIDGE_GROUP"
	bridgeGroupDefault = "tx-bridge-group"
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

// WorkerBridgeGroup register flag for kafka decoder group
func WorkerBridgeGroup(f *pflag.FlagSet) {
	WorkerConsumerGroup(f, decoderGroupEnv, decoderGroupDefault)
}

var (
	workerInFlag     = "worker-in"
	workerInViperKey = "engine.in"
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
	workerOutViperKey = "engine.out"
)

// WorkerOutTopic register flag for kafka output topic
func WorkerOutTopic(f *pflag.FlagSet, env string, defaultValue string) {
	desc := fmt.Sprintf(`Kafka topic to send message to after processing.
Environment variable: %q`, env)
	f.String(workerOutFlag, defaultValue, desc)
	viper.BindPFlag(workerOutViperKey, f.Lookup(workerOutFlag))
	viper.BindEnv(workerOutViperKey, env)
}

// Kafka topic environment variables

// TxCrafterInTopic register flag for kafka input topic on tx crafter
func TxCrafterInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txCrafterTopicEnv, txCrafterTopicDefault)
}

// TxCrafterOutTopic register flag for kafka output topic on tx crafter
func TxCrafterOutTopic(f *pflag.FlagSet) {
	WorkerOutTopic(f, txCrafterTopicEnv, txCrafterTopicDefault)
}

// TxNonceInTopic register flag for kafka input topic on tx nonce
func TxNonceInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txNonceTopicEnv, txNonceTopicDefault)
}

// TxNonceOutTopic register flag for kafka output topic on tx nonce
func TxNonceOutTopic(f *pflag.FlagSet) {
	WorkerOutTopic(f, txNonceTopicEnv, txNonceTopicDefault)
}

// TxSignerInTopic register flag for kafka input topic on tx signer
func TxSignerInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txSignerTopicEnv, txSignerTopicDefault)
}

// TxSignerOutTopic register flag for kafka output topic on tx signer
func TxSignerOutTopic(f *pflag.FlagSet) {
	WorkerOutTopic(f, txSignerTopicEnv, txSignerTopicDefault)
}

// TxSenderInTopic register flag for kafka input topic on tx sender
func TxSenderInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txSenderTopicEnv, txSenderTopicDefault)
}

// TxSenderOutTopic register flag for kafka output topic on tx sender
func TxSenderOutTopic(f *pflag.FlagSet) {
	WorkerOutTopic(f, txSenderTopicEnv, txSenderTopicDefault)
}

// TxDecoderInTopic register flag for kafka input topic on tx decoder
func TxDecoderInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txDecoderTopicEnv, txDecoderTopicDefault)
}

// TxDecoderOutTopic register flag for kafka output topic on tx decoder
func TxDecoderOutTopic(f *pflag.FlagSet) {
	WorkerOutTopic(f, txDecoderTopicEnv, txDecoderTopicDefault)
}

// TxDecodedInTopic register flag for kafka input topic on tx decoded
func TxDecodedInTopic(f *pflag.FlagSet) {
	WorkerInTopic(f, txRecoverTopicEnv, txRecoverTopicDefault)
}

// TxDecodedOutTopic register flag for kafka output topic on tx decoded
func TxDecodedOutTopic(f *pflag.FlagSet) {
	WorkerOutTopic(f, txRecoverTopicEnv, txRecoverTopicDefault)
}

var (
	workerGroupFlag     = "worker-group"
	workerGroupViperKey = "engine.group"
)
