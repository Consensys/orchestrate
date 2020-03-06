package chainregistry

import (
	"fmt"
	"time"

	"github.com/containous/traefik/v2/pkg/config/static"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"
)

const (
	DefaultInternalEntryPointName = "orchestrate"
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

func NewConfig() *static.Configuration {
	orchestrateEp := &static.EntryPoint{
		Address: rest.URL(),
	}
	orchestrateEp.SetDefaults()

	metricsEp := &static.EntryPoint{
		Address: metrics.URL(),
	}
	metricsEp.SetDefaults()

	return &static.Configuration{
		Providers: &static.Providers{
			// TODO: make it configurable
			ProvidersThrottleDuration: traefiktypes.Duration(time.Second),
		},
		EntryPoints: static.EntryPoints{
			DefaultInternalEntryPointName: orchestrateEp,
			"metrics":                     metricsEp,
		},
		API: &static.API{
			Dashboard: true,
		},
		Metrics: &traefiktypes.Metrics{
			Prometheus: &traefiktypes.Prometheus{
				EntryPoint:           "metrics",
				Buckets:              []float64{0.1, 0.3, 1.2, 5},
				AddEntryPointsLabels: true,
				AddServicesLabels:    true,
			},
		},
		ServersTransport: &static.ServersTransport{
			MaxIdleConnsPerHost: 200,
			InsecureSkipVerify:  true,
		},
		Log: &traefiktypes.TraefikLog{
			Level:  viper.GetString(logger.LogLevelViperKey),
			Format: viperToTraefikLogFormat(viper.GetString(logger.LogFormatViperKey)),
		},
		AccessLog: &traefiktypes.AccessLog{
			Filters: &traefiktypes.AccessLogFilters{
				StatusCodes: []string{"100-199", "400-428", "430-599"},
			},
			Format: viperToTraefikLogFormat(viper.GetString(logger.LogFormatViperKey)),
		},
	}
}

func viperToTraefikLogFormat(format string) string {
	switch format {
	case "text":
		return "common"
	case "json":
		return "json"
	default:
		return "json"
	}
}
