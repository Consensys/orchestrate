package abi

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// BindArgs cast string arguments into expected go-ethereum types
func BindArgs(arguments *ethabi.Arguments, args ...string) ([]interface{}, error) {
	if arguments.LengthNonIndexed() != len(args) {
		return nil,
			errors.InvalidArgsCountError(
				"invalid arguments count (expected %v but got %v)",
				arguments.LengthNonIndexed(), len(args),
			)
	}

	boundArgs := make([]interface{}, 0)
	for i := range arguments.NonIndexed() {
		boundArg, err := BindArg(&arguments.NonIndexed()[i].Type, args[i])
		if err != nil {
			return nil, err
		}
		boundArgs = append(boundArgs, boundArg)
	}

	return boundArgs, nil
}

// BindArg cast string argument into expected go-ethereum type
func BindArg(t *ethabi.Type, arg string) (interface{}, error) {
	return bindArg(t, arg)
}

func bindArg(t *ethabi.Type, arg string) (interface{}, error) {
	switch t.T {
	case ethabi.AddressTy:
		if !ethcommon.IsHexAddress(arg) {
			return nil, errors.InvalidArgError("invalid ethereum address %q", arg)
		}
		return ethcommon.HexToAddress(arg), nil

	case ethabi.FixedBytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, errors.InvalidArgError("invalid bytes %q", arg)
		}
		array := reflect.New(t.Type).Elem()

		data = ethcommon.LeftPadBytes(data, t.Size)
		reflect.Copy(array, reflect.ValueOf(data[0:t.Size]))

		return array.Interface(), nil

	case ethabi.BytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, errors.InvalidArgError("invalid bytes %q", arg)
		}
		return data, nil

	case ethabi.IntTy:
		switch t.Size {
		// only int of size 8, 16, 32, 64 should be bind in int
		// other ones should be in *big.Int see packNum in go-ethereum/accounts/abi/pack.go https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/pack.go
		case 8, 16, 32, 64:
			return bindIntArg(t, arg)
		default:
			return bindBigIntArg(t, arg)
		}

	case ethabi.UintTy:
		switch t.Size {
		// only uint of size 8, 16, 32, 64 should be bind in uint
		// other ones should be in *big.Int see packNum in go-ethereum/accounts/abi/pack.go https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/pack.go
		case 8, 16, 32, 64:
			return bindUintArg(t, arg)
		default:
			return bindBigIntArg(t, arg)
		}

	case ethabi.BoolTy:
		switch arg {
		case "0x1", "true", "1":
			return true, nil
		case "0x0", "false", "0":
			return false, nil
		default:
			return nil, errors.InvalidArgError("invalid boolean %q (expected one of %q)", arg, []string{"0x0", "false", "0", "0x1", "true", "1"})
		}

	case ethabi.StringTy:
		return arg, nil

	case ethabi.ArrayTy, ethabi.SliceTy:
		return bindArrayArg(t, arg)

	case ethabi.TupleTy:
		return nil, errors.FeatureNotSupportedError("solidity tuple not supported yet")

	default:
		return nil, errors.FeatureNotSupportedError("solidity type %q not supported", t.T)
	}
}

func bindBigIntArg(t *ethabi.Type, arg string) (interface{}, error) {
	// Check that it is a pointer to big int
	if t.Kind != reflect.Ptr {
		return nil, errors.InvalidArgError("bindBigIntArg: invalid type for %s - expected type kind %s but got %s", arg, reflect.Ptr, t.Kind)
	}

	// If arg is negative
	if arg != "" && arg[0] == '-' {
		raw, _, err := checkNumber(arg)
		if err != nil {
			return nil, errors.InvalidArgError("bindBigIntArg: invalid negative invalid number %q", err)
		}

		i := new(big.Int)
		i, ok := i.SetString(raw, 16)
		if !ok {
			return nil, errors.FromError(fmt.Errorf("bindBigIntArg: could not decode negative value of %s", arg))
		}
		return i, nil
	}

	data, err := hexutil.DecodeBig(arg)
	if err != nil {
		return data, errors.InvalidArgError("bindBigIntArg invalid number %q", arg)
	}
	return data, nil
}

func bindIntArg(t *ethabi.Type, arg string) (interface{}, error) {
	raw, _, err := checkNumber(arg)
	if err != nil {
		return nil, errors.InvalidArgError("bindIntArg: invalid number %q", err)
	}

	number, err := strconv.ParseInt(raw, 16, t.Size)
	if err != nil {
		return nil, errors.InvalidArgError("bindIntArg: could not parse number %q", err)
	}

	switch t.Size {
	case 8:
		return int8(number), nil
	case 16:
		return int16(number), nil
	case 32:
		return int32(number), nil
	case 64:
		return number, nil
	default:
		return nil, errors.InvalidArgError("bindIntArg: invalid size")
	}
}

func bindUintArg(t *ethabi.Type, arg string) (interface{}, error) {
	raw, isNegative, err := checkNumber(arg)
	if err != nil && isNegative {
		return nil, errors.InvalidArgError("bindUintArg: invalid number %q", err)
	}

	number, err := strconv.ParseUint(raw, 16, t.Size)
	if err != nil {
		return nil, errors.InvalidArgError("bindUintArg: could not parse number %q", err)
	}

	switch t.Size {
	case 8:
		return uint8(number), nil
	case 16:
		return uint16(number), nil
	case 32:
		return uint32(number), nil
	case 64:
		return number, nil
	default:
		return nil, errors.InvalidArgError("bindUintArg: invalid size")
	}
}

// Code inspired by github.com/ethereum/go-ethereum/common/hexutil/hexutil.go
func checkNumber(input string) (raw string, isNegative bool, err error) {
	if input == "" {
		return "", isNegative, hexutil.ErrEmptyString
	}
	if input[0] == '-' {
		isNegative = true
		input = input[1:]
	}
	if !has0xPrefix(input) {
		return "", isNegative, hexutil.ErrMissingPrefix
	}
	input = input[2:]
	if input == "" {
		return "", isNegative, hexutil.ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return "", isNegative, hexutil.ErrLeadingZero
	}
	if isNegative {
		input = "-" + input
	}
	return input, isNegative, nil
}

// Code inspired by github.com/ethereum/go-ethereum/common/hexutil/hexutil.go
func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func bindArrayArg(t *ethabi.Type, arg string) (interface{}, error) {
	elemType, _ := ethabi.NewType(t.Elem.String(), "", nil)
	slice := reflect.MakeSlice(reflect.SliceOf(elemType.Type), 0, 0)

	var argArray []string
	err := json.Unmarshal([]byte(arg), &argArray)
	if err != nil {
		return nil, errors.FromError(err)
	}

	// If t.Size == 0, then it is a dynamic array. We accept any length in this case.
	if t.Size != 0 && len(argArray) != t.Size {
		return nil,
			errors.InvalidArgError(
				"invalid size array %q (expected length %v but got %v)",
				arg, t.Size, len(argArray),
			)
	}
	for _, v := range argArray {
		typedArg, err := bindArg(&elemType, v)
		if err != nil {
			return nil, err
		}
		slice = reflect.Append(slice, reflect.ValueOf(typedArg))
	}
	return slice.Interface(), nil
}
