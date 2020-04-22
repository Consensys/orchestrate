package utils

import (
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(RetryInitialIntervalViperKey, 500*time.Millisecond)
	viper.SetDefault(RetryRandomFactorViperKey, 0.5)
	viper.SetDefault(RetryMultiplierViperKey, 1.5)
	viper.SetDefault(RetryMaxIntervalViperKey, 30*time.Second)
	viper.SetDefault(maxElapsedTimeViperKey, maxElapsedTimeDefault)
}

const (
	RetryInitialIntervalViperKey = "backOff.initinterval"
	RetryRandomFactorViperKey    = "backOff.randomfactor"
	RetryMultiplierViperKey      = "backOff.multiplier"
	RetryMaxIntervalViperKey     = "backOff.maxinterval"
)

const (
	maxElapsedTimeViperKey = "backOff.maxelapsedtime"
	maxElapsedTimeDefault  = 1 * time.Hour
)

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
func NewRetryConfig(vipr *viper.Viper) *RetryConfig {
	config := &RetryConfig{}
	config.InitialInterval = vipr.GetDuration(RetryInitialIntervalViperKey)
	config.RandomizationFactor = vipr.GetFloat64(RetryRandomFactorViperKey)
	config.Multiplier = vipr.GetFloat64(RetryMultiplierViperKey)
	config.MaxInterval = vipr.GetDuration(RetryMaxIntervalViperKey)
	config.MaxElapsedTime = vipr.GetDuration(maxElapsedTimeViperKey)

	return config
}

type Config struct {
	Retry *RetryConfig
}

// NewConfig creates a new default config
func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		Retry: NewRetryConfig(vipr),
	}
}

// NewBackOff creates a new Exponential backoff
func NewBackOff(conf *Config) backoff.BackOff {
	return &backoff.ExponentialBackOff{
		InitialInterval:     conf.Retry.InitialInterval,
		RandomizationFactor: conf.Retry.RandomizationFactor,
		Multiplier:          conf.Retry.Multiplier,
		MaxInterval:         conf.Retry.MaxInterval,
		MaxElapsedTime:      conf.Retry.MaxElapsedTime,
		Clock:               backoff.SystemClock,
		Stop:                backoff.Stop,
	}
}
