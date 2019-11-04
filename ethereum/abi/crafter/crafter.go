package crafter

import (
	"reflect"
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

	case ethabi.IntTy, ethabi.UintTy:
		// In current version we bind all types of integers to *big.Int
		// Meaning we do not yet support int8, int16, int32, int64, uint8, uin16, uint32, uint64
		data, err := hexutil.DecodeBig(arg)
		if err != nil {
			return data, errors.InvalidArgError("invalid int/uint %q", arg).SetComponent(component)
		}
		return data, nil

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

func bindArrayArg(t *ethabi.Type, arg string) (interface{}, error) {
	elemType, _ := ethabi.NewType(t.Elem.String(), nil)
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
			inputType, err := ethabi.NewType(arg, nil)
			if err != nil {
				return nil, errors.InvalidSignatureError(methodSig).SetMessage("invalid method signature (%v)", err).SetComponent(component)
			}
			method.Inputs = append(method.Inputs, ethabi.Argument{Type: inputType})
		}
	}

	return method, nil
}
