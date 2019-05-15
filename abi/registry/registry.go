package registry

import (
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
)

// Registry is an interface to manage ABIs
type Registry interface {
	// Register a new contract in the ABI
	RegisterContract(contract *abi.Contract) error

	// Retrieve method using 4 bytes unique selector
	GetMethodBySelector(selector string) (*ethabi.Method, error)
	// Retrieve method using signature
	GetMethodBySig(contract, signature string) (*ethabi.Method, error)

	// Retrieve event using 4 bytes unique selector
	GetEventBySelector(selector string) (*ethabi.Event, error)
	// Retrieve event using signature
	GetEventBySig(contract, signature string) (*ethabi.Event, error)

	// Returns the new bytecode
	GetBytecodeByID(id string) ([]byte, error)
}
