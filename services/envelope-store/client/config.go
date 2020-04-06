package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
)

type Config struct {
	envelopeStoreURL string
	serviceName string
}

const (
	envelopeStoreURLFlag     = "envelope-store-url"
	EnvelopeStoreURLViperKey = "envelope.store.url"
	envelopeStoreURLDefault  = "localhost:8080"
	envelopeStoreURLEnv      = "ENVELOPE_STORE_URL"
)

func NewConfig(serviceName string, envelopeStoreURL string) Config {
	return Config{
		envelopeStoreURL,
		serviceName,
	}
}

func NewConfigFromViper(vipr *viper.Viper) Config {
	return NewConfig(
		vipr.GetString(EnvelopeStoreURLViperKey),
		viper.GetString(jaeger.ServiceNameViperKey),
	)
}

func Flags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (GRPC target) Envelope Store (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, envelopeStoreURLEnv)
	f.String(envelopeStoreURLFlag, envelopeStoreURLDefault, desc)
	viper.SetDefault(EnvelopeStoreURLViperKey, envelopeStoreURLDefault)
	_ = viper.BindPFlag(EnvelopeStoreURLViperKey, f.Lookup(envelopeStoreURLFlag))
	_ = viper.BindEnv(EnvelopeStoreURLViperKey, envelopeStoreURLEnv)
}