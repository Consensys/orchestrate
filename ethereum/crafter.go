package ethereum

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func bindArg(stringKind string, arg string) (interface{}, error) {
	// In current version we assume every arguments to be strings
	switch {
	case stringKind == "address":
		if !common.IsHexAddress(arg) {
			return nil, fmt.Errorf("bindArg: %q is not a valid ethereum address", arg)
		}
		return common.HexToAddress(arg), nil

	case stringKind == "bytes":
		// We do not yet support bytesx (e.g. bytes1, bytes2...)
		return hexutil.Decode(arg)

	case strings.HasPrefix(stringKind, "int") || strings.HasPrefix(stringKind, "uint"):
		// In current version we bind all types of integers to *big.Int
		// Meaning we do not yet support int8, int16, int32, int64, uint8, uin16, uint32, uint64
		return hexutil.DecodeBig(arg)

	case stringKind == "bool":
		b, err := hexutil.DecodeBig(arg)
		if err != nil {
			return nil, err
		}
		return b.Int64() > 0, nil

	case strings.HasPrefix(stringKind, "string"):
		return arg, nil

	// In current version we only cover basic types
	default:
		return nil, nil
	}
}

// bindArgs cast string arguments to expected go-ethereum type before crafting
func bindArgs(method *abi.Method, args []string) ([]interface{}, error) {
	if method.Inputs.LengthNonIndexed() != len(args) {
		return nil, fmt.Errorf("BindArgs: expected %v inputs but got %v", method.Inputs.LengthNonIndexed(), len(args))
	}
	boundArgs := make([]interface{}, 0)
	for i, arg := range method.Inputs.NonIndexed() {
		boundArg, err := bindArg(arg.Type.String(), args[i])
		if err != nil {
			return nil, err
		}
		boundArgs = append(boundArgs, boundArg)
	}
	return boundArgs, nil
}

// CraftPayload craft a transaction payload
func CraftPayload(method *abi.Method, args []string) ([]byte, error) {
	boundArgs, err := bindArgs(method, args)

	if err != nil {
		return nil, err
	}

	arguments, err := method.Inputs.Pack(boundArgs...)

	if err != nil {
		return nil, err
	}

	return append(method.Id(), arguments...), nil
}
