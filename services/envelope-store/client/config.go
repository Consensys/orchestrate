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
}

type Config struct {
	envelopeStoreURL string
	serviceName      string
}

const (
	envelopeStoreURLFlag     = "envelope-store-url"
	EnvelopeStoreURLViperKey = "envelope.store.url"
	envelopeStoreURLDefault  = "localhost:8080"
	envelopeStoreURLEnv      = "ENVELOPE_STORE_URL"
)

func NewConfig(serviceName, envelopeStoreURL string) Config {
	return Config{
		envelopeStoreURL,
		serviceName,
	}
}

func NewConfigFromViper(vipr *viper.Viper) Config {
	return NewConfig(
		vipr.GetString(jaeger.ServiceNameViperKey),
		vipr.GetString(EnvelopeStoreURLViperKey),
	)
}

func Flags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (GRPC target) Envelope Store (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, envelopeStoreURLEnv)
	f.String(envelopeStoreURLFlag, envelopeStoreURLDefault, desc)
	_ = viper.BindPFlag(EnvelopeStoreURLViperKey, f.Lookup(envelopeStoreURLFlag))
}
