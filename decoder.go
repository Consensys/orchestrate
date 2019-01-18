package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// EventDecoder ...
type EventDecoder struct {
	Inputs abi.Arguments
}

// FormatIndexed transforms a data to string
func FormatIndexedEvent(datatype string, data string) (string, error) {
	switch {
	case datatype == "address":
		return common.HexToAddress(data).Hex(), nil
	default:
		return fmt.Sprintf("%v", data), nil
	}
}

// FormatData transforms a data to string
func FormatNonIndexEvent(t abi.Type, data interface{}) (string, error) {

	switch v := data.(type) {
	case common.Address:
		return v.Hex(), nil
	case [8]byte:
	case [16]byte:
	case [32]byte:
		return hexutil.Encode(v[:]), nil
	}
	return fmt.Sprintf("%v", data), nil
}

// Decode event data to string
func (event *EventDecoder) Decode(txLog *types.Log) (map[string]string, error) {
	logMapping := make(map[string]string, len(event.Inputs))

	unpackValues, err := event.Inputs.UnpackValues(txLog.Data)
	if err != nil {
		return nil, err
	}

	var topicIndex = 1
	var unpackValuesIndex = 0
	for _, arg := range event.Inputs {
		var decoded string
		if arg.Indexed {
			decoded, _ = FormatIndexedEvent(fmt.Sprintf("%s", arg.Type), txLog.Topics[topicIndex].Hex())
			topicIndex++
		} else {
			decoded, _ = FormatNonIndexEvent(arg.Type, unpackValues[unpackValuesIndex])
			unpackValuesIndex++
		}
		logMapping[arg.Name] = decoded
	}

	return logMapping, nil
}
