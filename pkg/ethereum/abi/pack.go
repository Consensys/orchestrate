package abi

import (
	ethabi "github.com/ConsenSys/orchestrate/pkg/go-ethereum/v1_9_12/accounts/abi"
)

// Pack automatically cast string args into correct Solidity type and pack arguments
func Pack(method *ethabi.Method, args ...string) ([]byte, error) {
	// Cast arguments
	boundArgs, err := BindArgs(&method.Inputs, args...)
	if err != nil {
		return nil, err
	}

	// Pack arguments
	arguments, err := method.Inputs.Pack(boundArgs...)
	if err != nil {
		return nil, err
	}

	return arguments, nil
}
