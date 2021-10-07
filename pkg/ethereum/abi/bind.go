package abi

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/consensys/orchestrate/pkg/errors"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const trueStr = "true"
const falseStr = "false"

// BindArgs cast string arguments into expected go-ethereum types
func BindArgs(arguments *ethabi.Arguments, args ...string) ([]interface{}, error) {
	if len(arguments.NonIndexed()) != len(args) {
		return nil,
			errors.InvalidArgsCountError(
				"invalid arguments count (expected %v but got %v)",
				len(arguments.NonIndexed()), len(args),
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
	switch t.T {
	case ethabi.AddressTy:
		if !ethcommon.IsHexAddress(arg) {
			return nil, errors.InvalidArgError("invalid ethereum address %q", arg)
		}
		return ethcommon.HexToAddress(arg), nil

	case ethabi.FixedBytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, errors.InvalidArgError("invalid fixed bytes %q", arg)
		}
		if len(data) > t.Size {
			return nil, errors.InvalidArgError("invalid fixed bytes %s of size %d - too big for %s", arg, len(data), t.String())
		}
		array := reflect.New(t.GetType()).Elem()

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
		isNegative, isHex, err := checkArgNumber(arg)
		if err != nil {
			return nil, err
		}
		switch t.Size {
		// only int of size 8, 16, 32, 64 should be bind in int
		// other ones should be in *big.Int see packNum in go-ethereum/accounts/abi/pack.go https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/pack.go
		case 8, 16, 32, 64:
			return bindIntArg(t, arg, isNegative, isHex)
		default:
			return bindBigIntArg(t, arg, isNegative, isHex)
		}

	case ethabi.UintTy:
		isNegative, has0xPrefix, err := checkArgNumber(arg)
		if err != nil {
			return nil, err
		}
		if isNegative {
			return nil, errors.InvalidArgError("did not expected negative value %s for %s", arg, t.String())
		}
		switch t.Size {
		// only uint of size 8, 16, 32, 64 should be bind in uint
		// other ones should be in *big.Int see packNum in go-ethereum/accounts/abi/pack.go https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/pack.go
		case 8, 16, 32, 64:
			return bindUintArg(t, arg, has0xPrefix)
		default:
			return bindBigIntArg(t, arg, isNegative, has0xPrefix)
		}

	case ethabi.BoolTy:
		switch arg {
		case "0x1", trueStr, "1":
			return true, nil
		case "0x0", falseStr, "0":
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
		return nil, errors.FeatureNotSupportedError("solidity type %s not supported", t.String())
	}
}

func checkArgNumber(arg string) (isNeg, hasHexPrefix bool, err error) {
	if arg == "" {
		return false, false, errors.InvalidArgError("did not expected empty uint/int value")
	}
	neg := isNegative(arg)
	if neg {
		arg = arg[1:]
	}
	hasPrefix := has0xPrefix(arg)
	if hasPrefix {
		// Checks if it is a valid hex
		// inspired by github.com/ethereum/go-ethereum/common/hexutil/hexutil.go
		arg = arg[2:]
		if arg == "" {
			return neg, hasPrefix, errors.InvalidArgError("invalid number - no value after 0x prefix")
		}
		if len(arg) > 1 && arg[0] == '0' {
			return neg, hasPrefix, errors.InvalidArgError("invalid number - got: %s", hexutil.ErrLeadingZero)
		}

	}
	return neg, hasPrefix, nil
}

func bindBigIntArg(t *ethabi.Type, arg string, isNegative, has0xPrefix bool) (interface{}, error) {
	// Check that it is a pointer to big int
	if t.GetType().Kind() != reflect.Ptr {
		return nil, errors.InvalidArgError("invalid type for %s - expected type kind %s but got %s", arg, reflect.Ptr, t.String())
	}

	base := 10
	raw := arg
	if has0xPrefix {
		if isNegative {
			raw = fmt.Sprintf("-%s", arg[3:])
		} else {
			raw = arg[2:]
		}
		base = 16
	}

	bigNumber, ok := new(big.Int).SetString(raw, base)
	if !ok {
		return nil, errors.InvalidArgError("invalid %s %s", t.String(), arg)
	}
	return bigNumber, nil
}

func bindIntArg(t *ethabi.Type, arg string, isNegative, has0xPrefix bool) (interface{}, error) {
	base := 10
	raw := arg
	if has0xPrefix {
		if isNegative {
			raw = fmt.Sprintf("-%s", arg[3:])
		} else {
			raw = arg[2:]
		}
		base = 16
	}

	number, err := strconv.ParseInt(raw, base, t.Size)
	if err != nil {
		return nil, errors.InvalidArgError("could not parse number %s - got %q", arg, err)
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
		return nil, errors.InvalidArgError("invalid size %d - could not bind %s", t.Size, arg)
	}
}

func bindUintArg(t *ethabi.Type, arg string, has0xPrefix bool) (interface{}, error) {
	base := 10
	raw := arg
	if has0xPrefix {
		raw = arg[2:]
		base = 16
	}

	number, err := strconv.ParseUint(raw, base, t.Size)
	if err != nil {
		return nil, errors.InvalidArgError("could not parse number %s - got %q", arg, err)
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
		return nil, errors.InvalidArgError("invalid size %d - could not bind %s", t.Size, arg)
	}
}

func bindArrayArg(t *ethabi.Type, arg string) (interface{}, error) {
	elemType, _ := ethabi.NewType(t.Elem.String(), "", nil)
	slice := reflect.MakeSlice(reflect.SliceOf(elemType.GetType()), 0, 0)

	var argArray []string
	err := json.Unmarshal([]byte(arg), &argArray)
	if err != nil {
		return nil, errors.InvalidArgError("could not parse array %s for %s - got %v", arg, t.String(), err)
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
		typedArg, err := BindArg(&elemType, v)
		if err != nil {
			return nil, err
		}
		slice = reflect.Append(slice, reflect.ValueOf(typedArg))
	}
	return slice.Interface(), nil
}

// Code inspired by github.com/ethereum/go-ethereum/common/hexutil/hexutil.go
func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func isNegative(input string) bool {
	return input[0] == '-'
}
