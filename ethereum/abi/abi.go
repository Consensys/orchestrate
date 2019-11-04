package abi

import (
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// MarshalABI marshal an Ethereum ABI object
func MarshalABI(abi *ethabi.ABI) ([]byte, error) {
	var fields []field

	// Add constructor to slice
	fields = append(fields, field{
		Type:   "constructor",
		Inputs: fromArguments(abi.Constructor.Inputs),
	})

	// Add Methods to slice
	for _, method := range abi.Methods {
		fields = append(fields, field{
			Type:     "function",
			Name:     method.Name,
			Inputs:   fromArguments(method.Inputs),
			Outputs:  fromArguments(method.Outputs),
			Constant: method.Const,
		})
	}

	// Add events to slice
	for _, event := range abi.Events {
		fields = append(fields, field{
			Type:      "event",
			Name:      event.Name,
			Anonymous: event.Anonymous,
			Inputs:    fromArguments(event.Inputs),
		})
	}

	return json.Marshal(fields)
}

type field struct {
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Constant  bool       `json:"constant"`
	Anonymous bool       `json:"anonymous"`
	Inputs    []argument `json:"inputs"`
	Outputs   []argument `json:"outputs"`
}

type argument struct {
	Argument ethabi.Argument
}

func (arg *argument) MarshalJSON() ([]byte, error) {
	if strings.Contains(arg.Argument.Type.String(), "(") {
		return []byte{}, errors.FeatureNotSupportedError("Cannot marshal tuple type argument")
	}

	return json.Marshal(&ethabi.ArgumentMarshaling{
		Name:    arg.Argument.Name,
		Indexed: arg.Argument.Indexed,
		Type:    arg.Argument.Type.String(),
	})
}

func fromArguments(args []ethabi.Argument) []argument {
	var rv []argument
	for i := range args {
		rv = append(rv, argument{args[i]})
	}
	return rv
}
