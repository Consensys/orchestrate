package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(grpcTargetContractRegistryViperKey, grpcTargetContractRegistryDefault)
	_ = viper.BindEnv(grpcTargetContractRegistryViperKey, grpcTargetContractRegistryEnv)
}

var (
	grpcTargetContractRegistryFlag     = "grpc-target-contract-registry"
	grpcTargetContractRegistryViperKey = "grpc.target.contract.registry"
	grpcTargetContractRegistryDefault  = "localhost:8080"
	grpcTargetContractRegistryEnv      = "GRPC_TARGET_CONTRACT_REGISTRY"
)

// ContractRegistryGRPCTarget register flag for Ethereum client URLs
func ContractRegistryGRPCTarget(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`GRPC target Contract Registry (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, grpcTargetContractRegistryEnv)
	f.String(grpcTargetContractRegistryFlag, grpcTargetContractRegistryDefault, desc)
	viper.SetDefault(grpcTargetContractRegistryViperKey, grpcTargetContractRegistryDefault)
	_ = viper.BindPFlag(grpcTargetContractRegistryViperKey, f.Lookup(grpcTargetContractRegistryFlag))
	_ = viper.BindEnv(grpcTargetContractRegistryViperKey, grpcTargetContractRegistryEnv)
}
