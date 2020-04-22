package contractregistry

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
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
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Smart Contract ABIs to register for crafting (expected format %v)
Environment variable: %q`, `<contract>:<abi>:<bytecode>:<deployedBytecode>`, abiEnv)
	f.StringSlice(abiFlag, abiDefault, desc)
	_ = viper.BindPFlag(ABIViperKey, f.Lookup(abiFlag))
}

// Flags register flags for Postgres database
func Flags(f *pflag.FlagSet) {
	Type(f)
	store.Flags(f)
	http.Flags(f)
	grpc.Flags(f)
}
