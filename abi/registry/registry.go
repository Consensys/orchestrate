package abi

import (
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
)

// Registry is an interface to manage ABIs
type Registry interface {
	// Retrieve method using unique identifier
	GetMethodByID(ID string) (ethabi.Method, error)
	// Retrieve method using 4 bytes signature
	GetMethodBySig(sig string) (ethabi.Method, error)
	// Retrieve event using unique identifier
	GetEventByID(ID string) (ethabi.Event, error)
	// Retrieve event using 4 bytes signature
	GetEventBySig(sig string) (ethabi.Event, error)
	// Register a new contract in the ABI
	RegisterContract(contract *abi.Contract) error
	// Returns the new bytecode
	GetBytecodeByID(id string) ([]byte, error)
}
