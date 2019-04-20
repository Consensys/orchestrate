package abi

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

// FormatIndexedArg transforms a data to string
func FormatIndexedArg(t abi.Type, arg common.Hash) (string, error) {

	switch t.T {
	case abi.BoolTy, abi.StringTy:
		return fmt.Sprintf("%v", arg), nil
	case abi.IntTy, abi.UintTy:
		num := new(big.Int).SetBytes(arg[:])
		return fmt.Sprintf("%v", num), nil
	case abi.AddressTy:
		return common.HexToAddress(arg.Hex()).Hex(), nil
	case abi.FixedBytesTy:
		return fmt.Sprintf("%v", hexutil.Encode(arg[common.HashLength-t.Type.Size():])), nil
	case abi.BytesTy, abi.ArrayTy, abi.TupleTy:
		return "", fmt.Errorf("unable to decode %v type", t.Kind)
	default:
		return fmt.Sprintf("%v", arg), nil
	}
}

// ArrayToByteSlice creates a new byte slice with the exact same size as value
// and copies the bytes in value to the new slice.
func ArrayToByteSlice(value reflect.Value) reflect.Value {
	slice := reflect.MakeSlice(reflect.TypeOf([]byte{}), value.Len(), value.Len())
	reflect.Copy(slice, value)
	return slice
}

// FormatNonIndexedArrayArg transforms a data to string
func FormatNonIndexedArrayArg(t abi.Type, arg interface{}) (string, error) {

	elemType, _ := abi.NewType(t.Elem.String(), nil)

	var arrayArgString []string
	for i := 0; i < t.Size; i++ {
		argString, _ := FormatNonIndexedArg(elemType, reflect.ValueOf(arg).Index(i).Interface())
		arrayArgString = append(arrayArgString, argString)
	}

	return "[" + strings.Join(arrayArgString, ",") + "]", nil
}

// FormatNonIndexedArg transforms a data to string
func FormatNonIndexedArg(t abi.Type, arg interface{}) (string, error) {

	switch t.T {
	case abi.IntTy, abi.UintTy, abi.BoolTy, abi.StringTy:
		return fmt.Sprintf("%v", arg), nil
	case abi.AddressTy:
		return arg.(common.Address).Hex(), nil
	case abi.FixedBytesTy:
		slice := ArrayToByteSlice(reflect.ValueOf(arg))
		return hexutil.Encode(slice.Bytes()), nil
	case abi.BytesTy:
		return hexutil.Encode(reflect.ValueOf(arg).Bytes()), nil
	case abi.ArrayTy:
		return FormatNonIndexedArrayArg(t, arg)
	case abi.TupleTy:
	default:
		return fmt.Sprintf("%v", arg), nil
	}
	return "", fmt.Errorf("unable to decode %v type", t.Kind)
}

// Decode event data to string
func Decode(event *abi.Event, txLog *ethpb.Log) (map[string]string, error) {
	expectedTopics := len(event.Inputs) - event.Inputs.LengthNonIndexed()
	if expectedTopics != len(txLog.Topics)-1 {
		return nil, fmt.Errorf("Error: Topics length does not match with abi event: expected %v but got %v", expectedTopics, len(txLog.Topics)-1)
	}

	unpackValues, err := event.Inputs.UnpackValues(txLog.Data)
	if err != nil {
		return nil, err
	}

	var (
		topicIndex        = 1
		unpackValuesIndex = 0
	)
	logMapping := make(map[string]string, len(event.Inputs))

	for _, arg := range event.Inputs {
		var decoded string
		if arg.Indexed {
			decoded, _ = FormatIndexedArg(arg.Type, common.HexToHash(txLog.Topics[topicIndex]))
			topicIndex++
		} else {
			decoded, _ = FormatNonIndexedArg(arg.Type, unpackValues[unpackValuesIndex])
			unpackValuesIndex++
		}
		logMapping[arg.Name] = decoded
	}

	return logMapping, nil
}
