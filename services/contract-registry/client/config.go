package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ContractRegistryURLViperKey, contractRegistryURLDefault)
	_ = viper.BindEnv(ContractRegistryURLViperKey, contractRegistryURLEnv)
}

const (
	contractRegistryURLFlag     = "contract-registry-url"
	ContractRegistryURLViperKey = "contract.registry.url"
	contractRegistryURLDefault  = "localhost:8080"
	contractRegistryURLEnv      = "CONTRACT_REGISTRY_URL"
)

// ContractRegistryURL register flag for Ethereum client URLs
func ContractRegistryURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (GRPC target) of the Contract Registry (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, contractRegistryURLEnv)
	f.String(contractRegistryURLFlag, contractRegistryURLDefault, desc)
	viper.SetDefault(ContractRegistryURLViperKey, contractRegistryURLDefault)
	_ = viper.BindPFlag(ContractRegistryURLViperKey, f.Lookup(contractRegistryURLFlag))
	_ = viper.BindEnv(ContractRegistryURLViperKey, contractRegistryURLEnv)
}
