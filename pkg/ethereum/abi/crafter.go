package abi

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

type Crafter interface {
	CraftCall(sig string, args ...string) ([]byte, error)
	CraftConstructor(bytecode []byte, sig string, args ...string) ([]byte, error)
}

type BaseCrafter struct{}

// CraftCall craft a transaction call payload
func (c *BaseCrafter) CraftCall(sig string, args ...string) ([]byte, error) {
	method, err := ParseMethodSignature(sig)
	if err != nil {
		return nil, err
	}

	// Pack arguments
	arguments, err := Pack(method, args...)
	if err != nil {
		return nil, err
	}

	return append(method.ID(), arguments...), nil
}

// CraftConstructor craft contract creation a transaction payload
func (c *BaseCrafter) CraftConstructor(bytecode []byte, sig string, args ...string) ([]byte, error) {
	method, err := ParseMethodSignature(sig)
	if err != nil {
		return nil, err
	}

	if len(bytecode) == 0 {
		return nil, errors.SolidityError("empty bytecode")
	}

	// Pack arguments
	arguments, err := Pack(method, args...)
	if err != nil {
		return nil, err
	}

	return append(bytecode, arguments...), nil
}
