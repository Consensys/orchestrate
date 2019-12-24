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
}

func Flags(f *pflag.FlagSet) {
	ProviderRefreshInterval(f)
}

const (
	providerRefreshIntervalFlag     = "provider-refresh-interval"
	ProviderRefreshIntervalViperKey = "provider.refresh-interval"
	providerRefreshIntervalDefault  = 5 * time.Second
	providerRefreshIntervalEnv      = "PROVIDER_REFRESH_INTERVAL"
)

// ProviderRefreshInterval register flag for refresh interval duration
func ProviderRefreshInterval(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Time interval for refreshing the configuration from the chain-registry
Environment variable: %q`, providerRefreshIntervalEnv)
	f.Duration(providerRefreshIntervalFlag, providerRefreshIntervalDefault, desc)
	_ = viper.BindPFlag(ProviderRefreshIntervalViperKey, f.Lookup(providerRefreshIntervalFlag))
}
