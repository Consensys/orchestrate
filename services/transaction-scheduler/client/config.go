package client

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(URLViperKey, urlDefault)
	_ = viper.BindEnv(URLViperKey, urlEnv)
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
}

const (
	urlFlag     = "transaction-scheduler-url"
	URLViperKey = "transaction.scheduler.url"
	urlDefault  = "localhost:8081"
	urlEnv      = "TRANSACTION_SCHEDULER_URL"
)

const (
	MetricsURLViperKey = "transaction.scheduler.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "TRANSACTION_SCHEDULER_METRICS_URL"
)

var defaultClientBackOff = backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 0)

// ChainRegistryURL register flag for the URL of the Chain Registry
func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Transaction Scheduler HTTP endpoint. 
Environment variable: %q`, urlEnv)
	f.String(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(URLViperKey, f.Lookup(urlFlag))
}

func Flags(f *pflag.FlagSet) {
	URL(f)
}

type Config struct {
	URL        string
	MetricsURL string
	backOff    backoff.BackOff
}

func NewConfig(url string, backOff backoff.BackOff) *Config {
	if backOff == nil {
		backOff = defaultClientBackOff
	}

	return &Config{
		URL:        url,
		MetricsURL: metricsURLDefault,
		backOff:    backOff,
	}
}

func NewConfigFromViper(vipr *viper.Viper, backOff backoff.BackOff) *Config {
	if backOff == nil {
		backOff = defaultClientBackOff
	}

	return &Config{
		URL:        vipr.GetString(URLViperKey),
		MetricsURL: vipr.GetString(MetricsURLViperKey),
		backOff:    backOff,
	}
}
