package http

import (
	"fmt"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
)

func init() {
	viper.SetDefault(HostnameViperKey, hostnameDefault)
	_ = viper.BindEnv(HostnameViperKey, hostnameEnv)
	viper.SetDefault(HTTPPortViperKey, httpPortDefault)
	_ = viper.BindEnv(HTTPPortViperKey, httpPortEnv)

	viper.SetDefault(MetricsHostnameViperKey, metricsHostnameDefault)
	_ = viper.BindEnv(MetricsHostnameViperKey, metricsHostnameEnv)
	viper.SetDefault(MetricsPortViperKey, metricsPortDefault)
	_ = viper.BindEnv(MetricsPortViperKey, metricsPortEnv)
}

const (
	hostnameFlag     = "rest-hostname"
	HostnameViperKey = "rest.hostname"
	hostnameDefault  = ""
	hostnameEnv      = "REST_HOSTNAME"
)

// Hostname register a flag for HTTP server address
func Hostname(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hostname to expose REST services
Environment variable: %q`, hostnameEnv)
	f.String(hostnameFlag, hostnameDefault, desc)
	_ = viper.BindPFlag(HostnameViperKey, f.Lookup(hostnameFlag))
}

const (
	httpPortFlag     = "rest-port"
	HTTPPortViperKey = "rest.port"
	httpPortDefault  = uint(8081)
	httpPortEnv      = "REST_PORT"
)

// Port register a flag for HTTp server port
func Port(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose REST services
Environment variable: %q`, httpPortEnv)
	f.Uint(httpPortFlag, httpPortDefault, desc)
	_ = viper.BindPFlag(HTTPPortViperKey, f.Lookup(httpPortFlag))
}

const (
	metricsHostnameFlag     = "metrics-hostname"
	MetricsHostnameViperKey = "metrics.hostname"
	metricsHostnameDefault  = ""
	metricsHostnameEnv      = "METRICS_HOSTNAME"
)

// MetricsHostname register a flag for metrics server hostname
func MetricsHostname(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hostname to expose metrics services
Environment variable: %q`, metricsHostnameEnv)
	f.String(metricsHostnameFlag, metricsHostnameDefault, desc)
	_ = viper.BindPFlag(MetricsHostnameViperKey, f.Lookup(metricsHostnameFlag))
}

const (
	metricsPortFlag     = "metrics-port"
	MetricsPortViperKey = "metrics.port"
	metricsPortDefault  = uint(8082)
	metricsPortEnv      = "METRICS_PORT"
)

// MetricsPort register a flag for metrics server port
func MetricsPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose metrics services
Environment variable: %q`, metricsPortEnv)
	f.Uint(metricsPortFlag, metricsPortDefault, desc)
	_ = viper.BindPFlag(MetricsPortViperKey, f.Lookup(metricsPortFlag))
}

func Flags(f *pflag.FlagSet) {
	Hostname(f)
	Port(f)
	MetricsHostname(f)
	MetricsPort(f)
}

func URL(hostname string, port uint) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}

func NewEPsConfig(vipr *viper.Viper) traefikstatic.EntryPoints {
	httpEp := &traefikstatic.EntryPoint{
		Address: URL(vipr.GetString(HostnameViperKey), vipr.GetUint(HTTPPortViperKey)),
	}
	httpEp.SetDefaults()

	metricsEp := &traefikstatic.EntryPoint{
		Address: URL(vipr.GetString(MetricsHostnameViperKey), vipr.GetUint(MetricsPortViperKey)),
	}
	metricsEp.SetDefaults()

	return traefikstatic.EntryPoints{
		DefaultHTTPEntryPoint:    httpEp,
		DefaultMetricsEntryPoint: metricsEp,
	}
}

func NewConfig(vipr *viper.Viper) *traefikstatic.Configuration {
	return &traefikstatic.Configuration{
		EntryPoints: NewEPsConfig(vipr),
		Metrics: &traefiktypes.Metrics{
			Prometheus: &traefiktypes.Prometheus{
				EntryPoint:           DefaultMetricsEntryPoint,
				Buckets:              []float64{0.1, 0.3, 1.2, 5},
				AddEntryPointsLabels: true,
				AddServicesLabels:    true,
			},
		},
		API: &traefikstatic.API{},
		ServersTransport: &traefikstatic.ServersTransport{
			MaxIdleConnsPerHost: 200,
			InsecureSkipVerify:  true,
		},
		Log: &traefiktypes.TraefikLog{
			Level:  vipr.GetString(logger.LogLevelViperKey),
			Format: viperToTraefikLogFormat(vipr.GetString(logger.LogFormatViperKey)),
		},
		AccessLog: &traefiktypes.AccessLog{
			Filters: &traefiktypes.AccessLogFilters{
				StatusCodes: []string{"100-199", "400-428", "430-599"},
			},
			Format: viperToTraefikLogFormat(vipr.GetString(logger.LogFormatViperKey)),
		},
	}
}

func DefaultConfig() *traefikstatic.Configuration {
	return NewConfig(viper.New())
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
