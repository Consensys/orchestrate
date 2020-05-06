package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
)

func init() {
	viper.SetDefault(EnvelopeStoreURLViperKey, envelopeStoreURLDefault)
	_ = viper.BindEnv(EnvelopeStoreURLViperKey, envelopeStoreURLEnv)
	viper.SetDefault(EnvelopeStoreMetricsURLViperKey, envelopeStoreMetricsURLDefault)
	_ = viper.BindEnv(EnvelopeStoreMetricsURLViperKey, envelopeStoreMetricsURLEnv)
	viper.SetDefault(EnvelopeStoreHTTPURLViperKey, envelopeStoreHTTPURLDefault)
	_ = viper.BindEnv(EnvelopeStoreHTTPURLViperKey, envelopeStoreHTTPURLEnv)
}

const (
	envelopeStoreURLFlag     = "envelope-store-url"
	EnvelopeStoreURLViperKey = "envelope.store.url"
	envelopeStoreURLDefault  = "localhost:8080"
	envelopeStoreURLEnv      = "ENVELOPE_STORE_URL"
)

func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (GRPC target) Envelope Store (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, envelopeStoreURLEnv)
	f.String(envelopeStoreURLFlag, envelopeStoreURLDefault, desc)
	_ = viper.BindPFlag(EnvelopeStoreURLViperKey, f.Lookup(envelopeStoreURLFlag))
}

const (
	envelopeStoreMetricsURLFlag     = "envelope-store-metrics-url"
	EnvelopeStoreMetricsURLViperKey = "envelope.store.metrics.url"
	envelopeStoreMetricsURLDefault  = "localhost:8082"
	envelopeStoreMetricsURLEnv      = "ENVELOPE_STORE_METRICS_URL"
)

func MetricsURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of Envelope Store metrics endpoint
Environment variable: %q`, envelopeStoreMetricsURLEnv)
	f.String(envelopeStoreMetricsURLFlag, envelopeStoreMetricsURLDefault, desc)
	_ = viper.BindPFlag(EnvelopeStoreMetricsURLViperKey, f.Lookup(envelopeStoreMetricsURLFlag))
}

const (
	envelopeStoreHTTPURLFlag     = "envelope-store-http-url"
	EnvelopeStoreHTTPURLViperKey = "envelope.store.http.url"
	envelopeStoreHTTPURLDefault  = "localhost:8081"
	envelopeStoreHTTPURLEnv      = "ENVELOPE_STORE_HTTP_URL"
)

// EnvelopeStoreURL register flag for Ethereum client URLs
func HTTPURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Envelope Store HTTP endpoint
Environment variable: %q`, envelopeStoreHTTPURLEnv)
	f.String(envelopeStoreHTTPURLFlag, envelopeStoreHTTPURLDefault, desc)
	_ = viper.BindPFlag(EnvelopeStoreHTTPURLViperKey, f.Lookup(envelopeStoreHTTPURLFlag))
}

func Flags(f *pflag.FlagSet) {
	URL(f)
}

type Config struct {
	URL         string
	MetricsURL  string
	ServiceName string
}

func NewConfig(serviceName, url string) *Config {
	return &Config{
		URL:         url,
		ServiceName: serviceName,
	}
}

func NewConfigFromViper(vipr *viper.Viper) *Config {
	return &Config{
		URL:         vipr.GetString(EnvelopeStoreURLViperKey),
		ServiceName: vipr.GetString(jaeger.ServiceNameViperKey),
	}
}
