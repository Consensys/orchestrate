package parsers

import (
	"encoding/json"
	"regexp"

	pkgjson "github.com/ConsenSys/orchestrate/pkg/encoding/json"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var contractRegexp = `^(?P<contract>[a-zA-Z0-9]+)(?:\[(?P<tag>[0-9a-zA-Z-.]+)\])?(?::(?P<abi>\[.+\]))?(?::(?P<bytecode>0[xX][a-fA-F0-9]+))?(?::(?P<deployedBytecode>0[xX][a-fA-F0-9]+))?$`
var contractPattern = regexp.MustCompile(contractRegexp)

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

// StringToContract computes a Contract from is short representation
func StringToContract(s string) (*entities.Contract, error) {
	parts := contractPattern.FindStringSubmatch(s)

	if len(parts) != 6 {
		return nil, errors.InvalidFormatError("invalid contract (expected format %s) %q", contractRegexp, s)
	}

	c := &entities.Contract{
		Name: parts[1],
		Tag:  parts[2],
	}

	// Make sure bytecode is valid and set bytecode
	if parts[4] == "" {
		parts[4] = "0x"
	}
	_, err := hexutil.Decode(parts[4])
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract bytecode on %q", c.Short())
	}
	c.Bytecode = parts[4]

	// Make sure deployedBytecode is valid and set deployedBytecode
	if parts[5] == "" {
		parts[5] = "0x"
	}
	_, err = hexutil.Decode(parts[5])
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract deployed bytecode on %q", c.Short())
	}
	c.DeployedBytecode = parts[5]

	// Set ABI and make sure it is valid
	c.ABI = parts[3]
	_, err = c.ToABI()
	if err != nil {
		return nil, errors.InvalidFormatError("invalid contract ABI on %q", c.Short())
	}

	return c, nil
}
