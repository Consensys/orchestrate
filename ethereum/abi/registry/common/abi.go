package common

import (
	"encoding/json"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	pkgJson "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

const component = "registry-common"

// SigHashToSelector returns the selector associated to a signature hash
func SigHashToSelector(data []byte) (res [4]byte) {
	copy(res[:], data)
	return res
}

// ParseJSONABI returns a decoded ABI object
func ParseJSONABI(data []byte) (methods, events map[string][]byte, err error) {
	// Retrieve fields types & names
	type Arguments struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Indexed bool   `json:"indexed"`
	}

	var parsedFields []struct {
		Type      string      `json:"type"`
		Name      string      `json:"name"`
		Constant  bool        `json:"constant"`
		Anonymous bool        `json:"anonymous"`
		Inputs    []Arguments `json:"inputs"`
		Outputs   []Arguments `json:"outputs"`
	}
	err = pkgJson.Unmarshal(data, &parsedFields)
	if err != nil {
		return nil, nil, errors.FromError(err).ExtendComponent(component)
	}

	// Retrieve raw JSONs
	normalizedJSON, err := pkgJson.Marshal(parsedFields)
	if err != nil {
		return nil, nil, errors.FromError(err).ExtendComponent(component)
	}
	var rawFields []json.RawMessage
	err = pkgJson.Unmarshal(normalizedJSON, &rawFields)
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
			methods[parsedFields[i].Name] = fieldJSON
		case "event":
			events[parsedFields[i].Name] = fieldJSON
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
