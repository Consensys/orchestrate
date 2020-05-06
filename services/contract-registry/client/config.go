package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ContractRegistryURLViperKey, contractRegistryURLDefault)
	_ = viper.BindEnv(ContractRegistryURLViperKey, contractRegistryURLEnv)
	viper.SetDefault(ContractRegistryMetricsURLViperKey, contractRegistryMetricsURLDefault)
	_ = viper.BindEnv(ContractRegistryMetricsURLViperKey, contractRegistryMetricsURLEnv)
	viper.SetDefault(ContractRegistryHTTPURLViperKey, contractRegistryHTTPURLDefault)
	_ = viper.BindEnv(ContractRegistryHTTPURLViperKey, contractRegistryHTTPURLEnv)
}

const (
	contractRegistryURLFlag     = "contract-registry-url"
	ContractRegistryURLViperKey = "contract.registry.url"
	contractRegistryURLDefault  = "localhost:8080"
	contractRegistryURLEnv      = "CONTRACT_REGISTRY_URL"
)

// ContractRegistryURL register flag for Ethereum client URLs
func ContractRegistryURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (GRPC target) of the Contract Registry (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, contractRegistryURLEnv)
	f.String(contractRegistryURLFlag, contractRegistryURLDefault, desc)
	_ = viper.BindPFlag(ContractRegistryURLViperKey, f.Lookup(contractRegistryURLFlag))
}

const (
	contractRegistryMetricsURLFlag     = "contract-registry-metrics-url"
	ContractRegistryMetricsURLViperKey = "contract.registry.metrics.url"
	contractRegistryMetricsURLDefault  = "localhost:8082"
	contractRegistryMetricsURLEnv      = "CONTRACT_REGISTRY_METRICS_URL"
)

// ContractRegistryURL register flag for Ethereum client URLs
func MetricsURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Contract Registry metrics endpoint
Environment variable: %q`, contractRegistryMetricsURLEnv)
	f.String(contractRegistryMetricsURLFlag, contractRegistryMetricsURLDefault, desc)
	_ = viper.BindPFlag(ContractRegistryMetricsURLViperKey, f.Lookup(contractRegistryMetricsURLFlag))
}

const (
	contractRegistryHTTPURLFlag     = "contract-registry-http-url"
	ContractRegistryHTTPURLViperKey = "contract.registry.http.url"
	contractRegistryHTTPURLDefault  = "localhost:8081"
	contractRegistryHTTPURLEnv      = "CONTRACT_REGISTRY_HTTP_URL"
)

// ContractRegistryURL register flag for Ethereum client URLs
func HTTPURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Contract Registry HTTP endpoint
Environment variable: %q`, contractRegistryHTTPURLEnv)
	f.String(contractRegistryHTTPURLFlag, contractRegistryHTTPURLDefault, desc)
	_ = viper.BindPFlag(ContractRegistryHTTPURLViperKey, f.Lookup(contractRegistryHTTPURLFlag))
}

func Flags(f *pflag.FlagSet) {
	ContractRegistryURL(f)
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
		URL:        vipr.GetString(ContractRegistryURLViperKey),
		MetricsURL: vipr.GetString(ContractRegistryMetricsURLViperKey),
	}
}
