package contractregistry

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
)

func init() {
	_ = viper.BindEnv(ABIViperKey, abiEnv)
	viper.SetDefault(ABIViperKey, abiDefault)
}

const (
	abiFlag     = "abi"
	ABIViperKey = "abis"
	abiEnv      = "ABI"
)

var abiDefault []string

// bindABIFlag register flag for ABI
func ABIs(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Smart Contract ABIs to register for crafting (expected format %v)
Environment variable: %q`, `<contract>:<abi>:<bytecode>:<deployedBytecode>`, abiEnv)
	f.StringSlice(abiFlag, abiDefault, desc)
	_ = viper.BindPFlag(ABIViperKey, f.Lookup(abiFlag))
}

// Flags register flags for Postgres database
func Flags(f *pflag.FlagSet) {
	ABIs(f)
	store.Flags(f)
	http.Flags(f)
	grpc.Flags(f)
}

type Config struct {
	App          *app.Config
	Store        *store.Config
	ABIs         []string // Chains defined in ENV
	Multitenancy bool
}

func NewConfig(vipr *viper.Viper) *Config {
	cfg := &Config{
		App:          app.NewConfig(vipr),
		Store:        store.NewConfig(vipr),
		ABIs:         viper.GetStringSlice(ABIViperKey),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
	}

	// Create static configuration for GRPC server
	cfg.App.GRPC.Static = &grpcstatic.Configuration{
		Services: &grpcstatic.Services{
			Contracts: &grpcstatic.Contracts{},
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
