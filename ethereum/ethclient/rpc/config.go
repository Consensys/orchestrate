package rpc

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(urlViperKey, urlEnv)
	viper.SetDefault(urlViperKey, urlDefault)
	viper.SetDefault(RetryInitialIntervalViperKey, 500*time.Millisecond)
	viper.SetDefault(RetryRandomFactorViperKey, 0.5)
	viper.SetDefault(RetryMultiplierViperKey, 1.5)
	viper.SetDefault(RetryMaxIntervalViperKey, 30*time.Second)
	_ = viper.BindEnv(maxElapsedTimeViperKey, maxElapsedTimeEnv)
	viper.SetDefault(maxElapsedTimeViperKey, maxElapsedTimeDefault)
}

var (
	urlFlag     = "eth-client-url"
	urlViperKey = "eth.client.url"
	urlDefault  []string
	urlEnv      = "ETH_CLIENT_URL"
)

func Flags(f *pflag.FlagSet) {
	URLs(f)
	MaxElapsedTime(f)
}

// URLs register flag for Ethereum client urls
func URLs(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Ethereum client url
Environment variable: %q`, urlEnv)
	f.StringSlice(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(urlViperKey, f.Lookup(urlFlag))
}

const (
	RetryInitialIntervalViperKey = "ethclient.retry.initinterval"
	RetryRandomFactorViperKey    = "ethclient.retry.randomfactor"
	RetryMultiplierViperKey      = "ethclient.retry.multiplier"
	RetryMaxIntervalViperKey     = "ethclient.retry.maxinterval"
)

const (
	maxElapsedTimeFlag     = "ethclient-retry-maxelapsedtime"
	maxElapsedTimeViperKey = "ethclient.retry.maxelapsedtime"
	maxElapsedTimeDefault  = 1 * time.Hour
	maxElapsedTimeEnv      = "ETH_CLIENT_RETRY_MAX_ELAPSED_TIME"
)

// MaxElapsedTime register flag for maximum elapsed time to retry RPC calls on Ethereum clients
func MaxElapsedTime(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Max elapsed time to retry rpc calls on Ethereum clients
Environment variable: %q`, maxElapsedTimeEnv)
	f.Duration(maxElapsedTimeFlag, maxElapsedTimeDefault, desc)
	_ = viper.BindPFlag(maxElapsedTimeViperKey, f.Lookup(maxElapsedTimeFlag))
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
	config.InitialInterval = viper.GetDuration(RetryInitialIntervalViperKey)
	config.RandomizationFactor = viper.GetFloat64(RetryRandomFactorViperKey)
	config.Multiplier = viper.GetFloat64(RetryMultiplierViperKey)
	config.MaxInterval = viper.GetDuration(RetryMaxIntervalViperKey)
	config.MaxElapsedTime = viper.GetDuration(maxElapsedTimeViperKey)

	return config
}

type Config struct {
	Retry *RetryConfig
}

// NewConfig creates a new default config
func NewConfig() *Config {
	return &Config{
		Retry: NewRetryConfig(),
	}
}
