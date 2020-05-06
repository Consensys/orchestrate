package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ChainRegistryURLViperKey, chainRegistryURLDefault)
	_ = viper.BindEnv(ChainRegistryURLViperKey, chainRegistryURLEnv)
	viper.SetDefault(ChainRegistryMetricsURLViperKey, chainRegistryMetricsURLDefault)
	_ = viper.BindEnv(ChainRegistryMetricsURLViperKey, chainRegistryMetricsURLEnv)
}

const (
	chainRegistryURLFlag     = "chain-registry-url"
	ChainRegistryURLViperKey = "chain.registry.url"
	chainRegistryURLDefault  = "localhost:8081"
	chainRegistryURLEnv      = "CHAIN_REGISTRY_URL"
)

// ChainRegistryURL register flag for the URL of the Chain Registry
func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Chain Registry HTTP endpoint. 
Environment variable: %q`, chainRegistryURLEnv)
	f.String(chainRegistryURLFlag, chainRegistryURLDefault, desc)
	_ = viper.BindPFlag(ChainRegistryURLViperKey, f.Lookup(chainRegistryURLFlag))
}

const (
	chainRegistryMetricsURLFlag     = "chain-registry-metrics-url"
	ChainRegistryMetricsURLViperKey = "chain.registry.metrics.url"
	chainRegistryMetricsURLDefault  = "localhost:8082"
	chainRegistryMetricsURLEnv      = "CHAIN_REGISTRY_METRICS_URL"
)

func MetricsURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Chain Registry Metrics endpoint.
Environment variable: %q`, chainRegistryMetricsURLEnv)
	f.String(chainRegistryMetricsURLFlag, chainRegistryMetricsURLDefault, desc)
	_ = viper.BindPFlag(ChainRegistryMetricsURLViperKey, f.Lookup(chainRegistryMetricsURLFlag))
}

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
		URL:        vipr.GetString(ChainRegistryURLViperKey),
		MetricsURL: vipr.GetString(ChainRegistryMetricsURLViperKey),
	}
}
