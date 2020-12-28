package sarama

import (
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	// Kafka general parameters
	viper.SetDefault(KafkaURLViperKey, kafkaURLDefault)
	_ = viper.BindEnv(KafkaURLViperKey, KafkaURLEnv)
	viper.SetDefault(KafkaGroupViperKey, kafkaGroupDefault)
	_ = viper.BindEnv(KafkaGroupViperKey, kafkaGroupEnv)

	// Kafka topics for the tx workflow
	viper.SetDefault(TxSenderViperKey, txSenderTopicDefault)
	_ = viper.BindEnv(TxSenderViperKey, txSenderTopicEnv)
	viper.SetDefault(TxDecodedViperKey, txDecodedTopicDefault)
	_ = viper.BindEnv(TxDecodedViperKey, txDecodedTopicEnv)
	viper.SetDefault(TxRecoverViperKey, txRecoverTopicDefault)
	_ = viper.BindEnv(TxRecoverViperKey, txRecoverTopicEnv)

	// Kafka consumer groups for tx workflow
	viper.SetDefault(CrafterGroupViperKey, crafterGroupDefault)
	_ = viper.BindEnv(CrafterGroupViperKey, crafterGroupEnv)
	viper.SetDefault(SignerGroupViperKey, signerGroupDefault)
	_ = viper.BindEnv(SignerGroupViperKey, signerGroupEnv)
	viper.SetDefault(SenderGroupViperKey, senderGroupDefault)
	_ = viper.BindEnv(SenderGroupViperKey, senderGroupEnv)
	viper.SetDefault(DecodedGroupViperKey, decodedGroupDefault)
	_ = viper.BindEnv(DecodedGroupViperKey, decodedGroupEnv)

	// Kafka SASL
	viper.SetDefault(kafkaSASLEnabledViperKey, kafkaSASLEnabledDefault)
	_ = viper.BindEnv(kafkaSASLEnabledViperKey, kafkaSASLEnabledEnv)
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

	// Kafka
	viper.SetDefault(kafkaVersionViperKey, kafkaVersionDefault)
	_ = viper.BindEnv(kafkaVersionViperKey, kafkaVersionEnv)

	// Kafka Consumer
	viper.SetDefault(kafkaConsumerMaxWaitTimeViperKey, kafkaConsumerMaxWaitTimeDefault)
	_ = viper.BindEnv(kafkaConsumerMaxWaitTimeViperKey, kafkaConsumerMaxWaitTimeEnv)
	viper.SetDefault(kafkaConsumerMaxProcessingTimeViperKey, kafkaConsumerMaxProcessingTimeDefault)
	_ = viper.BindEnv(kafkaConsumerMaxProcessingTimeViperKey, kafkaConsumerMaxProcessingTimeEnv)
	viper.SetDefault(kafkaConsumerGroupSessionTimeoutViperKey, kafkaConsumerGroupSessionTimeoutDefault)
	_ = viper.BindEnv(kafkaConsumerGroupSessionTimeoutViperKey, kafkaConsumerGroupSessionTimeoutEnv)
	viper.SetDefault(kafkaConsumerGroupHeartbeatIntervalViperKey, kafkaConsumerGroupHeartbeatIntervalDefault)
	_ = viper.BindEnv(kafkaConsumerGroupHeartbeatIntervalViperKey, kafkaConsumerGroupHeartbeatIntervalEnv)
	viper.SetDefault(kafkaConsumerGroupRebalanceTimeoutViperKey, kafkaConsumerGroupRebalanceTimeoutDefault)
	_ = viper.BindEnv(kafkaConsumerGroupRebalanceTimeoutViperKey, kafkaConsumerGroupRebalanceTimeoutEnv)
}

// InitKafkaFlags
func InitKafkaFlags(f *pflag.FlagSet) {
	KafkaURL(f)
	KafkaGroup(f)
	InitKafkaSASLFlags(f)
	InitKafkaTLSFlags(f)
	KafkaConsumerMaxWaitTime(f)
	kafkaConsumerMaxProcessingTime(f)
	kafkaConsumerGroupSessionTimeout(f)
	kafkaConsumerGroupHeartbeatInterval(f)
	kafkaConsumerGroupRebalanceTimeout(f)
	kafkaVersion(f)
	kafkaConsumerGroupRebalanceStrategy(f)
}

var (
	kafkaURLFlag     = "kafka-url"
	KafkaURLViperKey = "kafka.url"
	kafkaURLDefault  = []string{"localhost:9092"}
	KafkaURLEnv      = "KAFKA_URL"
)

// KafkaURL register flag for Kafka server
func KafkaURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (addresses) of Kafka server(s) to connect to.
Environment variable: %q`, KafkaURLEnv)
	f.StringSlice(kafkaURLFlag, kafkaURLDefault, desc)
	_ = viper.BindPFlag(KafkaURLViperKey, f.Lookup(kafkaURLFlag))
}

const (
	kafkaGroupFlag     = "kafka-group"
	KafkaGroupViperKey = "kafka.group"
	kafkaGroupEnv      = "KAFKA_GROUP"
	kafkaGroupDefault  = "group-e2e"
)

// KafkaGroup register flag for Kafka group
func KafkaGroup(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Kafka server to connect to.
Environment variable: %q`, kafkaGroupEnv)
	f.String(kafkaGroupFlag, kafkaGroupDefault, desc)
	_ = viper.BindPFlag(KafkaGroupViperKey, f.Lookup(kafkaGroupFlag))
}

const (
	txSenderFlag         = "topic-tx-sender"
	TxSenderViperKey     = "topic.tx.sender"
	txSenderTopicEnv     = "TOPIC_TX_SENDER"
	txSenderTopicDefault = "topic-tx-sender"

	txDecodedFlag         = "topic-tx-decoded"
	TxDecodedViperKey     = "topic.tx.decoded"
	txDecodedTopicEnv     = "TOPIC_TX_DECODED"
	txDecodedTopicDefault = "topic-tx-decoded"

	txRecoverFlag         = "topic-tx-recover"
	TxRecoverViperKey     = "topic.tx.recover"
	txRecoverTopicEnv     = "TOPIC_TX_RECOVER"
	txRecoverTopicDefault = "topic-tx-recover"
)

type KafkaTopicConfig struct {
	Sender  string
	Decoded string
	Recover string
}

func NewKafkaTopicConfig(vipr *viper.Viper) *KafkaTopicConfig {
	return &KafkaTopicConfig{
		Sender:  vipr.GetString(TxSenderViperKey),
		Decoded: vipr.GetString(TxDecodedViperKey),
		Recover: vipr.GetString(TxRecoverViperKey),
	}
}

// TODO: implement test for all Topics flags & Group flags

// KafkaTopicTxSender register flag for Kafka topic
func KafkaTopicTxSender(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their transaction sent
Environment variable: %q`, txSenderTopicEnv)
	f.String(txSenderFlag, txSenderTopicDefault, desc)
	_ = viper.BindPFlag(TxSenderViperKey, f.Lookup(txSenderFlag))
}

// KafkaTopicTxRecover register flag for Kafka topic
func KafkaTopicTxRecover(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their transaction recovered
Environment variable: %q`, txRecoverTopicEnv)
	f.String(txRecoverFlag, txRecoverTopicDefault, desc)
	_ = viper.BindPFlag(TxRecoverViperKey, f.Lookup(txRecoverFlag))
}

// KafkaTopicTxDecoded register flag for Kafka topic
func KafkaTopicTxDecoded(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for messages which receipt has been decoded
Environment variable: %q`, txDecodedTopicEnv)
	f.String(txDecodedFlag, txDecodedTopicDefault, desc)
	_ = viper.BindPFlag(TxDecodedViperKey, f.Lookup(txDecodedFlag))
}

// Kafka Consumer group environment variables
const (
	crafterGroupFlag     = "group-crafter"
	CrafterGroupViperKey = "kafka.group.crafter"
	crafterGroupEnv      = "KAFKA_GROUP_CRAFTER"
	crafterGroupDefault  = "group-crafter"

	signerGroupFlag     = "group-signer"
	SignerGroupViperKey = "kafka.group.signer"
	signerGroupEnv      = "KAFKA_GROUP_SIGNER"
	signerGroupDefault  = "group-signer"

	senderGroupFlag     = "group-sender"
	SenderGroupViperKey = "kafka.group.sender"
	senderGroupEnv      = "KAFKA_GROUP_SENDER"
	senderGroupDefault  = "group-sender"

	decodedGroupFlag     = "group-decoded"
	DecodedGroupViperKey = "kafka.group.decoded"
	decodedGroupEnv      = "KAFKA_GROUP_DECODED"
	decodedGroupDefault  = "group-decoded"
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
	consumerGroupFlag(f, crafterGroupFlag, CrafterGroupViperKey, crafterGroupEnv, crafterGroupDefault)
}

// SignerGroup register flag for kafka signer group
func SignerGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, signerGroupFlag, SignerGroupViperKey, signerGroupEnv, signerGroupDefault)
}

// SenderGroup register flag for kafka sender group
func SenderGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, senderGroupFlag, SenderGroupViperKey, senderGroupEnv, senderGroupDefault)
}

// DecoderGroup register flag for kafka decoder group
func DecoderGroup(f *pflag.FlagSet) {
	consumerGroupFlag(f, decodedGroupFlag, DecodedGroupViperKey, decodedGroupEnv, decodedGroupDefault)
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
const (
	kafkaSASLEnabledFlag     = "kafka-sasl-enabled"
	kafkaSASLEnabledViperKey = "kafka.sasl.enabled"
	kafkaSASLEnabledEnv      = "KAFKA_SASL_ENABLED"
	kafkaSASLEnabledDefault  = false
)

// KafkaSASLEnable register flag
func KafkaSASLEnable(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Whether or not to use SASL authentication when connecting to the broker
Environment variable: %q`, kafkaSASLEnabledEnv)
	f.Bool(kafkaSASLEnabledFlag, kafkaSASLEnabledDefault, desc)
	_ = viper.BindPFlag(kafkaSASLEnabledViperKey, f.Lookup(kafkaSASLEnabledFlag))
}

// Kafka SASL mechanism environment variables
const (
	kafkaSASLMechanismFlag     = "kafka-sasl-mechanism"
	kafkaSASLMechanismViperKey = "kafka.sasl.mechanism"
	kafkaSASLMechanismEnv      = "KAFKA_SASL_MECHANISM"
	kafkaSASLMechanismDefault  = ""
)

// KafkaSASLMechanism register flag
func KafkaSASLMechanism(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`SASLMechanism is the name of the enabled SASL mechanism. Possible values: OAUTHBEARER, PLAIN (defaults to PLAIN).
Environment variable: %q`, kafkaSASLMechanismEnv)
	f.String(kafkaSASLMechanismFlag, kafkaSASLMechanismDefault, desc)
	_ = viper.BindPFlag(kafkaSASLMechanismViperKey, f.Lookup(kafkaSASLMechanismFlag))
}

// Kafka SASL Handshake environment variables
const (
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
const (
	kafkaSASLUserFlag     = "kafka-sasl-user"
	kafkaSASLUserViperKey = "kafka.sasl.user"
	kafkaSASLUserEnv      = "KAFKA_SASL_USER"
	kafkaSASLUserDefault  = ""
)

// KafkaSASLUser register flag
func KafkaSASLUser(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Username for SASL/PLAIN or SASL/SCRAM auth.
Environment variable: %q`, kafkaSASLUserEnv)
	f.String(kafkaSASLUserFlag, kafkaSASLUserDefault, desc)
	_ = viper.BindPFlag(kafkaSASLUserViperKey, f.Lookup(kafkaSASLUserFlag))
}

// Kafka SASL Password environment variables
const (
	kafkaSASLPasswordFlag     = "kafka-sasl-password"
	kafkaSASLPasswordViperKey = "kafka.sasl.password"
	kafkaSASLPasswordEnv      = "KAFKA_SASL_PASSWORD"
	kafkaSASLPasswordDefault  = ""
)

// KafkaSASLPassword register flag
func KafkaSASLPassword(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Password for SASL/PLAIN or SASL/SCRAM auth.
Environment variable: %q`, kafkaSASLPasswordEnv)
	f.String(kafkaSASLPasswordFlag, kafkaSASLPasswordDefault, desc)
	_ = viper.BindPFlag(kafkaSASLPasswordViperKey, f.Lookup(kafkaSASLPasswordFlag))
}

// Kafka SASL SCRAMAuthzID environment variables
const (
	kafkaSASLSCRAMAuthzIDFlag     = "kafka-sasl-scramauthzid"
	kafkaSASLSCRAMAuthzIDViperKey = "kafka.sasl.scramauthzid"
	kafkaSASLSCRAMAuthzIDEnv      = "KAFKA_SASL_SCRAMAUTHZID"
	kafkaSASLSCRAMAuthzIDDefault  = ""
)

// KafkaSASLSCRAMAuthzID register flag
func KafkaSASLSCRAMAuthzID(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Authz id used for SASL/SCRAM authentication
Environment variable: %q`, kafkaSASLSCRAMAuthzIDEnv)
	f.String(kafkaSASLSCRAMAuthzIDFlag, kafkaSASLSCRAMAuthzIDDefault, desc)
	_ = viper.BindPFlag(kafkaSASLSCRAMAuthzIDViperKey, f.Lookup(kafkaSASLSCRAMAuthzIDFlag))
}

// InitKafkaTLSFlags register flags for SASL and SSL
func InitKafkaTLSFlags(f *pflag.FlagSet) {
	KafkaTLSEnable(f)
	KafkaTLSInsecureSkipVerify(f)
	KafkaTLSClientCertFilePath(f)
	KafkaTLSClientKeyFilePath(f)
	KafkaTLSCaCertFilePath(f)
}

// Kafka TLS Enable environment variables
const (
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
const (
	kafkaTLSInsecureSkipVerifyFlag     = "kafka-tls-insecure-skip-verify"
	kafkaTLSInsecureSkipVerifyViperKey = "kafka.tls.insecure.skip.verify"
	kafkaTLSInsecureSkipVerifyEnv      = "KAFKA_TLS_INSECURE_SKIP_VERIFY"
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
const (
	kafkaTLSClientCertFilePathFlag     = "kafka-tls-client-cert-file"
	kafkaTLSClientCertFilePathViperKey = "kafka.tls.client.cert.file"
	kafkaTLSClientCertFilePathEnv      = "KAFKA_TLS_CLIENT_CERT_FILE"
	kafkaTLSClientCertFilePathDefault  = ""
)

// KafkaTLSClientCertFilePath register flag
func KafkaTLSClientCertFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Client Cert File Path.
Environment variable: %q`, kafkaTLSClientCertFilePathEnv)
	f.String(kafkaTLSClientCertFilePathFlag, kafkaTLSClientCertFilePathDefault, desc)
	_ = viper.BindPFlag(kafkaTLSClientCertFilePathViperKey, f.Lookup(kafkaTLSClientCertFilePathFlag))
}

// Kafka TLS ClientKeyFilePath environment variables
const (
	kafkaTLSClientKeyFilePathFlag     = "kafka-tls-client-key-file"
	kafkaTLSClientKeyFilePathViperKey = "kafka.tls.client.key.file"
	kafkaTLSClientKeyFilePathEnv      = "KAFKA_TLS_CLIENT_KEY_FILE"
	kafkaTLSClientKeyFilePathDefault  = ""
)

// KafkaTLSClientKeyFilePath register flag
func KafkaTLSClientKeyFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Client key file Path.
Environment variable: %q`, kafkaTLSClientKeyFilePathEnv)
	f.String(kafkaTLSClientKeyFilePathFlag, kafkaTLSClientKeyFilePathDefault, desc)
	_ = viper.BindPFlag(kafkaTLSClientKeyFilePathViperKey, f.Lookup(kafkaTLSClientKeyFilePathFlag))
}

// Kafka TLS CACertFilePath environment variables
const (
	kafkaTLSCACertFilePathFlag     = "kafka-tls-ca-cert-file"
	kafkaTLSCACertFilePathViperKey = "kafka.tls.ca.cert.file"
	kafkaTLSCACertFilePathEnv      = "KAFKA_TLS_CA_CERT_FILE"
	kafkaTLSCACertFilePathDefault  = ""
)

// KafkaTLSCaCertFilePath register flag
func KafkaTLSCaCertFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`CA cert file Path.
Environment variable: %q`, kafkaTLSCACertFilePathEnv)
	f.String(kafkaTLSCACertFilePathFlag, kafkaTLSCACertFilePathDefault, desc)
	_ = viper.BindPFlag(kafkaTLSCACertFilePathViperKey, f.Lookup(kafkaTLSCACertFilePathFlag))
}

// Kafka Consumer MaxWaitTime wait time environment variables
const (
	kafkaConsumerMaxWaitTimeViperFlag = "kafka-consumer-max-wait-time"
	kafkaConsumerMaxWaitTimeViperKey  = "kafka.consumer.max.wait.time"
	kafkaConsumerMaxWaitTimeEnv       = "KAFKA_CONSUMER_MAX_WAIT_TIME"
	kafkaConsumerMaxWaitTimeDefault   = time.Millisecond * 250
)

// KafkaConsumerMaxWaitTime configuration
func KafkaConsumerMaxWaitTime(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`The maximum amount of time the broker will wait for Consumer.Fetch.Min bytes to become available before it returns fewer than that anyways.
Environment variable: %q in ms`, kafkaConsumerMaxWaitTimeEnv)
	f.Duration(kafkaConsumerMaxWaitTimeViperFlag, kafkaConsumerMaxWaitTimeDefault, desc)
	_ = viper.BindPFlag(kafkaConsumerMaxWaitTimeViperKey, f.Lookup(kafkaConsumerMaxWaitTimeViperFlag))
}

// kafkaConsumerMaxProcessingTime environment variables
const (
	kafkaConsumerMaxProcessingTimeFlag     = "kafka-consumer-max-processing-time"
	kafkaConsumerMaxProcessingTimeViperKey = "kafka.consumer.maxprocessingtime"
	kafkaConsumerMaxProcessingTimeEnv      = "KAFKA_CONSUMER_MAXPROCESSINGTIME"
	kafkaConsumerMaxProcessingTimeDefault  = time.Millisecond * 100
)

// kafkaConsumerMaxProcessingTime configuration
func kafkaConsumerMaxProcessingTime(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`The maximum amount of time the consumer expects a message takes to process for the user. If writing to the Messages channel takes longer than this, that partition will stop fetching more messages until it can proceed again.
Environment variable: %q`, kafkaConsumerMaxProcessingTimeEnv)
	f.Duration(kafkaConsumerMaxProcessingTimeFlag, kafkaConsumerMaxProcessingTimeDefault, desc)
	_ = viper.BindPFlag(kafkaConsumerMaxProcessingTimeViperKey, f.Lookup(kafkaConsumerMaxProcessingTimeFlag))
}

// kafkaConsumerGroupSessionTimeout environment variables
const (
	kafkaConsumerGroupSessionTimeoutFlag     = "kafka-consumer-group-session-timeout"
	kafkaConsumerGroupSessionTimeoutViperKey = "kafka.consumer.group.session.timeout"
	kafkaConsumerGroupSessionTimeoutEnv      = "KAFKA_CONSUMER_GROUP_SESSION_TIMEOUT"
	kafkaConsumerGroupSessionTimeoutDefault  = time.Second * 10
)

// kafkaConsumerGroupSessionTimeout configuration
func kafkaConsumerGroupSessionTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`The timeout used to detect consumer failures when using Kafka's group management facility. The consumer sends periodic heartbeats to indicate its liveness to the broker. If no heartbeats are received by the broker before the expiration of this session timeout, then the broker will remove this consumer from the group and initiate a rebalance.
Environment variable: %q`, kafkaConsumerGroupSessionTimeoutEnv)
	f.Duration(kafkaConsumerGroupSessionTimeoutFlag, kafkaConsumerGroupSessionTimeoutDefault, desc)
	_ = viper.BindPFlag(kafkaConsumerGroupSessionTimeoutViperKey, f.Lookup(kafkaConsumerGroupSessionTimeoutFlag))
}

// kafkaConsumerGroupHeartbeatInterval environment variables
const (
	kafkaConsumerGroupHeartbeatIntervalFlag     = "kafka-consumer-group-heartbeat-interval"
	kafkaConsumerGroupHeartbeatIntervalViperKey = "kafka.consumer.group.heartbeat.interval"
	kafkaConsumerGroupHeartbeatIntervalEnv      = "KAFKA_CONSUMER_GROUP_HEARTBEAT_INTERVAL"
	kafkaConsumerGroupHeartbeatIntervalDefault  = time.Second * 3
)

// kafkaConsumerGroupHeartbeatInterval configuration
func kafkaConsumerGroupHeartbeatInterval(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`The expected time between heartbeats to the consumer coordinator when using Kafka's group management facilities. Heartbeats are used to ensure that the consumer's session stays active and to facilitate rebalancing when new consumers join or leave the group.
Environment variable: %q in ms`, kafkaConsumerGroupHeartbeatIntervalEnv)
	f.Duration(kafkaConsumerGroupHeartbeatIntervalFlag, kafkaConsumerGroupHeartbeatIntervalDefault, desc)
	_ = viper.BindPFlag(kafkaConsumerGroupHeartbeatIntervalViperKey, f.Lookup(kafkaConsumerGroupHeartbeatIntervalFlag))
}

// kafkaConsumerGroupRebalanceTimeout environment variables
const (
	kafkaConsumerGroupRebalanceTimeoutFlag     = "kafka-consumer-group-rebalance-timeout"
	kafkaConsumerGroupRebalanceTimeoutViperKey = "kafka.consumer.group.rebalance.timeout"
	kafkaConsumerGroupRebalanceTimeoutEnv      = "KAFKA_CONSUMER_GROUP_REBALANCE_TIMEOUT"
	kafkaConsumerGroupRebalanceTimeoutDefault  = time.Second * 60
)

// kafkaConsumerGroupRebalanceTimeout configuration
func kafkaConsumerGroupRebalanceTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`The maximum allowed time for each worker to join the group once a rebalance has begun. This is basically a limit on the amount of time needed for all tasks to flush any pending data and commit offsets. If the timeout is exceeded, then the worker will be removed from the group, which will cause offset commit failures
Environment variable: %q`, kafkaConsumerGroupRebalanceTimeoutEnv)
	f.Duration(kafkaConsumerGroupRebalanceTimeoutFlag, kafkaConsumerGroupRebalanceTimeoutDefault, desc)
	_ = viper.BindPFlag(kafkaConsumerGroupRebalanceTimeoutViperKey, f.Lookup(kafkaConsumerGroupRebalanceTimeoutFlag))
}

// kafkaVersion environment variables
const (
	kafkaVersionFlag     = "kafka-version"
	kafkaVersionViperKey = "kafka.version"
	kafkaVersionEnv      = "KAFKA_VERSION"
	kafkaVersionDefault  = "1.0.0"
)

// kafkaVersion configuration
func kafkaVersion(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`The version of Kafka that Sarama will assume it is running against. Defaults to the oldest supported stable version. Since Kafka provides backwards-compatibility, setting it to a version older than you have will not break anything, although it may prevent you from using the latest features. Setting it to a version greater than you are actually running may lead to random breakage.
Environment variable: %q`, kafkaConsumerGroupRebalanceTimeoutEnv)
	f.String(kafkaVersionFlag, kafkaVersionDefault, desc)
	_ = viper.BindPFlag(kafkaVersionViperKey, f.Lookup(kafkaVersionFlag))
}

// kafkaConsumerGroupRebalanceStrategy environment variables
const (
	kafkaConsumerGroupRebalanceStrategyFlag     = "kafka-consumer-group-rebalance-strategy"
	kafkaConsumerGroupRebalanceStrategyViperKey = "kafka.consumer.group.rebalance.strategy"
	kafkaConsumerGroupRebalanceStrategyEnv      = "KAFKA_CONSUMER_GROUP_REBALANCE_STRATEGY"
	kafkaConsumerGroupRebalanceStrategyDefault  = "Range"
)

// kafkaConsumerGroupRebalanceStrategy configuration
func kafkaConsumerGroupRebalanceStrategy(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Strategy for allocating topic partitions to members (one of %q).
Environment variable: %q`, reflect.ValueOf(rebalanceStrategy).MapKeys(), kafkaConsumerGroupRebalanceStrategyEnv)
	f.String(kafkaConsumerGroupRebalanceStrategyFlag, kafkaConsumerGroupRebalanceStrategyDefault, desc)
	_ = viper.BindPFlag(kafkaConsumerGroupRebalanceStrategyViperKey, f.Lookup(kafkaConsumerGroupRebalanceStrategyFlag))
}
