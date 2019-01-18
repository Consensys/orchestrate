package services

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ABIRegistry is an interface to manage ABIs
// TODO: extend including
//     - SetMethod() (ID string, err error)
//     - Functions for event ABIs
//     - Functions for contract ABIs
type ABIRegistry interface {
	// Retrieve method using unique identifier
	GetMethodByID(ID string) (abi.Method, error)
	// Retrieve method using 4 bytes signature
	GetMethodBySig(sig string) (abi.Method, error)
}

// Crafter takes a method abi and args to craft a transaction
type Crafter interface {
	Craft(method abi.Method, args ...string) ([]byte, error)
}
