package envelopestore

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
)

const EnvelopeStoreEnabledKey = "ENVELOPE_STORE_ENABLED"

func Flags(f *pflag.FlagSet) {
	store.Flags(f)
	http.Flags(f)
	grpc.Flags(f)
}

type Config struct {
	App          *app.Config
	Store        *store.Config
	Multitenancy bool
}

func NewConfig(vipr *viper.Viper) *Config {
	cfg := &Config{
		App:          app.NewConfig(vipr),
		Store:        store.NewConfig(vipr),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
	}

	// Create static configuration for GRPC server
	cfg.App.GRPC.Static = &grpcstatic.Configuration{
		Services: &grpcstatic.Services{
			Envelopes: &grpcstatic.Envelopes{},
		},
		Interceptors: []*grpcstatic.Interceptor{
			{Tags: &grpcstatic.Tags{}},
			{Logrus: &grpcstatic.Logrus{}},
			{Auth: &grpcstatic.Auth{}},
			{Error: &grpcstatic.Error{}},
			{Recovery: &grpcstatic.Recovery{}},
		},
	}

	return cfg
}
