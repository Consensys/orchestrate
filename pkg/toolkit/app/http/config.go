package http

import (
	"fmt"
	"time"

	"github.com/consensys/orchestrate/pkg/multitenancy"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/key"
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(hostnameViperKey, hostnameDefault)
	_ = viper.BindEnv(hostnameViperKey, hostnameEnv)

	viper.SetDefault(httpPortViperKey, httpPortDefault)
	_ = viper.BindEnv(httpPortViperKey, httpPortEnv)

	viper.SetDefault(metricsHostnameViperKey, metricsHostnameDefault)
	_ = viper.BindEnv(metricsHostnameViperKey, metricsHostnameEnv)

	viper.SetDefault(metricsPortViperKey, metricsPortDefault)
	_ = viper.BindEnv(metricsPortViperKey, metricsPortEnv)
}

const (
	hostnameFlag     = "rest-hostname"
	hostnameViperKey = "rest.hostname"
	hostnameDefault  = ""
	hostnameEnv      = "REST_HOSTNAME"
)

// Hostname register a flag for HTTP server address
func hostname(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hostname to expose REST services
Environment variable: %q`, hostnameEnv)
	f.String(hostnameFlag, hostnameDefault, desc)
	_ = viper.BindPFlag(hostnameViperKey, f.Lookup(hostnameFlag))
}

const (
	httpPortFlag     = "rest-port"
	httpPortViperKey = "rest.port"
	httpPortDefault  = 8081
	httpPortEnv      = "REST_PORT"
)

// Port register a flag for HTTp server port
func port(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose REST services
Environment variable: %q`, httpPortEnv)
	f.Uint(httpPortFlag, httpPortDefault, desc)
	_ = viper.BindPFlag(httpPortViperKey, f.Lookup(httpPortFlag))
}

const (
	metricsHostnameFlag     = "metrics-hostname"
	metricsHostnameViperKey = "metrics.hostname"
	metricsHostnameDefault  = ""
	metricsHostnameEnv      = "METRICS_HOSTNAME"
)

// metricsHostname register a flag for metrics server hostname
func metricsHostname(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hostname to expose metrics services
Environment variable: %q`, metricsHostnameEnv)
	f.String(metricsHostnameFlag, metricsHostnameDefault, desc)
	_ = viper.BindPFlag(metricsHostnameViperKey, f.Lookup(metricsHostnameFlag))
}

const (
	metricsPortFlag     = "metrics-port"
	metricsPortViperKey = "metrics.port"
	metricsPortDefault  = 8082
	metricsPortEnv      = "METRICS_PORT"
)

// MetricsPort register a flag for metrics server port
func metricsPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose metrics services
Environment variable: %q`, metricsPortEnv)
	f.Uint(metricsPortFlag, metricsPortDefault, desc)
	_ = viper.BindPFlag(metricsPortViperKey, f.Lookup(metricsPortFlag))
}

func Flags(f *pflag.FlagSet) {
	hostname(f)
	port(f)
}

func MetricFlags(f *pflag.FlagSet) {
	metricsHostname(f)
	metricsPort(f)
}

func url(hostname string, port uint) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}

func NewEPsConfig(vipr *viper.Viper) traefikstatic.EntryPoints {
	httpEp := &traefikstatic.EntryPoint{
		Address: url(vipr.GetString(hostnameViperKey), vipr.GetUint(httpPortViperKey)),
	}
	httpEp.SetDefaults()

	metricsEp := &traefikstatic.EntryPoint{
		Address: url(vipr.GetString(metricsHostnameViperKey), vipr.GetUint(metricsPortViperKey)),
	}
	metricsEp.SetDefaults()

	return traefikstatic.EntryPoints{
		DefaultHTTPAppEntryPoint: httpEp,
		DefaultMetricsEntryPoint: metricsEp,
	}
}

type Config struct {
	Timeout               time.Duration
	KeepAlive             time.Duration
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	APIKey                string
	MaxIdleConnsPerHost   int
	MultiTenancy          bool
	AuthHeaderForward     bool
}

func NewConfig(vipr *viper.Viper) *Config {
	cfg := NewDefaultConfig()
	if vipr != nil {
		cfg.MultiTenancy = vipr.GetBool(multitenancy.EnabledViperKey)
		cfg.APIKey = vipr.GetString(key.APIKeyViperKey)
	}

	return cfg
}

func NewDefaultConfig() *Config {
	return &Config{
		MaxIdleConnsPerHost:   200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Timeout:               30 * time.Second,
		KeepAlive:             30 * time.Second,
		AuthHeaderForward:     true,
	}
}
