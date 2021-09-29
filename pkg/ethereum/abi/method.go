package abi

import (
	"encoding/json"
	"strings"

	"github.com/consensys/orchestrate/pkg/errors"
	ethabi "github.com/consensys/orchestrate/pkg/go-ethereum/v1_9_12/accounts/abi"
)

func ParseMethod(methodABI []byte) (*ethabi.Method, error) {
	var method *ethabi.Method
	err := json.Unmarshal(methodABI, &method)
	if err != nil {
		return nil, err
	}
	return method, nil
}

// ParseMethodSignature create a method from a method signature string
func ParseMethodSignature(methodSig string) (*ethabi.Method, error) {
	splt := strings.Split(methodSig, "(")
	if len(splt) != 2 || splt[0] == "" || splt[1] == "" { // || splt[1][len(splt[1])-1:] != ")" {
		return nil, errors.InvalidSignatureError("Invalid method signature %q", methodSig)
	}

	method := &ethabi.Method{
		RawName: splt[0],
		Const:   false,
		Inputs:  ethabi.Arguments{},
	}

	inputArgs := splt[1][:len(splt[1])-1]
	if inputArgs != "" {
		for _, arg := range strings.Split(inputArgs, ",") {
			inputType, err := ethabi.NewType(arg, "", nil)
			if err != nil {
				return nil, errors.InvalidSignatureError(err.Error())
			}
			method.Inputs = append(method.Inputs, ethabi.Argument{Type: inputType})
		}
	}

	return method, nil
}
