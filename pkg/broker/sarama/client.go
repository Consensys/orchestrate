package sarama

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// NewClient creates a new sarama client and connects to one of the given broker addresses
func NewClient(addrs []string, conf *sarama.Config) (sarama.Client, error) {
	if err := conf.Validate(); err != nil {
		return nil, errors.ConfigError(err.Error()).SetComponent(component)
	}

	client, err := sarama.NewClient(addrs, conf)
	if err != nil {
		return nil, errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}

	return client, nil
}

func NewSaramaConfig() (*sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V1_0_0_0
	cfg.Consumer.Return.Errors = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	cfg.Consumer.MaxWaitTime = time.Duration(viper.GetInt64(kafkaConsumerMaxWaitTimeViperKey)) * time.Millisecond

	cfg.Net.SASL.Enable = viper.GetBool(kafkaSASLEnabledViperKey)
	cfg.Net.SASL.Mechanism = sarama.SASLMechanism(viper.GetString(kafkaSASLMechanismViperKey))
	cfg.Net.SASL.Handshake = viper.GetBool(kafkaSASLHandshakeViperKey)
	cfg.Net.SASL.User = viper.GetString(kafkaSASLUserViperKey)
	cfg.Net.SASL.Password = viper.GetString(kafkaSASLPasswordViperKey)
	cfg.Net.SASL.SCRAMAuthzID = viper.GetString(kafkaSASLSCRAMAuthzIDViperKey)

	cfg.Net.TLS.Enable = viper.GetBool(kafkaTLSEnableViperKey)
	if cfg.Net.TLS.Enable {
		tlsConfig, err := NewTLSConfig(
			viper.GetString(kafkaTLSClientCertFilePathViperKey),
			viper.GetString(kafkaTLSClientKeyFilePathViperKey),
			viper.GetString(kafkaTLSCACertFilePathViperKey),
		)
		// Fatal if get error from NewTLSConfig
		if err != nil {
			return nil, err
		}
		cfg.Net.TLS.Config = tlsConfig
	}

	return cfg, nil
}
