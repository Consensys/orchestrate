package contractregistry

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/config/static"
	grpcmetrics2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	httpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/metrics"
	metricregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/store/multi"
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
	multi.Flags(f)
	http.Flags(f)
	grpc.Flags(f)
	metricregistry.Flags(f, httpmetrics.ModuleName, grpcmetrics2.ModuleName, tcpmetrics.ModuleName)
}

type Config struct {
	App          *app.Config
	Store        *multi.Config
	ABIs         []string // Chains defined in ENV
	Multitenancy bool
}

func NewConfig(vipr *viper.Viper) *Config {
	cfg := &Config{
		App:          app.NewConfig(vipr),
		Store:        multi.NewConfig(vipr),
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
