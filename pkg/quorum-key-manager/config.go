package quorumkeymanager

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(URLViperKey, urlDefault)
	_ = viper.BindEnv(URLViperKey, urlEnv)
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
	viper.SetDefault(StoreNameViperKey, StoreNameDefault)
	_ = viper.BindEnv(StoreNameViperKey, storeNameEnv)
}

const (
	urlFlag     = "key-manager-url"
	URLViperKey = "key.manager.url"
	urlDefault  = "http://localhost:8081"
	urlEnv      = "KEY_MANAGER_URL"
)

const (
	metricsURLFlag     = "key-manager-metrics-url"
	MetricsURLViperKey = "key.manager.metrics.url"
	metricsURLDefault  = "http://localhost:8082"
	metricsURLEnv      = "KEY_MANAGER_METRICS_URL"
)

const (
	storeNameFlag     = "key-manager-store-name"
	StoreNameViperKey = "key.manager.store.name"
	StoreNameDefault  = ""
	storeNameEnv      = "KEY_MANAGER_STORE_NAME"
)

func url(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Key Manager HTTP endpoint.
Environment variable: %q`, urlEnv)
	f.String(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(URLViperKey, f.Lookup(urlFlag))
}

func metricsURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Key Manager HTTP metrics endpoint.
Environment variable: %q`, metricsURLEnv)
	f.String(metricsURLFlag, metricsURLDefault, desc)
	_ = viper.BindPFlag(MetricsURLViperKey, f.Lookup(metricsURLFlag))
}

func storeName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager ethereum account store name.
Environment variable: %q`, storeNameEnv)
	f.String(storeNameFlag, metricsURLDefault, desc)
	_ = viper.BindPFlag(StoreNameViperKey, f.Lookup(storeNameFlag))
}

func Flags(f *pflag.FlagSet) {
	url(f)
	metricsURL(f)
	storeName(f)
}

type Config struct {
	URL        string
	MetricsURL string
	StoreName  string
}

func NewConfigFromViper(vipr *viper.Viper) *Config {
	return &Config{
		URL:        vipr.GetString(URLViperKey),
		MetricsURL: vipr.GetString(MetricsURLViperKey),
		StoreName:  vipr.GetString(StoreNameViperKey),
	}
}
