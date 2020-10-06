package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(GRPCURLViperKey, grpcURLDefault)
	_ = viper.BindEnv(GRPCURLViperKey, grpcURLEnv)
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
	viper.SetDefault(HTTPURLViperKey, httpURLDefault)
	_ = viper.BindEnv(HTTPURLViperKey, httpURLEnv)
}

const (
	grpcURLFlag     = "contract-registry-url"
	GRPCURLViperKey = "contract.registry.url"
	grpcURLDefault  = "localhost:8080"
	grpcURLEnv      = "CONTRACT_REGISTRY_URL"
)

// ContractRegistryURL register flag for Ethereum client URLs
func ContractRegistryURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (GRPC target) of the Contract Registry (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, grpcURLEnv)
	f.String(grpcURLFlag, grpcURLDefault, desc)
	_ = viper.BindPFlag(GRPCURLViperKey, f.Lookup(grpcURLFlag))
}

const (
	MetricsURLViperKey = "contract.registry.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "CONTRACT_REGISTRY_METRICS_URL"
)

const (
	HTTPURLViperKey = "contract.registry.http.url"
	httpURLDefault  = "localhost:8081"
	httpURLEnv      = "CONTRACT_REGISTRY_HTTP_URL"
)

type Config struct {
	URL        string
	MetricsURL string
}

func NewConfig(url string) *Config {
	return &Config{
		URL: url,
	}
}
