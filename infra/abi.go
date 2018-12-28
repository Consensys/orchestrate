package infra

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ABIRegistry is an interface to manage ABIs
// TODO: extend including
//     - SetMethod() (ID string, err error)
//     - Functions for event ABIs
//     - Functions for contract ABIs
type ABIRegistry interface {
	// Retrieve method using unique identifier
	GetMethodByID(ID string) (*abi.Method, error)
	// Retrieve method using 4 bytes signature
	GetMethodBySig(sig string) (*abi.Method, error)
}

// DummyABIRegistry always return the same ABI method (useful for testing purpose)
type DummyABIRegistry struct {
	method *abi.Method
}

// NewDummyABIRegistry creates a new DummyABIgetter
func NewDummyABIRegistry(methodABI []byte) *DummyABIRegistry {
	var method abi.Method
	json.Unmarshal(methodABI, &method)
	return &DummyABIRegistry{&method}
}

// GetMethodByID return method ABI
func (g *DummyABIRegistry) GetMethodByID(ID string) (*abi.Method, error) {
	return g.method, nil
}

// GetMethodBySig return method ABI
func (g *DummyABIRegistry) GetMethodBySig(sig string) (*abi.Method, error) {
	return g.method, nil
}
