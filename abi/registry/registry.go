package registry

import (
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

// Registry is an interface to manage ABIs
type Registry interface {
	// Register a new contract in the ABI
	RegisterContract(contract *abi.Contract) error
	// Retrieve contract ABI
	GetContractABI(contract *abi.Contract) ([]byte, error)
	// Returns the bytecode
	GetContractBytecode(contract *abi.Contract) ([]byte, error)
	// Returns the deployed bytecode
	GetContractDeployedBytecode(contract *abi.Contract) ([]byte, error)

	// Retrieve method using 4 bytes unique selector
	GetMethodsBySelector(selector [4]byte, contract common.AccountInstance) (*ethAbi.Method, []*ethAbi.Method, error)
	// Retrieve event using its signature hash
	GetEventsBySigHash(sigHash ethCommon.Hash, contract common.AccountInstance, indexedInputCount uint) (*ethAbi.Event, []*ethAbi.Event, error)

	// Request an update of the codehash of the contract address
	RequestAddressUpdate(contract common.AccountInstance) error
}
