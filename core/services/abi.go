package services

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	abipb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/abi"
)

// ABIRegistry is an interface to manage ABIs
type ABIRegistry interface {
	// Retrieve method using unique identifier
	GetMethodByID(ID string) (abi.Method, error)
	// Retrieve method using 4 bytes signature
	GetMethodBySig(sig string) (abi.Method, error)
	// Retrieve event using unique identifier
	GetEventByID(ID string) (abi.Event, error)
	// Retrieve event using 4 bytes signature
	GetEventBySig(sig string) (abi.Event, error)
	// Register a new contract in the ABI
	RegisterContract(contract *abipb.Contract) error
	// Returns the new bytecode
	GetBytecodeByID(id string) ([]byte, error)
}

// Crafter takes a method abi and args to craft a transaction
type Crafter interface {
	Craft(method abi.Method, args ...string) ([]byte, error)
	CraftConstructor(method abi.Method, args ...string) ([]byte, error)
}
