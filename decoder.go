package ethereum

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// FormatIndexedArg transforms a data to string
func FormatIndexedArg(t abi.Type, arg common.Hash) (string, error) {

	switch {
	case t.Type == reflect.TypeOf(&big.Int{}):
		num := new(big.Int).SetBytes(arg[:])
		return fmt.Sprintf("%v", num), nil
	case t.Type == reflect.TypeOf(common.Address{}):
		return common.HexToAddress(arg.Hex()).Hex(), nil
	default:
		switch {
		case t.T == abi.FixedBytesTy:
			return fmt.Sprintf("%v", hexutil.Encode(arg[common.HashLength-t.Type.Size():])), nil
		}
		return fmt.Sprintf("%v", arg), nil
	}
}

// FormatNonIndexedArg transforms a data to string
func FormatNonIndexedArg(t abi.Type, arg interface{}) (string, error) {

	switch v := arg.(type) {
	case common.Address:
		return v.Hex(), nil
	case [1]byte:
		return hexutil.Encode(v[:]), nil
	case [2]byte:
		return hexutil.Encode(v[:]), nil
	case [3]byte:
		return hexutil.Encode(v[:]), nil
	case [4]byte:
		return hexutil.Encode(v[:]), nil
	case [5]byte:
		return hexutil.Encode(v[:]), nil
	case [6]byte:
		return hexutil.Encode(v[:]), nil
	case [7]byte:
		return hexutil.Encode(v[:]), nil
	case [8]byte:
		return hexutil.Encode(v[:]), nil
	case [9]byte:
		return hexutil.Encode(v[:]), nil
	case [10]byte:
		return hexutil.Encode(v[:]), nil
	case [11]byte:
		return hexutil.Encode(v[:]), nil
	case [12]byte:
		return hexutil.Encode(v[:]), nil
	case [13]byte:
		return hexutil.Encode(v[:]), nil
	case [14]byte:
		return hexutil.Encode(v[:]), nil
	case [15]byte:
		return hexutil.Encode(v[:]), nil
	case [16]byte:
		return hexutil.Encode(v[:]), nil
	case [17]byte:
		return hexutil.Encode(v[:]), nil
	case [18]byte:
		return hexutil.Encode(v[:]), nil
	case [19]byte:
		return hexutil.Encode(v[:]), nil
	case [20]byte:
		return hexutil.Encode(v[:]), nil
	case [21]byte:
		return hexutil.Encode(v[:]), nil
	case [22]byte:
		return hexutil.Encode(v[:]), nil
	case [23]byte:
		return hexutil.Encode(v[:]), nil
	case [24]byte:
		return hexutil.Encode(v[:]), nil
	case [25]byte:
		return hexutil.Encode(v[:]), nil
	case [26]byte:
		return hexutil.Encode(v[:]), nil
	case [27]byte:
		return hexutil.Encode(v[:]), nil
	case [28]byte:
		return hexutil.Encode(v[:]), nil
	case [29]byte:
		return hexutil.Encode(v[:]), nil
	case [30]byte:
		return hexutil.Encode(v[:]), nil
	case [31]byte:
		return hexutil.Encode(v[:]), nil
	case [32]byte:
		return hexutil.Encode(v[:]), nil
	case []byte:
		return hexutil.Encode(v[:]), nil
	}
	return fmt.Sprintf("%v", arg), nil
}

// Decode event data to string
func Decode(event *abi.Event, txLog *types.Log) (map[string]string, error) {
	logMapping := make(map[string]string, len(event.Inputs))

	unpackValues, err := event.Inputs.UnpackValues(txLog.Data)
	if err != nil {
		return nil, err
	}

	var (
		topicIndex        = 1
		unpackValuesIndex = 0
	)
	for _, arg := range event.Inputs {
		var decoded string
		if arg.Indexed {
			decoded, _ = FormatIndexedArg(arg.Type, txLog.Topics[topicIndex])
			topicIndex++
		} else {

			decoded, _ = FormatNonIndexedArg(arg.Type, unpackValues[unpackValuesIndex])
			unpackValuesIndex++
		}
		logMapping[arg.Name] = decoded
	}

	return logMapping, nil
}
