package configwatcher

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ProvidersThrottleDurationViperKey, providersThrottleDurationDefault)
	_ = viper.BindEnv(ProvidersThrottleDurationViperKey, providersThrottleDurationEnv)
}

func Flags(f *pflag.FlagSet) {
	ProvidersThrottleDuration(f)
}

const (
	providersThrottleDurationFlag     = "providers-throttle-duration"
	ProvidersThrottleDurationViperKey = "providers.throttle.duration"
	providersThrottleDurationDefault  = time.Second
	providersThrottleDurationEnv      = "PROVIDERS_THROTTLE_DURATION"
)

// ProvidersThrottleDuration register flag for throttle time duration
func ProvidersThrottleDuration(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Duration to wait for, after a configuration reload, before taking into account any new configuration
Environment variable: %q`, providersThrottleDurationEnv)
	f.Duration(providersThrottleDurationFlag, providersThrottleDurationDefault, desc)
	_ = viper.BindPFlag(ProvidersThrottleDurationViperKey, f.Lookup(providersThrottleDurationFlag))
}

type Config struct {
	ProvidersThrottleDuration time.Duration
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		ProvidersThrottleDuration: vipr.GetDuration(ProvidersThrottleDurationViperKey),
	}
}
