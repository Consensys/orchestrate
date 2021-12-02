package parsers

import (
	"encoding/json"

	pkgjson "github.com/consensys/orchestrate/pkg/encoding/json"
	ethabi "github.com/consensys/orchestrate/pkg/ethereum/abi"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

// ParseJSONABI returns a decoded ABI object
func ParseJSONABI(data string) (methods, events map[string]string, err error) {
	var parsedFields []entities.RawABI
	err = pkgjson.Unmarshal([]byte(data), &parsedFields)
	if err != nil {
		return nil, nil, err
	}

	// Retrieve raw JSONs
	normalizedJSON, err := pkgjson.Marshal(parsedFields)
	if err != nil {
		return nil, nil, err
	}
	var rawFields []json.RawMessage
	err = pkgjson.Unmarshal(normalizedJSON, &rawFields)
	if err != nil {
		return nil, nil, err
	}

	methods = make(map[string]string)
	events = make(map[string]string)
	for i := 0; i < len(rawFields) && i < len(parsedFields); i++ {
		fieldJSON, err := rawFields[i].MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		switch parsedFields[i].Type {
		case "function", "":
			var m *ethabi.Method
			err := json.Unmarshal(fieldJSON, &m)
			if err != nil {
				return nil, nil, err
			}
			methods[m.Name+m.Sig()] = string(fieldJSON)
		case "event":
			var e *ethabi.Event
			err := json.Unmarshal(fieldJSON, &e)
			if err != nil {
				return nil, nil, err
			}
			events[e.Name+e.Sig()] = string(fieldJSON)
		}
	}
	return methods, events, nil
}
