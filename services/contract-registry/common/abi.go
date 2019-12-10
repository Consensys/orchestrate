package common

import (
	"encoding/json"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	pkgjson "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const component = "contract-registry-common"

// SigHashToSelector returns the selector associated to a signature hash
func SigHashToSelector(data []byte) (res [4]byte) {
	copy(res[:], data)
	return res
}

// Retrieve fields types & names
type Arguments struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed"`
}

type RawABI struct {
	Type      string      `json:"type,omitempty"`
	Name      string      `json:"name,omitempty"`
	Constant  bool        `json:"constant,omitempty"`
	Anonymous bool        `json:"anonymous,omitempty"`
	Inputs    []Arguments `json:"inputs,omitempty"`
	Outputs   []Arguments `json:"outputs,omitempty"`
}

// ParseJSONABI returns a decoded ABI object
func ParseJSONABI(data []byte) (methods, events map[string][]byte, err error) {
	var parsedFields []RawABI
	err = pkgjson.Unmarshal(data, &parsedFields)
	if err != nil {
		return nil, nil, errors.FromError(err).ExtendComponent(component)
	}

	// Retrieve raw JSONs
	normalizedJSON, err := pkgjson.Marshal(parsedFields)
	if err != nil {
		return nil, nil, errors.FromError(err).ExtendComponent(component)
	}
	var rawFields []json.RawMessage
	err = pkgjson.Unmarshal(normalizedJSON, &rawFields)
	if err != nil {
		return nil, nil, errors.FromError(err).ExtendComponent(component)
	}

	methods = make(map[string][]byte)
	events = make(map[string][]byte)
	for i := 0; i < len(rawFields) && i < len(parsedFields); i++ {
		fieldJSON, err := rawFields[i].MarshalJSON()
		if err != nil {
			return nil, nil, errors.FromError(err).ExtendComponent(component)
		}
		switch parsedFields[i].Type {
		case "function", "":
			var m *ethabi.Method
			err := json.Unmarshal(fieldJSON, &m)
			if err != nil {
				return nil, nil, errors.FromError(err).ExtendComponent(component)
			}
			methods[m.Name+m.Sig()] = fieldJSON
		case "event":
			var e *ethabi.Event
			err := json.Unmarshal(fieldJSON, &e)
			if err != nil {
				return nil, nil, errors.FromError(err).ExtendComponent(component)
			}
			events[e.Name+e.Sig()] = fieldJSON
		}
	}
	return methods, events, nil
}

// GetIndexedCount returns the count of indexed inputs in the event
func GetIndexedCount(event ethabi.Event) (indexedInputCount uint) {
	for i := range event.Inputs {
		if event.Inputs[i].Indexed {
			indexedInputCount++
		}
	}
	return indexedInputCount
}
