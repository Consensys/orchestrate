package abi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PayloadCrafter is a structure that can Craft payloads
type PayloadCrafter struct{}

func bindArg(t abi.Type, arg string) (interface{}, error) {

	switch t.T {
	case abi.AddressTy:
		if !common.IsHexAddress(arg) {
			return nil, fmt.Errorf("bindArg: %q is not a valid ethereum address", arg)
		}
		return common.HexToAddress(arg), nil

	case abi.FixedBytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, err
		}
		array := reflect.New(t.Type).Elem()

		data = common.LeftPadBytes(data, t.Size)
		reflect.Copy(array, reflect.ValueOf(data[0:t.Size]))

		return array.Interface(), nil

	case abi.BytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, err
		}
		return data, nil

	case abi.IntTy, abi.UintTy:
		// In current version we bind all types of integers to *big.Int
		// Meaning we do not yet support int8, int16, int32, int64, uint8, uin16, uint32, uint64
		return hexutil.DecodeBig(arg)

	case abi.BoolTy:
		switch arg {
		case "0x1", "true", "1":
			return true, nil
		case "0x0", "false", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("bindArg: %v is not a bool", arg)
		}

	case abi.StringTy:
		return arg, nil

	case abi.ArrayTy:
		return bindArrayArg(t, arg)

	case abi.SliceTy:
		return bindArrayArg(t, arg)

	// TODO: handle tuple (struct in solidity)

	// In current version we only cover basic types (in particular we do not support arrays)
	default:
		return nil, fmt.Errorf("Arg format %v not known", t.T)
	}
}

func bindArrayArg(t abi.Type, arg string) (interface{}, error) {

	elemType, _ := abi.NewType(t.Elem.String(), nil)
	slice := reflect.MakeSlice(reflect.SliceOf(elemType.Type), 0, 0)

	arg = strings.TrimSuffix(strings.TrimPrefix(arg, "["), "]")
	argArray := strings.Split(arg, ",")

	// If t.Size == 0, then it is a dynamic array. We accept any length in this case.
	if len(argArray) != t.Size && t.Size != 0 {
		return nil, fmt.Errorf("Craft array error: %q is not well separated", argArray)
	}
	for _, v := range argArray {
		typedArg, err := bindArg(elemType, v)
		if err != nil {
			return nil, fmt.Errorf("Craft array error: %v", err)
		}
		slice = reflect.Append(slice, reflect.ValueOf(typedArg))
	}
	return slice.Interface(), nil
}

// bindArgs cast string arguments into expected go-ethereum types
func bindArgs(method abi.Method, args ...string) ([]interface{}, error) {
	if method.Inputs.LengthNonIndexed() != len(args) {
		return nil, fmt.Errorf("Expected %v inputs but got %v", method.Inputs.LengthNonIndexed(), len(args))
	}

	boundArgs := make([]interface{}, 0)
	for i, arg := range method.Inputs.NonIndexed() {
		boundArg, err := bindArg(arg.Type, args[i])
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

// CraftConstructor craft contract creation a transaction payload
func (c *PayloadCrafter) CraftConstructor(method abi.Method, args ...string) ([]byte, error) {
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

	return arguments, nil
}
