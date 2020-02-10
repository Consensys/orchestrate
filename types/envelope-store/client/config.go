package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(EnvelopeStoreURLViperKey, envelopeStoreURLDefault)
	_ = viper.BindEnv(EnvelopeStoreURLViperKey, envelopeStoreURLEnv)
}

const (
	envelopeStoreURLFlag     = "envelope-store-url"
	EnvelopeStoreURLViperKey = "envelope.store.url"
	envelopeStoreURLDefault  = "localhost:8080"
	envelopeStoreURLEnv      = "ENVELOPE_STORE_URL"
)

// EnvelopeStoreURL register flag for Ethereum client URLs
func EnvelopeStoreURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (GRPC target) Builder Store (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, envelopeStoreURLEnv)
	f.String(envelopeStoreURLFlag, envelopeStoreURLDefault, desc)
	viper.SetDefault(EnvelopeStoreURLViperKey, envelopeStoreURLDefault)
	_ = viper.BindPFlag(EnvelopeStoreURLViperKey, f.Lookup(envelopeStoreURLFlag))
	_ = viper.BindEnv(EnvelopeStoreURLViperKey, envelopeStoreURLEnv)
}
