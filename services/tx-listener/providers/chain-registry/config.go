package chainregistry

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
)

func init() {
	viper.SetDefault(ProviderRefreshIntervalViperKey, providerRefreshIntervalDefault)
	_ = viper.BindEnv(ProviderRefreshIntervalViperKey, providerRefreshIntervalEnv)
}

func Flags(f *pflag.FlagSet) {
	ProviderRefreshInterval(f)
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

type Config struct {
	RefreshInterval  time.Duration
	ChainRegistryURL string
}

func NewConfig() *Config {
	return &Config{
		ChainRegistryURL: viper.GetString(chainregistry.URLViperKey),
		RefreshInterval:  viper.GetDuration(ProviderRefreshIntervalViperKey),
	}
}
