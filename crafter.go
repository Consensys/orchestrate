package ethereum

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var bytesTypes [33]reflect.Type

func copyBytesSliceToArray(b []byte, size int) interface{} {
	if size == 32 {
		rv := [32]byte{}
		copy(rv[:], b[0:size])
		return rv
	}

	if size == 16 {
		rv := [16]byte{}
		copy(rv[:], b[0:size])
		return rv
	}

	if size == 8 {
		rv := [8]byte{}
		copy(rv[:], b[0:size])
		return rv
	}

	if size == 1 {
		rv := [1]byte{}
		copy(rv[:], b[0:size])
		return rv
	}

	return nil
}

// PayloadCrafter is a structure that can Craft payloads
type PayloadCrafter struct{}

func bindArg(stringKind string, arg string) (interface{}, error) {
	switch {
	case stringKind == "address":
		if !common.IsHexAddress(arg) {
			return nil, fmt.Errorf("bindArg: %q is not a valid ethereum address", arg)
		}
		return common.HexToAddress(arg), nil

	case strings.HasPrefix(stringKind, "bytes"):
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, err
		}

		parts := regexp.MustCompile(`bytes([0-9]*)`).FindStringSubmatch(stringKind)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Arg format %q not known", stringKind)
		}

		if parts[1] == "" {
			return data, nil
		}

		size, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("Arg format %q not known", stringKind)
		}

		data = common.LeftPadBytes(data, size)
		b := copyBytesSliceToArray(data, size)

		if b == nil {
			return nil, fmt.Errorf("Arg format %q not known", stringKind)
		}

		return b, nil

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

	// In current version we only cover basic types (in particular we do not support arrays)
	default:
		return nil, fmt.Errorf("Arg format %q not known", stringKind)
	}
}

// bindArgs cast string arguments into expected go-ethereum types
func bindArgs(method abi.Method, args ...string) ([]interface{}, error) {
	if method.Inputs.LengthNonIndexed() != len(args) {
		return nil, fmt.Errorf("Expected %v inputs but got %v", method.Inputs.LengthNonIndexed(), len(args))
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

// Craft craft a transaction payload
func (c *PayloadCrafter) Craft(method abi.Method, args ...string) ([]byte, error) {
	// Cast arguments
	boundArgs, err := bindArgs(method, args...)
	if err != nil {
		return nil, err
	}

	// Pack arguments
	arguments, err := method.Inputs.Pack(boundArgs...)
	if err != nil {
		return nil, err
	}

	return append(method.Id(), arguments...), nil
}
