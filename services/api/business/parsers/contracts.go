package parsers

import (
	"encoding/json"

	ethabi "github.com/consensys/orchestrate/pkg/ethereum/abi"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

// TODO: Remove this function as parsing the events from the ABI should not be done on Orchestrate as we do not have control on how the events are represented in the ABI

// ParseEvents returns a map of events given an ABI
func ParseEvents(data string) (map[string]string, error) {
	var parsedFields []entities.RawABI

	err := json.Unmarshal([]byte(data), &parsedFields)
	if err != nil {
		return nil, err
	}

	// Retrieve raw JSONs
	normalizedJSON, err := json.Marshal(parsedFields)
	if err != nil {
		return nil, err
	}

	var rawFields []json.RawMessage
	err = json.Unmarshal(normalizedJSON, &rawFields)
	if err != nil {
		return nil, err
	}

	events := make(map[string]string)
	for i := 0; i < len(rawFields) && i < len(parsedFields); i++ {
		fieldJSON, err := rawFields[i].MarshalJSON()
		if err != nil {
			return nil, err
		}

		if parsedFields[i].Type == "event" {
			e := &ethabi.Event{}
			err := json.Unmarshal(fieldJSON, e)
			if err != nil {
				return nil, err
			}

			events[e.Name+e.Sig()] = string(fieldJSON)
		}
	}

	return events, nil
}
