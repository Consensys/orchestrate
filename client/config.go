package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(grpcTargetEnvelopeStoreViperKey, grpcTargetEnvelopeStoreDefault)
	_ = viper.BindEnv(grpcTargetEnvelopeStoreViperKey, grpcTargetEnvelopeStoreEnv)
}

var (
	grpcTargetEnvelopeStoreFlag     = "grpc-target-envelope-store"
	grpcTargetEnvelopeStoreViperKey = "grpc.target.envelope.store"
	grpcTargetEnvelopeStoreDefault  = "localhost:8080"
	grpcTargetEnvelopeStoreEnv      = "GRPC_TARGET_ENVELOPE_STORE"
)

// EnvelopeStoreGRPCTarget register flag for Ethereum client URLs
func EnvelopeStoreGRPCTarget(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`GRPC target Envelope Store (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, grpcTargetEnvelopeStoreEnv)
	f.String(grpcTargetEnvelopeStoreFlag, grpcTargetEnvelopeStoreDefault, desc)
	viper.SetDefault(grpcTargetEnvelopeStoreViperKey, grpcTargetEnvelopeStoreDefault)
	_ = viper.BindPFlag(grpcTargetEnvelopeStoreViperKey, f.Lookup(grpcTargetEnvelopeStoreFlag))
	_ = viper.BindEnv(grpcTargetEnvelopeStoreViperKey, grpcTargetEnvelopeStoreEnv)
}
