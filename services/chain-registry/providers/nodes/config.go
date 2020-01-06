package nodes

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ProviderRefreshIntervalViperKey, providerRefreshIntervalDefault)
	_ = viper.BindEnv(ProviderRefreshIntervalViperKey, providerRefreshIntervalEnv)
}

func Flags(f *pflag.FlagSet) {
	ProviderRefreshInterval(f)
}

const (
	providerRefreshIntervalFlag     = "chain-registry-provider-nodes-refresh-interval"
	ProviderRefreshIntervalViperKey = "chain-registry.provider.nodes.refresh-interval"
	providerRefreshIntervalDefault  = 1 * time.Second
	providerRefreshIntervalEnv      = "CHAIN_REGISTRY_PROVIDER_NODES_REFRESH_INTERVAL"
)

// ProviderRefreshInterval register flag for refresh interval duration
func ProviderRefreshInterval(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Time interval for refreshing the list of nodes from storage
Environment variable: %q`, providerRefreshIntervalEnv)
	f.Duration(providerRefreshIntervalFlag, providerRefreshIntervalDefault, desc)
	_ = viper.BindPFlag(ProviderRefreshIntervalViperKey, f.Lookup(providerRefreshIntervalFlag))
}
