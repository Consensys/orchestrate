package txsigner

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
)

func init() {
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)

	viper.SetDefault(kafkaTopicViperKey, kafkaTopicDefault)
	_ = viper.BindEnv(kafkaTopicViperKey, kafkaTopicEnv)

	viper.SetDefault(kafkaSenderTopicViperKey, kafkaSenderTopicDefault)
	_ = viper.BindEnv(kafkaSenderTopicViperKey, kafkaSenderTopicEnv)

	viper.SetDefault(kafkaRecoverTopicViperKey, kafkaRecoverTopicDefault)
	_ = viper.BindEnv(kafkaRecoverTopicViperKey, kafkaRecoverTopicEnv)
}

const (
	MetricsURLViperKey = "tx-signer.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "TX_SIGNER_METRICS_URL"

	kafkaTopicFlag     = "topic-tx-signer"
	kafkaTopicViperKey = "topic.tx.signer"
	kafkaTopicEnv      = "TOPIC_TX_SIGNER"
	kafkaTopicDefault  = "topic-tx-signer"

	kafkaSenderTopicFlag     = "topic-tx-sender"
	kafkaSenderTopicViperKey = "topic.tx.sender"
	kafkaSenderTopicEnv      = "TOPIC_TX_SENDER"
	kafkaSenderTopicDefault  = "topic-tx-sender"

	kafkaRecoverTopicFlag     = "topic-tx-recover"
	kafkaRecoverTopicViperKey = "topic.tx.recover"
	kafkaRecoverTopicEnv      = "TOPIC_TX_RECOVER"
	kafkaRecoverTopicDefault  = "topic-tx-recover"
)

// Flags register flags for tx sentry
func Flags(f *pflag.FlagSet) {
	broker.InitKafkaFlags(f)
	keymanager.Flags(f)
	kafkaTopicTxSignerFlag(f)
	kafkaTopicSenderFlag(f)
	kafkaTopicRecoverFlag(f)
}

type Config struct {
	App           *app.Config
	GroupName     string
	ListenerTopic string
	SenderTopic   string
	RecoverTopic  string
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:           app.NewConfig(vipr),
		GroupName:     "group-signer",
		ListenerTopic: vipr.GetString(kafkaTopicViperKey),
		SenderTopic:   vipr.GetString(kafkaSenderTopicViperKey),
		RecoverTopic:  vipr.GetString(kafkaRecoverTopicViperKey),
	}
}

// kafkaTopicTxSigner register flag for Kafka topic
func kafkaTopicTxSignerFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes waiting for their transaction signed. Environment variable: %q`, kafkaTopicEnv)
	f.String(kafkaTopicFlag, kafkaTopicDefault, desc)
	_ = viper.BindPFlag(kafkaTopicViperKey, f.Lookup(kafkaTopicFlag))
}

// kafkaTopicSenderFlag register flag for Kafka sender topic
func kafkaTopicSenderFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes send to the tx-sender. Environment variable: %q`, kafkaSenderTopicEnv)
	f.String(kafkaSenderTopicFlag, kafkaSenderTopicDefault, desc)
	_ = viper.BindPFlag(kafkaSenderTopicViperKey, f.Lookup(kafkaSenderTopicFlag))
}

// kafkaTopicRecoverFlag register flag for Kafka recover topic
func kafkaTopicRecoverFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic for envelopes of failed transactions. Environment variable: %q`, kafkaRecoverTopicEnv)
	f.String(kafkaRecoverTopicFlag, kafkaRecoverTopicDefault, desc)
	_ = viper.BindPFlag(kafkaRecoverTopicViperKey, f.Lookup(kafkaRecoverTopicFlag))
}
