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

func Flags(f *pflag.FlagSet) {
	ChainRegistryURL(f)
}

const (
	chainRegistryURLFlag     = "chain-registry-url"
	ChainRegistryURLViperKey = "chain.registry.url"
	chainRegistryURLDefault  = "localhost:8081"
	chainRegistryURLEnv      = "CHAIN_REGISTRY_URL"
)

// ChainRegistryURL register flag for the URL of the Chain Registry
func ChainRegistryURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Chain Registry
Environment variable: %q`, chainRegistryURLEnv)
	f.String(chainRegistryURLFlag, chainRegistryURLDefault, desc)
	viper.SetDefault(ChainRegistryURLViperKey, chainRegistryURLDefault)
	_ = viper.BindPFlag(ChainRegistryURLViperKey, f.Lookup(chainRegistryURLFlag))
}

type Config struct {
	URL string
}

func NewConfig() *Config {
	return &Config{
		URL: viper.GetString(ChainRegistryURLViperKey),
	}
}
