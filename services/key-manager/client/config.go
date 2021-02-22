package client

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/backoff"
)

func init() {
	viper.SetDefault(URLViperKey, urlDefault)
	_ = viper.BindEnv(URLViperKey, urlEnv)
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
}

const (
	urlFlag     = "key-manager-url"
	URLViperKey = "key.manager.url"
	urlDefault  = "http://localhost:8081"
	urlEnv      = "KEY_MANAGER_URL"
)

const (
	metricsURLFlag     = "key-manager-metrics-url"
	MetricsURLViperKey = "key.manager.metrics.url"
	metricsURLDefault  = "http://localhost:8082"
	metricsURLEnv      = "KEY_MANAGER_METRICS_URL"
)

var defaultClientBackOff = backoff.ConstantBackOffWithMaxRetries(time.Second, 0)

func url(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Key Manager HTTP endpoint.
Environment variable: %q`, urlEnv)
	f.String(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(URLViperKey, f.Lookup(urlFlag))
}

func metricsURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Key Manager HTTP metrics endpoint.
Environment variable: %q`, metricsURLEnv)
	f.String(metricsURLFlag, metricsURLDefault, desc)
	_ = viper.BindPFlag(MetricsURLViperKey, f.Lookup(metricsURLDefault))
}

func Flags(f *pflag.FlagSet) {
	url(f)
	metricsURL(f)
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
