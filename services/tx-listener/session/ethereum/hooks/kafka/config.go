package kafka

import (
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
)

type Config struct {
	TopicTxDecoder string
}

func NewConfig() *Config {
	return &Config{
		TopicTxDecoder: viper.GetString(broker.TxDecoderViperKey),
	}
}
