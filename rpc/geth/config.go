package geth

import (
	"time"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("ethclient.retry.initinterval", 500*time.Millisecond)
	viper.SetDefault("ethclient.retry.randomfactor", 0.5)
	viper.SetDefault("ethclient.retry.multiplier", 1.5)
	viper.SetDefault("ethclient.retry.maxinterval", 10*time.Second)
	viper.SetDefault("ethclient.retry.maxelapsedtime", 60*time.Second)
}

// TODO: move everything related to Retry and Backoff to pkg
// RetryConfig is a configuration for Exponential Backoff
type RetryConfig struct {
	// We use an exponential backoff retry strategy when fetching from an Eth Client
	// See https://github.com/cenkalti/backoff/blob/master/exponential.go
	InitialInterval     time.Duration
	RandomizationFactor float64
	Multiplier          float64
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration
}

// NewRetryConfig creates a New Configuration for an Ex
func NewRetryConfig() *RetryConfig {
	config := &RetryConfig{}
	config.InitialInterval = viper.GetDuration("ethclient.initinterval")
	config.RandomizationFactor = viper.GetFloat64("ethclient.retry.randomfactor")
	config.Multiplier = viper.GetFloat64("ethclient.retry.multiplier")
	config.MaxInterval = viper.GetDuration("ethclient.retry.maxinterval")
	config.MaxElapsedTime = viper.GetDuration("ethclient.retry.maxelapsedtime")

	return config
}

type Config struct {
	Retry RetryConfig
}

// NewConfig creates a new default config
func NewConfig() *Config {
	return &Config{
		Retry: *NewRetryConfig(),
	}
}
