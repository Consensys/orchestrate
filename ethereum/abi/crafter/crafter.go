package crafter

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// Crafter takes a method abi and args to craft a transaction
type Crafter interface {
	// CraftCall craft a call Transaction
	CraftCall(method *ethabi.Method, args ...string) ([]byte, error)

	// CraftConstructor craft a Contract Deployment Transaction
	CraftConstructor(bytecode []byte, method *ethabi.Method, args ...string) ([]byte, error)
}

// PayloadCrafter is a structure that can Craft payloads
type PayloadCrafter struct{}

func bindArg(t *ethabi.Type, arg string) (interface{}, error) {
	switch t.T {
	case ethabi.AddressTy:
		if !ethcommon.IsHexAddress(arg) {
			return nil, errors.InvalidArgError("invalid ethereum address %q", arg).SetComponent(component)
		}
		return ethcommon.HexToAddress(arg), nil

	case ethabi.FixedBytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, errors.InvalidArgError("invalid bytes %q", arg).SetComponent(component)
		}
		array := reflect.New(t.Type).Elem()

		data = ethcommon.LeftPadBytes(data, t.Size)
		reflect.Copy(array, reflect.ValueOf(data[0:t.Size]))

		return array.Interface(), nil

	case ethabi.BytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, errors.InvalidArgError("invalid bytes %q", arg).SetComponent(component)
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
			return nil, errors.InvalidArgError("invalid boolean %q (expected one of %q)", arg, []string{"0x0", "false", "0", "0x1", "true", "1"}).
				SetComponent(component)
		}

	case ethabi.StringTy:
		return arg, nil

	case ethabi.ArrayTy, ethabi.SliceTy:
		return bindArrayArg(t, arg)

	case ethabi.TupleTy:
		return nil, errors.FeatureNotSupportedError("solidity tuple not supported yet").SetComponent(component)

	default:
		return nil, errors.FeatureNotSupportedError("solidity type %q not supported", t.T).SetComponent(component)
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

func bindBigIntArg(t *ethabi.Type, arg string) (interface{}, error) {
	// Check that it is a pointer to big int
	if t.Kind != reflect.Ptr {
		return nil, errors.InvalidArgError("bindBigIntArg: invalid type for %s - expected type kind %s but got %s", arg, reflect.Ptr, t.Kind).SetComponent(component)
	}

	// If arg is negative
	if arg != "" && arg[0] == '-' {
		raw, _, err := checkNumber(arg)
		if err != nil {
			return nil, errors.InvalidArgError("bindBigIntArg: invalid negative invalid number %q", err).SetComponent(component)
		}

		i := new(big.Int)
		i, ok := i.SetString(raw, 16)
		if !ok {
			return nil, errors.FromError(fmt.Errorf("bindBigIntArg: could not decode negative value of %s", arg)).SetComponent(component)
		}
		return i, nil
	}

	data, err := hexutil.DecodeBig(arg)
	if err != nil {
		return data, errors.InvalidArgError("bindBigIntArg invalid number %q", arg).SetComponent(component)
	}
	return data, nil
}

func bindIntArg(t *ethabi.Type, arg string) (interface{}, error) {
	raw, _, err := checkNumber(arg)
	if err != nil {
		return nil, errors.InvalidArgError("bindIntArg: invalid number %q", err).SetComponent(component)
	}

	number, err := strconv.ParseInt(raw, 16, t.Size)
	if err != nil {
		return nil, errors.InvalidArgError("bindIntArg: could not parse number %q", err).SetComponent(component)
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
		return nil, errors.InvalidArgError("bindIntArg: invalid size").SetComponent(component)
	}
}

func bindUintArg(t *ethabi.Type, arg string) (interface{}, error) {
	raw, isNegative, err := checkNumber(arg)
	if err != nil && isNegative {
		return nil, errors.InvalidArgError("bindUintArg: invalid number %q", err).SetComponent(component)
	}

	number, err := strconv.ParseUint(raw, 16, t.Size)
	if err != nil {
		return nil, errors.InvalidArgError("bindUintArg: could not parse number %q", err).SetComponent(component)
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
		return nil, errors.InvalidArgError("bindUintArg: invalid size").SetComponent(component)
	}
}

func bindArrayArg(t *ethabi.Type, arg string) (interface{}, error) {
	elemType, _ := ethabi.NewType(t.Elem.String(), "", nil)
	slice := reflect.MakeSlice(reflect.SliceOf(elemType.Type), 0, 0)

	var argArray []string
	err := json.Unmarshal([]byte(arg), &argArray)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// If t.Size == 0, then it is a dynamic array. We accept any length in this case.
	if t.Size != 0 && len(argArray) != t.Size {
		return nil,
			errors.InvalidArgError(
				"invalid size array %q (expected length %v but got %v)",
				arg, t.Size, len(argArray),
			).SetComponent(component)
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

// bindArgs cast string arguments into expected go-ethereum types
func bindArgs(method *ethabi.Method, args ...string) ([]interface{}, error) {
	if method.Inputs.LengthNonIndexed() != len(args) {
		return nil,
			errors.InvalidArgsCountError(
				"invalid arguments count (expected %v but got %v)",
				method.Inputs.LengthNonIndexed(), len(args),
			).SetComponent(component)
	}

	boundArgs := make([]interface{}, 0)
	for i := range method.Inputs.NonIndexed() {
		boundArg, err := bindArg(&method.Inputs.NonIndexed()[i].Type, args[i])
		if err != nil {
			return nil, err
		}
		boundArgs = append(boundArgs, boundArg)
	}
	return boundArgs, nil
}

// Pack automatically cast string args into correct Solidity type and pack arguments
func (c *PayloadCrafter) Pack(method *ethabi.Method, args ...string) ([]byte, error) {
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

// CraftCall craft a transaction call payload
func (c *PayloadCrafter) CraftCall(method *ethabi.Method, args ...string) ([]byte, error) {
	// Pack arguments
	arguments, err := c.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	return append(method.ID(), arguments...), nil
}

// CraftConstructor craft contract creation a transaction payload
func (c *PayloadCrafter) CraftConstructor(bytecode []byte, method *ethabi.Method, args ...string) ([]byte, error) {
	if len(bytecode) == 0 {
		return nil, errors.SolidityError("invalid empty bytecode").SetComponent(component)
	}

	// Pack arguments
	arguments, err := c.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	return append(bytecode, arguments...), nil
}

// SignatureToMethod create a method from a method signature string
func SignatureToMethod(methodSig string) (*ethabi.Method, error) {
	splt := strings.Split(methodSig, "(")
	if len(splt) != 2 || splt[0] == "" || splt[1] == "" { // || splt[1][len(splt[1])-1:] != ")" {
		return nil, errors.InvalidSignatureError(methodSig).SetComponent(component)
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
				return nil, errors.InvalidSignatureError(methodSig).SetMessage("invalid method signature (%v)", err).SetComponent(component)
			}
			method.Inputs = append(method.Inputs, ethabi.Argument{Type: inputType})
		}
	}

	return method, nil
}
