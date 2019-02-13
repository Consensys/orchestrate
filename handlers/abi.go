package handlers

// TODO all this file should me moved out of the project and be replaced by a
// standalone ABIRegistry API
import (
	InfEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
)

// LoadABIRegistry creates an ABI registry and register contracts passed in environment variable in it
func LoadABIRegistry(abis map[string]string) *InfEth.ContractABIRegistry {
	registry := InfEth.NewContractABIRegistry()
	for k, v := range abis {
		registry.RegisterContract(k, []byte(v))
	}
	return registry
}
