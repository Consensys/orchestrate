package chainregistry

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ProviderRefreshIntervalViperKey, providerRefreshIntervalDefault)
	_ = viper.BindEnv(ProviderRefreshIntervalViperKey, providerRefreshIntervalEnv)
	viper.SetDefault(ChainProxyURLViperKey, chainProxyURLDefault)
	_ = viper.BindEnv(ChainProxyURLViperKey, chainProxyURLEnv)
}

func Flags(f *pflag.FlagSet) {
	ProviderRefreshInterval(f)
	ChainProxyURL(f)
}

const (
	providerRefreshIntervalFlag     = "tx-listener-provider-refresh-interval"
	ProviderRefreshIntervalViperKey = "tx-listener-provider.refresh-interval"
	providerRefreshIntervalDefault  = 5 * time.Second
	providerRefreshIntervalEnv      = "TX_LISTENER_PROVIDER_REFRESH_INTERVAL"
)

// ProviderRefreshInterval register flag for refresh interval duration
func ProviderRefreshInterval(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Time interval for refreshing the configuration from the chain-registry
Environment variable: %q`, providerRefreshIntervalEnv)
	f.Duration(providerRefreshIntervalFlag, providerRefreshIntervalDefault, desc)
	_ = viper.BindPFlag(ProviderRefreshIntervalViperKey, f.Lookup(providerRefreshIntervalFlag))
}

const (
	chainProxyURLFlag     = "chain-proxy-url"
	ChainProxyURLViperKey = "chain.proxy.url"
	chainProxyURLDefault  = "localhost:8081"
	chainProxyURLEnv      = "CHAIN_PROXY_URL"
)

// ProviderRefreshInterval register flag for refresh interval duration
func ChainProxyURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Chain proxy
Environment variable: %q`, chainProxyURLEnv)
	f.String(chainProxyURLFlag, chainProxyURLDefault, desc)
	_ = viper.BindPFlag(ChainProxyURLViperKey, f.Lookup(chainProxyURLFlag))
}
