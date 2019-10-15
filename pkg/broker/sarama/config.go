package sarama

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {

	// Kafka general parameters
	viper.SetDefault(kafkaAddressViperKey, kafkaAddressDefault)
	_ = viper.BindEnv(kafkaAddressViperKey, kafkaAddressEnv)
	viper.SetDefault(kafkaGroupViperKey, kafkaGroupDefault)
	_ = viper.BindEnv(kafkaGroupViperKey, kafkaGroupEnv)

	// Kafka topics for the tx workflow
	viper.SetDefault(txCrafterViperKey, txCrafterTopicDefault)
	_ = viper.BindEnv(txCrafterViperKey, txCrafterTopicEnv)
	viper.SetDefault(txNonceViperKey, txNonceTopicDefault)
	_ = viper.BindEnv(txNonceViperKey, txNonceTopicEnv)
	viper.SetDefault(txSignerViperKey, txSignerTopicDefault)
	_ = viper.BindEnv(txSignerViperKey, txSignerTopicEnv)
	viper.SetDefault(txSenderViperKey, txSenderTopicDefault)
	_ = viper.BindEnv(txSenderViperKey, txSenderTopicEnv)
	viper.SetDefault(txDecoderViperKey, txDecoderTopicDefault)
	_ = viper.BindEnv(txDecoderViperKey, txDecoderTopicEnv)
	viper.SetDefault(txDecodedViperKey, txDecodedTopicDefault)
	_ = viper.BindEnv(txDecodedViperKey, txDecodedTopicEnv)
	viper.SetDefault(txRecoverViperKey, txRecoverTopicDefault)
	_ = viper.BindEnv(txRecoverViperKey, txRecoverTopicEnv)

	// Kafka topics for the wallet generation workflow
	viper.SetDefault(walletGeneratorViperKey, walletGeneratorDefault)
	_ = viper.BindEnv(walletGeneratorViperKey, walletGeneratorTopicEnv)
	viper.SetDefault(walletGeneratedViperKey, walletGeneratedDefault)
	_ = viper.BindEnv(walletGeneratedViperKey, walletGeneratedTopicEnv)

	// Kafka consumer groups for tx workflow
	viper.SetDefault(crafterGroupViperKey, crafterGroupDefault)
	_ = viper.BindEnv(crafterGroupViperKey, crafterGroupEnv)
	viper.SetDefault(nonceGroupViperKey, nonceGroupDefault)
	_ = viper.BindEnv(nonceGroupViperKey, nonceGroupEnv)
	viper.SetDefault(signerGroupViperKey, signerGroupDefault)
	_ = viper.BindEnv(signerGroupViperKey, signerGroupEnv)
	viper.SetDefault(senderGroupViperKey, senderGroupDefault)
	_ = viper.BindEnv(senderGroupViperKey, senderGroupEnv)
	viper.SetDefault(decoderGroupViperKey, decoderGroupDefault)
	_ = viper.BindEnv(decoderGroupViperKey, decoderGroupEnv)
	viper.SetDefault(bridgeGroupViperKey, bridgeGroupDefault)
	_ = viper.BindEnv(bridgeGroupViperKey, bridgeGroupEnv)

	// Kafka consumer group for wallet generation workflow
	viper.SetDefault(walletGeneratorGroupViperKey, walletGeneratorGroupDefault)
	_ = viper.BindEnv(walletGeneratorGroupViperKey, walletGeneratorGroupEnv)
	viper.SetDefault(walletGeneratedGroupViperKey, walletGeneratedGroupDefault)
	_ = viper.BindEnv(walletGeneratedGroupViperKey, walletGeneratedGroupEnv)

	// Kafka SASL
	viper.SetDefault(kafkaSASLEnableViperKey, kafkaSASLEnableDefault)
	_ = viper.BindEnv(kafkaSASLEnableViperKey, kafkaSASLEnableEnv)
	viper.SetDefault(kafkaSASLMechanismViperKey, kafkaSASLMechanismDefault)
	_ = viper.BindEnv(kafkaSASLMechanismViperKey, kafkaSASLMechanismEnv)
	viper.SetDefault(kafkaSASLHandshakeViperKey, kafkaSASLHandshakeDefault)
	_ = viper.BindEnv(kafkaSASLHandshakeViperKey, kafkaSASLHandshakeEnv)
	viper.SetDefault(kafkaSASLUserViperKey, kafkaSASLUserDefault)
	_ = viper.BindEnv(kafkaSASLUserViperKey, kafkaSASLUserEnv)
	viper.SetDefault(kafkaSASLPasswordViperKey, kafkaSASLPasswordDefault)
	_ = viper.BindEnv(kafkaSASLPasswordViperKey, kafkaSASLPasswordEnv)
	viper.SetDefault(kafkaSASLSCRAMAuthzIDViperKey, kafkaSASLSCRAMAuthzIDDefault)
	_ = viper.BindEnv(kafkaSASLSCRAMAuthzIDViperKey, kafkaSASLSCRAMAuthzIDEnv)

	// Kafka TLS
	viper.SetDefault(kafkaTLSEnableViperKey, kafkaTLSEnableDefault)
	_ = viper.BindEnv(kafkaTLSEnableViperKey, kafkaTLSEnableEnv)
	viper.SetDefault(kafkaTLSInsecureSkipVerifyViperKey, kafkaTLSInsecureSkipVerifyDefault)
	_ = viper.BindEnv(kafkaTLSInsecureSkipVerifyViperKey, kafkaTLSInsecureSkipVerifyEnv)
	viper.SetDefault(kafkaTLSClientCertFilePathViperKey, kafkaTLSClientCertFilePathDefault)
	_ = viper.BindEnv(kafkaTLSClientCertFilePathViperKey, kafkaTLSClientCertFilePathEnv)
	viper.SetDefault(kafkaTLSClientKeyFilePathViperKey, kafkaTLSClientKeyFilePathDefault)
	_ = viper.BindEnv(kafkaTLSClientKeyFilePathViperKey, kafkaTLSClientKeyFilePathEnv)
	viper.SetDefault(kafkaTLSCACertFilePathViperKey, kafkaTLSCACertFilePathDefault)
	_ = viper.BindEnv(kafkaTLSCACertFilePathViperKey, kafkaTLSCACertFilePathEnv)
}

var (
	kafkaAddressFlag     = "kafka-address"
	kafkaAddressViperKey = "kafka.addresses"
	kafkaAddressDefault  = []string{"localhost:9092"}
	kafkaAddressEnv      = "KAFKA_ADDRESS"
)

// KafkaAddresses register flag for Kafka server addresses
func KafkaAddresses(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Kafka server to connect to.
Environment variable: %q`, kafkaAddressEnv)
	f.StringSlice(kafkaAddressFlag, kafkaAddressDefault, desc)
	_ = viper.BindPFlag(kafkaAddressViperKey, f.Lookup(kafkaAddressFlag))
}

var (
	kafkaGroupFlag     = "kafka-group"
	kafkaGroupViperKey = "kafka.group"
	kafkaGroupDefault  = "group-e2e"
	kafkaGroupEnv      = "KAFKA_GROUP"
)

// KafkaGroup register flag for Kafka server addresses
func KafkaGroup(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Kafka server to connect to.
Environment variable: %q`, kafkaGroupEnv)
	f.String(kafkaGroupFlag, kafkaGroupDefault, desc)
	_ = viper.BindPFlag(kafkaGroupViperKey, f.Lookup(kafkaGroupEnv))
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

	txDecodedFlag         = "topic-decoded"
	txDecodedViperKey     = "kafka.topic.decoded"
	txDecodedTopicEnv     = "KAFKA_TOPIC_TX_DECODED"
	txDecodedTopicDefault = "topic-tx-decoded"

	txRecoverFlag         = "topic-recover"
	txRecoverViperKey     = "kafka.topic.recover"
	txRecoverTopicEnv     = "KAFKA_TOPIC_TX_RECOVER"
	txRecoverTopicDefault = "topic-tx-recover"

	walletGeneratorFlag     = "topic-wallet-generator"
	walletGeneratorViperKey = "kafka.topic.wallet.generator"
	walletGeneratorTopicEnv = "KAFKA_TOPIC_WALLET_GENERATOR"
	walletGeneratorDefault  = "topic-wallet-generator"

	walletGeneratedFlag     = "topic-wallet-generated"
	walletGeneratedViperKey = "kafka.topic.wallet.generated"
	walletGeneratedTopicEnv = "KAFKA_TOPIC_WALLET_GENERATED"
	walletGeneratedDefault  = "topic-wallet-generated"
)

// TODO: implement test for all Topics flags & Group flags

// KafkaTopicTxCrafter register flag for Kafka topic
func KafkaTopicTxCrafter(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their transaction payload crafted
Environment variable: %q`, txCrafterTopicEnv)
	f.String(txCrafterFlag, txCrafterTopicDefault, desc)
	_ = viper.BindPFlag(txCrafterViperKey, f.Lookup(txCrafterFlag))
}

// KafkaTopicTxNonce register flag for Kafka topic
func KafkaTopicTxNonce(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their transaction nonce set
Environment variable: %q`, txNonceTopicEnv)
	f.String(txNonceFlag, txNonceTopicDefault, desc)
	_ = viper.BindPFlag(txNonceViperKey, f.Lookup(txNonceFlag))
}

// KafkaTopicTxSigner register flag for Kafka topic
func KafkaTopicTxSigner(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their transaction signed
Environment variable: %q`, txSignerTopicEnv)
	f.String(txSignerFlag, txSignerTopicDefault, desc)
	_ = viper.BindPFlag(txSignerViperKey, f.Lookup(txSignerFlag))
}

// KafkaTopicTxSender register flag for Kafka topic
func KafkaTopicTxSender(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their transaction sent
Environment variable: %q`, txSenderTopicEnv)
	f.String(txSenderFlag, txSenderTopicDefault, desc)
	_ = viper.BindPFlag(txSenderViperKey, f.Lookup(txSenderFlag))
}

// KafkaTopicTxDecoder register flag for Kafka topic
func KafkaTopicTxDecoder(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their receipt decoded
Environment variable: %q`, txDecoderTopicEnv)
	f.String(txDecoderFlag, txDecoderTopicDefault, desc)
	_ = viper.BindPFlag(txDecoderViperKey, f.Lookup(txDecoderFlag))
}

// KafkaTopicTxRecover register flag for Kafka topic
func KafkaTopicTxRecover(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their transaction recovered
Environment variable: %q`, txRecoverTopicEnv)
	f.String(txRecoverFlag, txRecoverTopicDefault, desc)
	_ = viper.BindPFlag(txRecoverViperKey, f.Lookup(txRecoverFlag))
}

// KafkaTopicTxDecoded register flag for Kafka topic
func KafkaTopicTxDecoded(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages which receipt has been decoded
Environment variable: %q`, txDecodedTopicEnv)
	f.String(txDecodedFlag, txDecodedTopicDefault, desc)
	_ = viper.BindPFlag(txDecodedViperKey, f.Lookup(txDecodedFlag))
}

// KafkaTopicWalletGenerator register flag for Kafka topic
func KafkaTopicWalletGenerator(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for generating new wallets
Environment variable: %q`, walletGeneratorTopicEnv)
	f.String(walletGeneratorFlag, walletGeneratorDefault, desc)
	_ = viper.BindPFlag(walletGeneratorViperKey, f.Lookup(walletGeneratorFlag))
}

// KafkaTopicWalletGenerated register flag for Kafka topic
func KafkaTopicWalletGenerated(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for newly generated wallets
Environment variable: %q`, walletGeneratedTopicEnv)
	f.String(walletGeneratedFlag, walletGeneratedDefault, desc)
	_ = viper.BindPFlag(walletGeneratedViperKey, f.Lookup(walletGeneratedFlag))
}

// Kafka Consumer group environment variables
var (
	crafterGroupFlag     = "group-crafter"
	crafterGroupViperKey = "kafka.group.crafter"
	crafterGroupEnv      = "KAFKA_GROUP_CRAFTER"
	crafterGroupDefault  = "group-crafter"

	nonceGroupFlag     = "group-nonce"
	nonceGroupViperKey = "kafka.group.nonce"
	nonceGroupEnv      = "KAFKA_GROUP_NONCE"
	nonceGroupDefault  = "group-nonce"

	signerGroupFlag     = "group-signer"
	signerGroupViperKey = "kafka.group.signer"
	signerGroupEnv      = "KAFKA_GROUP_SIGNER"
	signerGroupDefault  = "group-signer"

	senderGroupFlag     = "group-sender"
	senderGroupViperKey = "kafka.group.sender"
	senderGroupEnv      = "KAFKA_GROUP_SENDER"
	senderGroupDefault  = "group-sender"

	decoderGroupFlag     = "group-decoder"
	decoderGroupViperKey = "kafka.group.decoder"
	decoderGroupEnv      = "KAFKA_GROUP_DECODER"
	decoderGroupDefault  = "group-decoder"

	bridgeGroupFlag     = "group-bridge"
	bridgeGroupViperKey = "kafka.group.bridge"
	bridgeGroupEnv      = "KAFKA_GROUP_BRIDGE"
	bridgeGroupDefault  = "group-bridge"

	walletGeneratorGroupFlag     = "group-wallet-generator"
	walletGeneratorGroupViperKey = "kafka.group.wallet.generator"
	walletGeneratorGroupEnv      = "KAFKA_GROUP_WALLET_GENERATOR"
	walletGeneratorGroupDefault  = "group-wallet-generator"

	walletGeneratedGroupFlag     = "group-wallet-generated"
	walletGeneratedGroupViperKey = "kafka.group.wallet.generated"
	walletGeneratedGroupEnv      = "KAFKA_GROUP_WALLET_GENERATed"
	walletGeneratedGroupDefault  = "group-wallet-generated"
)

// consumerGroupFlag register flag for a kafka consumer group
func consumerGroupFlag(f *pflag.FlagSet, flag, key, env, defaultValue string) {
	desc := fmt.Sprintf(`Kafka consumer group name
Environment variable: %q`, env)
	f.String(flag, defaultValue, desc)
	_ = viper.BindPFlag(key, f.Lookup(flag))
}

// CrafterGroup register flag for kafka crafter group
func CrafterGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, crafterGroupFlag, crafterGroupViperKey, crafterGroupEnv, crafterGroupDefault)
}

// NonceGroup register flag for kafka nonce group
func NonceGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, nonceGroupFlag, nonceGroupViperKey, nonceGroupEnv, nonceGroupDefault)
}

// SignerGroup register flag for kafka signer group
func SignerGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, signerGroupFlag, signerGroupViperKey, signerGroupEnv, signerGroupDefault)
}

// SenderGroup register flag for kafka sender group
func SenderGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, senderGroupFlag, senderGroupViperKey, senderGroupEnv, senderGroupDefault)
}

// DecoderGroup register flag for kafka decoder group
func DecoderGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, decoderGroupFlag, decoderGroupViperKey, decoderGroupEnv, decoderGroupDefault)
}

// BridgeGroup register flag for kafka decoder group
func BridgeGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, bridgeGroupFlag, bridgeGroupViperKey, bridgeGroupEnv, bridgeGroupDefault)
}

// WalletGeneratorGroup register flag for kafka decoder group
func WalletGeneratorGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, walletGeneratorGroupFlag, walletGeneratorGroupViperKey, walletGeneratorGroupEnv, walletGeneratorGroupDefault)
}

// WalletGeneratedGroup register flag for kafka decoder group
func WalletGeneratedGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, walletGeneratedGroupFlag, walletGeneratedGroupViperKey, walletGeneratedGroupEnv, walletGeneratedGroupDefault)
}

// InitKafkaSASLFlags register flags for SASL authentication
func InitKafkaSASLFlags(f *pflag.FlagSet) {
	KafkaSASLEnable(f)
	KafkaSASLMechanism(f)
	KafkaSASLHandshake(f)
	KafkaSASLUser(f)
	KafkaSASLPassword(f)
	KafkaSASLSCRAMAuthzID(f)
}

// Kafka SASL Enable environment variables
var (
	kafkaSASLEnableFlag     = "kafka-sasl-enable"
	kafkaSASLEnableViperKey = "kafka.sasl.enable"
	kafkaSASLEnableEnv      = "KAFKA_SASL_ENABLE"
	kafkaSASLEnableDefault  = false
)

// KafkaSASLEnable register flag
func KafkaSASLEnable(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Whether or not to use SASL authentication when connecting to the broker
Environment variable: %q`, kafkaSASLEnableEnv)
	f.Bool(kafkaSASLEnableFlag, kafkaSASLEnableDefault, desc)
	_ = viper.BindPFlag(kafkaSASLEnableViperKey, f.Lookup(kafkaSASLEnableFlag))
}

// Kafka SASL mechanism environment variables
var (
	kafkaSASLMechanismFlag     = "kafka-sasl-mechanism"
	kafkaSASLMechanismViperKey = "kafka.sasl.mechanism"
	kafkaSASLMechanismEnv      = "KAFKA_SASL_MECHANISM"
	kafkaSASLMechanismDefault  string
)

// KafkaSASLMechanism register flag
func KafkaSASLMechanism(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`SASLMechanism is the name of the enabled SASL mechanism. Possible values: OAUTHBEARER, PLAIN (defaults to PLAIN).
Environment variable: %q`, kafkaSASLMechanismEnv)
	f.String(kafkaSASLMechanismFlag, kafkaSASLMechanismDefault, desc)
	_ = viper.BindPFlag(kafkaSASLMechanismViperKey, f.Lookup(kafkaSASLMechanismFlag))
}

// Kafka SASL Handshake environment variables
var (
	kafkaSASLHandshakeFlag     = "kafka-sasl-handshake"
	kafkaSASLHandshakeViperKey = "kafka.sasl.handshake"
	kafkaSASLHandshakeEnv      = "KAFKA_SASL_HANDSHAKE"
	kafkaSASLHandshakeDefault  = true
)

// KafkaSASLHandshake register flag
func KafkaSASLHandshake(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Whether or not to send the Kafka SASL handshake first if enabled (defaults to true). You should only set this to false if you're using a non-Kafka SASL proxy.
Environment variable: %q`, kafkaSASLHandshakeEnv)
	f.Bool(kafkaSASLHandshakeFlag, kafkaSASLHandshakeDefault, desc)
	_ = viper.BindPFlag(kafkaSASLHandshakeViperKey, f.Lookup(kafkaSASLHandshakeFlag))
}

// Kafka SASL User environment variables
var (
	kafkaSASLUserFlag     = "kafka-sasl-user"
	kafkaSASLUserViperKey = "kafka.sasl.user"
	kafkaSASLUserEnv      = "KAFKA_SASL_USER"
	kafkaSASLUserDefault  string
)

// KafkaSASLUser register flag
func KafkaSASLUser(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Username for SASL/PLAIN or SASL/SCRAM authentication.
Environment variable: %q`, kafkaSASLUserEnv)
	f.String(kafkaSASLUserFlag, kafkaSASLUserDefault, desc)
	_ = viper.BindPFlag(kafkaSASLUserViperKey, f.Lookup(kafkaSASLUserFlag))
}

// Kafka SASL Password environment variables
var (
	kafkaSASLPasswordFlag     = "kafka-sasl-password"
	kafkaSASLPasswordViperKey = "kafka.sasl.password"
	kafkaSASLPasswordEnv      = "KAFKA_SASL_PASSWORD"
	kafkaSASLPasswordDefault  string
)

// KafkaSASLPassword register flag
func KafkaSASLPassword(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Password for SASL/PLAIN or SASL/SCRAM authentication.
Environment variable: %q`, kafkaSASLPasswordEnv)
	f.String(kafkaSASLPasswordFlag, kafkaSASLPasswordDefault, desc)
	_ = viper.BindPFlag(kafkaSASLPasswordViperKey, f.Lookup(kafkaSASLPasswordFlag))
}

// Kafka SASL SCRAMAuthzID environment variables
var (
	kafkaSASLSCRAMAuthzIDFlag     = "kafka-sasl-scramauthzid"
	kafkaSASLSCRAMAuthzIDViperKey = "kafka.sasl.scramauthzid"
	kafkaSASLSCRAMAuthzIDEnv      = "KAFKA_SASL_SCRAMAUTHZID"
	kafkaSASLSCRAMAuthzIDDefault  string
)

// KafkaSASLSCRAMAuthzID register flag
func KafkaSASLSCRAMAuthzID(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Authz id used for SASL/SCRAM authentication
Environment variable: %q`, kafkaSASLSCRAMAuthzIDEnv)
	f.String(kafkaSASLSCRAMAuthzIDFlag, kafkaSASLSCRAMAuthzIDDefault, desc)
	_ = viper.BindPFlag(kafkaSASLSCRAMAuthzIDViperKey, f.Lookup(kafkaSASLSCRAMAuthzIDFlag))
}

// InitKafkaSASLTLSFlags register flags for SASL and SSL
func InitKafkaSASLTLSFlags(f *pflag.FlagSet) {
	KafkaTLSEnable(f)
	KafkaTLSInsecureSkipVerify(f)
	KafkaTLSClientCertFilePath(f)
	KafkaTLSClientKeyFilePath(f)
	KafkaTLSCaCertFilePath(f)
}

// Kafka TLS Enable environment variables
var (
	kafkaTLSEnableFlag     = "kafka-tls-enabled"
	kafkaTLSEnableViperKey = "kafka.tls.enabled"
	kafkaTLSEnableEnv      = "KAFKA_TLS_ENABLED"
	kafkaTLSEnableDefault  = false
)

// KafkaTLSEnable register flag
func KafkaTLSEnable(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Whether or not to use TLS when connecting to the broker (defaults to false).
Environment variable: %q`, kafkaTLSEnableEnv)
	f.Bool(kafkaTLSEnableFlag, kafkaTLSEnableDefault, desc)
	_ = viper.BindPFlag(kafkaTLSEnableViperKey, f.Lookup(kafkaTLSEnableFlag))
}

// Kafka TLS InsecureSkipVerify environment variables
var (
	kafkaTLSInsecureSkipVerifyFlag     = "kafka-tls-insecureSkipVerify"
	kafkaTLSInsecureSkipVerifyViperKey = "kafka.tls.insecureSkipVerify"
	kafkaTLSInsecureSkipVerifyEnv      = "KAFKA_TLS_INSECURESKIPVERIFY"
	kafkaTLSInsecureSkipVerifyDefault  = false
)

// KafkaTLSInsecureSkipVerify register flag
func KafkaTLSInsecureSkipVerify(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Controls whether a client verifies the server's certificate chain and host name. If InsecureSkipVerify is true, TLS accepts any certificate presented by the server and any host name in that certificate. In this mode, TLS is susceptible to man-in-the-middle attacks. This should be used only for testing.
Environment variable: %q`, kafkaTLSInsecureSkipVerifyEnv)
	f.Bool(kafkaTLSInsecureSkipVerifyFlag, kafkaTLSInsecureSkipVerifyDefault, desc)
	_ = viper.BindPFlag(kafkaTLSInsecureSkipVerifyViperKey, f.Lookup(kafkaTLSInsecureSkipVerifyFlag))
}

// Kafka TLS ClientCertFilePath environment variables
var (
	kafkaTLSClientCertFilePathFlag     = "kafka-tls-clientcertfilepath"
	kafkaTLSClientCertFilePathViperKey = "kafka.tls.clientCertfilepath"
	kafkaTLSClientCertFilePathEnv      = "KAFKA_TLS_CLIENTCERTFILEPATH"
	kafkaTLSClientCertFilePathDefault  string
)

// KafkaTLSClientCertFilePath register flag
func KafkaTLSClientCertFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Client Cert File Path.
Environment variable: %q`, kafkaTLSClientCertFilePathEnv)
	f.String(kafkaTLSClientCertFilePathFlag, kafkaTLSClientCertFilePathDefault, desc)
	_ = viper.BindPFlag(kafkaTLSClientCertFilePathViperKey, f.Lookup(kafkaTLSClientCertFilePathFlag))
}

// Kafka TLS ClientKeyFilePath environment variables
var (
	kafkaTLSClientKeyFilePathFlag     = "kafka-tls-clientkeyfilepath"
	kafkaTLSClientKeyFilePathViperKey = "kafka.tls.clientkeyfilepath"
	kafkaTLSClientKeyFilePathEnv      = "KAFKA_TLS_CLIENTKEYFILEPATH"
	kafkaTLSClientKeyFilePathDefault  string
)

// KafkaTLSClientKeyFilePath register flag
func KafkaTLSClientKeyFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Client key file Path.
Environment variable: %q`, kafkaTLSClientKeyFilePathEnv)
	f.String(kafkaTLSClientKeyFilePathFlag, kafkaTLSClientKeyFilePathDefault, desc)
	_ = viper.BindPFlag(kafkaTLSClientKeyFilePathViperKey, f.Lookup(kafkaTLSClientKeyFilePathFlag))
}

// Kafka TLS CACertFilePath environment variables
var (
	kafkaTLSCACertFilePathFlag     = "kafka-tls-cacertfilepath"
	kafkaTLSCACertFilePathViperKey = "kafka.tls.cacertfilepath"
	kafkaTLSCACertFilePathEnv      = "KAFKA_TLS_CACERTFILEPATH"
	kafkaTLSCACertFilePathDefault  string
)

// KafkaTLSCaCertFilePath register flag
func KafkaTLSCaCertFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`CA cert file Path.
Environment variable: %q`, kafkaTLSCACertFilePathEnv)
	f.String(kafkaTLSCACertFilePathFlag, kafkaTLSCACertFilePathDefault, desc)
	_ = viper.BindPFlag(kafkaTLSCACertFilePathViperKey, f.Lookup(kafkaTLSCACertFilePathFlag))
}
