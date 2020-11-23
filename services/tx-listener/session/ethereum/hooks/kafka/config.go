package kafka

import (
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
)

type Config struct {
	OutTopic string
}

func NewConfig() *Config {
	return &Config{
		OutTopic: viper.GetString(broker.TxDecodedViperKey),
	}
}
