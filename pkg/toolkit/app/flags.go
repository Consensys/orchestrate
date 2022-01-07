package app

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(hostnameViperKey, hostnameDefault)
	_ = viper.BindEnv(hostnameViperKey, hostnameEnv)

	viper.SetDefault(httpPortViperKey, httpPortDefault)
	_ = viper.BindEnv(httpPortViperKey, httpPortEnv)

	viper.SetDefault(accessLogEnabledKey, accessLogEnabledDefault)
	_ = viper.BindEnv(accessLogEnabledKey, accessLogEnabledEnv)

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

const (
	accessLogEnabledFlag    = "accesslog-enabled"
	accessLogEnabledKey     = "accesslog.enabled"
	accessLogEnabledDefault = false
	accessLogEnabledEnv     = "ACCESSLOG_ENABLED"
)

func accessLogEnabled(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable http accesslog stout
Environment variable: %q`, accessLogEnabledEnv)
	f.Bool(accessLogEnabledFlag, accessLogEnabledDefault, desc)
	_ = viper.BindPFlag(accessLogEnabledKey, f.Lookup(accessLogEnabledFlag))
}

func Flags(f *pflag.FlagSet) {
	hostname(f)
	port(f)
	accessLogEnabled(f)
}

func MetricFlags(f *pflag.FlagSet) {
	metricsHostname(f)
	metricsPort(f)
}

func url(hostname string, port uint) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}
