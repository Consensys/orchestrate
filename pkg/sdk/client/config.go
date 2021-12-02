package client

import (
	"fmt"
	"time"

	"github.com/consensys/orchestrate/pkg/backoff"
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
	urlFlag     = "api-url"
	URLViperKey = "api.url"
	urlDefault  = "http://localhost:8081"
	urlEnv      = "API_URL"
)

const (
	MetricsURLViperKey = "api.metrics.url"
	metricsURLDefault  = "http://localhost:8082"
	metricsURLEnv      = "API_METRICS_URL"
)

var defaultClientBackOff = backoff.ConstantBackOffWithMaxRetries(time.Second, 0)

func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the API HTTP endpoint. Environment variable: %q`, urlEnv)
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
