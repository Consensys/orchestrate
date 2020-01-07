package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ChainRegistryURLViperKey, chainRegistryURLDefault)
	_ = viper.BindEnv(ChainRegistryURLViperKey, chainRegistryURLEnv)
	viper.SetDefault(ChainProxyURLViperKey, chainProxyURLDefault)
	_ = viper.BindEnv(ChainProxyURLViperKey, chainProxyURLEnv)
}

func Flags(f *pflag.FlagSet) {
	ChainRegistryURL(f)
	ChainProxyURL(f)
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

const (
	chainProxyURLFlag     = "chain-proxy-url"
	ChainProxyURLViperKey = "chain.proxy.url"
	chainProxyURLDefault  = "localhost:8080"
	chainProxyURLEnv      = "CHAIN_PROXY_URL"
)

// ProviderRefreshInterval register flag for refresh interval duration
func ChainProxyURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Chain proxy
Environment variable: %q`, chainProxyURLEnv)
	f.String(chainProxyURLFlag, chainProxyURLDefault, desc)
	_ = viper.BindPFlag(ChainProxyURLViperKey, f.Lookup(chainProxyURLFlag))
}
