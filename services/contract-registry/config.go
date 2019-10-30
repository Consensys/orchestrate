package contractregistry

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/abi"
)

func init() {
	_ = viper.BindEnv(typeViperKey, typeEnv)
	viper.SetDefault(typeViperKey, typeDefault)

	_ = viper.BindEnv(abiViperKey, abiEnv)
	viper.SetDefault(abiViperKey, abiDefault)
}

var (
	typeFlag     = "contract-registry-type"
	typeViperKey = "contract-registry.type"
	typeDefault  = "postgres"
	typeEnv      = "CONTRACT_REGISTRY_TYPE"
)

// Type register flag for the Contract Registry to select
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Contract Registry (one of %q)
Environment variable: %q`, []string{"mock", "postgres"}, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(typeViperKey, f.Lookup(typeFlag))
}

var (
	abiFlag     = "abi"
	abiViperKey = "abis"
	abiDefault  []string
	abiEnv      = "ABI"
)

// ABIs register flag for ABI
func ABIs(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Smart Contract ABIs to register for crafting (expected format %v)
Environment variable: %q`, `<contract>:<abi>:<bytecode>:<deployedBytecode>`, abiEnv)
	f.StringSlice(abiFlag, abiDefault, desc)
	_ = viper.BindPFlag(abiViperKey, f.Lookup(abiFlag))
}

// FromABIConfig read viper config and return contracts
func FromABIConfig() (contracts []*abi.Contract, err error) {
	for _, ABI := range viper.GetStringSlice(abiViperKey) {
		c, err := abi.StringToContract(ABI)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, c)
	}
	return contracts, nil
}
