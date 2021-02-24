package sarama

import (
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
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

var rebalanceStrategy = map[string]sarama.BalanceStrategy{
	"Range":      sarama.BalanceStrategyRange,
	"RoundRobin": sarama.BalanceStrategyRoundRobin,
	"Sticky":     sarama.BalanceStrategySticky,
}

func NewSaramaConfig() (*sarama.Config, error) {
	cfg := sarama.NewConfig()

	// // If not able to parse version then use Min version by default
	if version, err := sarama.ParseKafkaVersion(viper.GetString(kafkaVersionViperKey)); err == nil {
		cfg.Version = version
	} else {
		cfg.Version = sarama.V1_0_0_0
	}

	cfg.ClientID = "sarama-orchestrate"
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.AutoCommit.Enable = false
	cfg.Consumer.MaxWaitTime = viper.GetDuration(kafkaConsumerMaxWaitTimeViperKey)
	cfg.Consumer.MaxProcessingTime = viper.GetDuration(kafkaConsumerMaxProcessingTimeViperKey)
	cfg.Consumer.Group.Session.Timeout = viper.GetDuration(kafkaConsumerGroupSessionTimeoutViperKey)
	cfg.Consumer.Group.Heartbeat.Interval = viper.GetDuration(kafkaConsumerGroupHeartbeatIntervalViperKey)
	cfg.Consumer.Group.Rebalance.Timeout = viper.GetDuration(kafkaConsumerGroupRebalanceTimeoutViperKey)
	if strategy, ok := rebalanceStrategy[viper.GetString(kafkaConsumerGroupRebalanceStrategyViperKey)]; ok {
		cfg.Consumer.Group.Rebalance.Strategy = strategy
	}
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
