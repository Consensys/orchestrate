package abi

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

// FormatIndexedArg transforms a data to string
func FormatIndexedArg(t *abi.Type, arg common.Hash) (string, error) {

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

// GetElemType returns the underlying element type of an array or slice
func GetElemType(t *abi.Type) (abi.Type, error) {

	switch strings.Count(t.Elem.String(), "(") {
	case 0:
		// Not a struct - able to return the correct type
		return abi.NewType(t.Elem.String(), nil)
	case 1:
		// Simple struct containing elementary types
		nestedTypes := regexp.MustCompile(`\((.*?)\)`).FindStringSubmatch(fmt.Sprintf("%v", t.Elem))
		nestedTypesList := strings.Split(nestedTypes[1], ",")
		tupleArgs := make([]abi.ArgumentMarshaling, t.Type.Elem().NumField())

		for i := 0; i < t.Type.Elem().NumField(); i++ {
			tupleArgs[i] = abi.ArgumentMarshaling{
				Name: t.Type.Elem().Field(i).Name,
				Type: nestedTypesList[i],
			}
		}

		return abi.NewType("tuple", tupleArgs)
	}

	return abi.Type{}, fmt.Errorf("decoder: cannot get Elem type of %v", t)
}

// FormatNonIndexedArrayArg transforms a data to string
func FormatNonIndexedArrayArg(t *abi.Type, arg interface{}) (string, error) {

	elemType, _ := GetElemType(t)

	var arrayArgString []string
	for i := 0; i < t.Size; i++ {
		argString, _ := FormatNonIndexedArg(&elemType, reflect.ValueOf(arg).Index(i).Interface())
		arrayArgString = append(arrayArgString, argString)
	}

	jsonArgs, err := json.Marshal(arrayArgString)
	if err != nil {
		return "", err
	}
	return string(jsonArgs), nil
}

// FormatNonIndexedSliceArg transforms a slice data to string
func FormatNonIndexedSliceArg(t *abi.Type, arg interface{}) (string, error) {
	val := reflect.ValueOf(arg)

	elemType, _ := GetElemType(t)

	var sliceArgString []string
	for i := 0; i < val.Len(); i++ {
		argString, _ := FormatNonIndexedArg(&elemType, val.Index(i).Interface())
		sliceArgString = append(sliceArgString, argString)
	}

	jsonArgs, err := json.Marshal(sliceArgString)
	if err != nil {
		return "", err
	}

	return string(jsonArgs), nil
}

// FormatNonIndexedTupleArg transforms a struct data to string
func FormatNonIndexedTupleArg(t *abi.Type, arg interface{}) (string, error) {

	val := reflect.ValueOf(arg)

	tuple := make(map[string]string, len(t.TupleElems))
	for i, elemeType := range t.TupleElems {
		var decoded string
		decoded, _ = FormatNonIndexedArg(elemeType, val.Field(i).Interface())
		tuple[abi.ToCamelCase(t.TupleRawNames[i])] = decoded
	}
	jsonArgs, err := json.Marshal(tuple)
	if err != nil {
		return "", err
	}

	return string(jsonArgs), nil
}

// FormatNonIndexedArg transforms a data to string
func FormatNonIndexedArg(t *abi.Type, arg interface{}) (string, error) {

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
	case abi.SliceTy:
		return FormatNonIndexedSliceArg(t, arg)
	case abi.TupleTy:
		return FormatNonIndexedTupleArg(t, arg)
	default:
		return "", fmt.Errorf("unable to decode %v type", t.Kind)
	}
}

// Decode event data to string
func Decode(event *abi.Event, txLog *ethpb.Log) (map[string]string, error) {
	expectedTopics := len(event.Inputs) - event.Inputs.LengthNonIndexed()
	if expectedTopics != len(txLog.Topics)-1 {
		return nil, fmt.Errorf("decoder error: Topics length does not match with abi event: expected %v but got %v", expectedTopics, len(txLog.Topics)-1)
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

	for i := range event.Inputs {
		var decoded string
		input := event.Inputs[i]

		if input.Indexed {
			decoded, err = FormatIndexedArg(&input.Type, common.HexToHash(txLog.Topics[topicIndex]))
			topicIndex++
		} else {
			decoded, err = FormatNonIndexedArg(&input.Type, unpackValues[unpackValuesIndex])
			unpackValuesIndex++
		}
		if err != nil {
			return nil, err
		}
		logMapping[input.Name] = decoded
	}

	return logMapping, nil
}
