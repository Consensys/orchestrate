package client

import (
	"fmt"

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
	urlFlag     = "chain-registry-url"
	URLViperKey = "chain.registry.url"
	urlDefault  = "localhost:8081"
	urlEnv      = "CHAIN_REGISTRY_URL"
)

// ChainRegistryURL register flag for the URL of the Chain Registry
func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Chain Registry HTTP endpoint. 
Environment variable: %q`, urlEnv)
	f.String(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(URLViperKey, f.Lookup(urlFlag))
}

const (
	MetricsURLViperKey = "chain.registry.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "CHAIN_REGISTRY_METRICS_URL"
)

func Flags(f *pflag.FlagSet) {
	URL(f)
}

type Config struct {
	URL        string
	MetricsURL string
}

func NewConfig(url string) *Config {
	return &Config{
		URL: url,
	}
}

func NewConfigFromViper(vipr *viper.Viper) *Config {
	return &Config{
		URL:        vipr.GetString(URLViperKey),
		MetricsURL: vipr.GetString(MetricsURLViperKey),
	}
}
