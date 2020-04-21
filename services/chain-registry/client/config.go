package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ChainRegistryURLViperKey, chainRegistryURLDefault)
	_ = viper.BindEnv(ChainRegistryURLViperKey, chainRegistryURLEnv)
}

const (
	chainRegistryURLFlag     = "chain-registry-url"
	ChainRegistryURLViperKey = "chain.registry.url"
	chainRegistryURLDefault  = "localhost:8081"
	chainRegistryURLEnv      = "CHAIN_REGISTRY_URL"
)

// ChainRegistryURL register flag for the URL of the Chain Registry
func Flags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Chain Registry Environment variable: %q`, chainRegistryURLEnv)
	f.String(chainRegistryURLFlag, chainRegistryURLDefault, desc)
	_ = viper.BindPFlag(ChainRegistryURLViperKey, f.Lookup(chainRegistryURLFlag))
}

type Config struct {
	URL string
}

func NewConfig(url string) *Config {
	return &Config{
		URL: url,
	}
}

func NewConfigFromViper(vipr *viper.Viper) *Config {
	return NewConfig(vipr.GetString(ChainRegistryURLViperKey))
}
