package abi

// Pack automatically cast string args into correct Solidity type and pack arguments
func Pack(method *Method, args ...string) ([]byte, error) {
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
