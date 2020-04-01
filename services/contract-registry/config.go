package contractregistry

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/controllers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

const (
	typeFlag     = "contract-registry-type"
	TypeViperKey = "contract-registry.type"
	typeDefault  = postgresOpt
	typeEnv      = "CONTRACT_REGISTRY_TYPE"
	abiFlag      = "abi"
	abiViperKey  = "abis"
	abiEnv       = "ABI"
	postgresOpt  = "postgres"
)

var abiDefault []string

// bindTypeFlag register flag for the Contract Registry to select
func bindTypeFlag(flag *pflag.FlagSet) {
	description := fmt.Sprintf(`Type of Contract Registry (one of %q) Environment variable: %q`, []string{postgresOpt}, typeEnv)
	flag.String(typeFlag, typeDefault, description)
	_ = viper.BindPFlag(TypeViperKey, flag.Lookup(typeFlag))
}

// bindABIFlag register flag for ABI
func bindABIFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Smart Contract ABIs to register for crafting (expected format %v)
Environment variable: %q`, `<contract>:<abi>:<bytecode>:<deployedBytecode>`, abiEnv)
	f.StringSlice(abiFlag, abiDefault, desc)
	_ = viper.BindPFlag(abiViperKey, f.Lookup(abiFlag))
}

// initializeABIs Read ABIs from ABI viper configuration
func initializeABIs(ctx context.Context, contractRegistryController *controllers.ContractRegistryController) {
	var contracts []*abi.Contract
	for _, ABI := range viper.GetStringSlice(abiViperKey) {
		c, err := abi.StringToContract(ABI)
		if err != nil {
			log.WithError(err).Fatalf("could not initialize contract-registry")
			return
		}
		contracts = append(contracts, c)
	}

	// Register contracts
	for _, contract := range contracts {
		_, err := contractRegistryController.RegisterContract(ctx, &svc.RegisterContractRequest{Contract: contract})

		if err != nil {
			log.WithError(err).Fatalf("could not register ABI")
		}
	}
}
