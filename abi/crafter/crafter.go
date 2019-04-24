package crafter

import (
	"fmt"
	"reflect"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Crafter takes a method abi and args to craft a transaction
type Crafter interface {
	// CraftCall craft a call Transaction
	CraftCall(method ethabi.Method, args ...string) ([]byte, error)

	// CraftConstructor craft a Contract Deployment Transaction
	CraftConstructor(bytecode []byte, method ethabi.Method, args ...string) ([]byte, error)
}

// PayloadCrafter is a structure that can Craft payloads
type PayloadCrafter struct{}

func bindArg(t *ethabi.Type, arg string) (interface{}, error) {
	switch t.T {
	case ethabi.AddressTy:
		if !ethcommon.IsHexAddress(arg) {
			return nil, fmt.Errorf("bindArg: %q is not a valid ethereum address", arg)
		}
		return ethcommon.HexToAddress(arg), nil

	case ethabi.FixedBytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, err
		}
		array := reflect.New(t.Type).Elem()

		data = ethcommon.LeftPadBytes(data, t.Size)
		reflect.Copy(array, reflect.ValueOf(data[0:t.Size]))

		return array.Interface(), nil

	case ethabi.BytesTy:
		data, err := hexutil.Decode(arg)
		if err != nil {
			return data, err
		}
		return data, nil

	case ethabi.IntTy, ethabi.UintTy:
		// In current version we bind all types of integers to *big.Int
		// Meaning we do not yet support int8, int16, int32, int64, uint8, uin16, uint32, uint64
		return hexutil.DecodeBig(arg)

	case ethabi.BoolTy:
		switch arg {
		case "0x1", "true", "1":
			return true, nil
		case "0x0", "false", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("bindArg: %v is not a bool", arg)
		}

	case ethabi.StringTy:
		return arg, nil

	case ethabi.ArrayTy:
		return bindArrayArg(t, arg)

	case ethabi.SliceTy:
		return bindArrayArg(t, arg)

	// TODO: handle tuple (struct in solidity)

	default:
		return nil, fmt.Errorf("arg format %v not known", t.T)
	}
}

func bindArrayArg(t *ethabi.Type, arg string) (interface{}, error) {
	elemType, _ := ethabi.NewType(t.Elem.String(), nil)
	slice := reflect.MakeSlice(reflect.SliceOf(elemType.Type), 0, 0)

	arg = strings.TrimSuffix(strings.TrimPrefix(arg, "["), "]")
	argArray := strings.Split(arg, ",")

	// If t.Size == 0, then it is a dynamic array. We accept any length in this case.
	if len(argArray) != t.Size && t.Size != 0 {
		return nil, fmt.Errorf("craft array error: %q is not well separated", argArray)
	}
	for _, v := range argArray {
		typedArg, err := bindArg(&elemType, v)
		if err != nil {
			return nil, fmt.Errorf("craft array error: %v", err)
		}
		slice = reflect.Append(slice, reflect.ValueOf(typedArg))
	}
	return slice.Interface(), nil
}

// bindArgs cast string arguments into expected go-ethereum types
func bindArgs(method ethabi.Method, args ...string) ([]interface{}, error) {
	if method.Inputs.LengthNonIndexed() != len(args) {
		return nil, fmt.Errorf("expected %v inputs but got %v", method.Inputs.LengthNonIndexed(), len(args))
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
func (c *PayloadCrafter) Pack(method ethabi.Method, args ...string) ([]byte, error) {
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
func (c *PayloadCrafter) CraftCall(method ethabi.Method, args ...string) ([]byte, error) {
	// Pack arguments
	arguments, err := c.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	return append(method.Id(), arguments...), nil
}

// CraftConstructor craft contract creation a transaction payload
func (c *PayloadCrafter) CraftConstructor(bytecode []byte, method ethabi.Method, args ...string) ([]byte, error) {
	if len(bytecode) == 0 {
		return nil, fmt.Errorf("invalid empty bytecode")
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
	if len(splt) != 2 || splt[0] == "" || len(splt[1]) <= 1 {
		return nil, fmt.Errorf("invalid method signature %q, expected Function(type1,type2,...)", methodSig)
	}
	inputArgs := strings.Split(splt[1][:len(splt[1])-1], ",")

	method := &ethabi.Method{
		Name:  splt[0],
		Const: false,
	}
	for _, arg := range inputArgs {
		inputType, err := ethabi.NewType(arg, nil)
		if err != nil {
			return nil, fmt.Errorf("invalid method signature format, cannot cast type: %v", err)
		}
		method.Inputs = append(method.Inputs, ethabi.Argument{Type: inputType})
	}

	return method, nil
}
