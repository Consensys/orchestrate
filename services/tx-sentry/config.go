package txsentry

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	sentryRefreshIntervalFlag     = "tx-sentry-refresh-interval"
	sentryRefreshIntervalViperKey = "tx-sentry.refresh-interval"
	sentryRefreshIntervalDefault  = 5 * time.Second
	sentryRefreshIntervalEnv      = "TX_SENTRY_REFRESH_INTERVAL"
)

func init() {
	viper.SetDefault(sentryRefreshIntervalViperKey, sentryRefreshIntervalDefault)
	_ = viper.BindEnv(sentryRefreshIntervalViperKey, sentryRefreshIntervalEnv)
}

// Flags register flags for tx sentry
func Flags(f *pflag.FlagSet) {
	pendingDurationDesc := fmt.Sprintf(`Interval of time between checks for pending transactions. Environment variable: %q`, sentryRefreshIntervalEnv)
	f.Duration(sentryRefreshIntervalFlag, sentryRefreshIntervalDefault, pendingDurationDesc)
	_ = viper.BindPFlag(sentryRefreshIntervalViperKey, f.Lookup(sentryRefreshIntervalFlag))
}

type Config struct {
	RefreshInterval time.Duration
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		RefreshInterval: vipr.GetDuration(sentryRefreshIntervalViperKey),
	}
}
