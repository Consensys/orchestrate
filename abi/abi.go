package abi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ContractABIRegistry will contain all abis
type ContractABIRegistry struct {
	// TODO handle contract version
	abis           map[string]*abi.ABI
	abiMethodBySig map[string]abi.Method
	abiEventBySig  map[string]abi.Event
}

// RegisterContract allow to add a contract and its corresponding ABI to the registry
func (r *ContractABIRegistry) RegisterContract(name string, contractABI []byte) error {
	r.abis[name] = &abi.ABI{}
	// TODO differentiate registering vs updating
	err := r.abis[name].UnmarshalJSON(contractABI)
	if err != nil {
		return err
	}
	for _, method := range r.abis[name].Methods {
		sig := "0x" + hex.EncodeToString(method.Id())
		r.abiMethodBySig[sig] = method
	}
	for _, event := range r.abis[name].Events {
		sig := event.Id().Hex()
		r.abiEventBySig[sig] = event
	}
	return nil
}

func parseID(id string) (string, string, error) {
	s := strings.Split(id, "@")
	if len(s) != 2 {
		return "", "", errors.New("id input does not match the <FunctionOrEventName>@<ContractName> pattern")
	}
	return s[0], s[1], nil
}

func cleanSig(sig string, sigType string) (string, error) {
	var pattern string
	switch sigType {
	case "method":
		pattern = "^(0x)?[0-9a-fA-F]{8}$"
	case "event":
		pattern = "^(0x)?[0-9a-fA-F]{64}$"
	default:
		return "", fmt.Errorf("Cannot use this sigType %v", sigType)
	}
	matched, err := regexp.MatchString(pattern, sig)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("%v is not a valid %v signature, it should match the regex pattern \"%v\"", sig, sigType, pattern)
	}
	s := strings.ToLower(sig)
	if s[:2] == "0x" {
		return s, nil
	}
	return "0x" + s, nil
}

func (r *ContractABIRegistry) getContractABI(contractName string) (*abi.ABI, error) {
	contractABI, exist := r.abis[contractName]
	if !exist {
		return &abi.ABI{}, fmt.Errorf("Could not find contract %v in ABI registry", contractName)
	}
	return contractABI, nil
}

// GetMethodByID returns the abi for a given method of a contract
// id should match the following pattern "<MethodName>@<ContracName>"
func (r *ContractABIRegistry) GetMethodByID(id string) (abi.Method, error) {
	method, contract, err := parseID(id)
	if err != nil {
		return abi.Method{}, err
	}
	contractABI, err := r.getContractABI(contract)
	if err != nil {
		return abi.Method{}, err
	}
	// TODO handle constructor
	methodABI, exist := contractABI.Methods[method]
	if !exist {
		return abi.Method{}, fmt.Errorf("Could not find method %v in contract %v ABI", method, contract)
	}
	return methodABI, nil
}

// GetEventByID returns the abi for a given event of a contract
// id should match the following pattern "<EventName>@<ContracName>"
func (r *ContractABIRegistry) GetEventByID(id string) (abi.Event, error) {
	event, contract, err := parseID(id)
	if err != nil {
		return abi.Event{}, err
	}
	contractABI, err := r.getContractABI(contract)
	if err != nil {
		return abi.Event{}, err
	}
	eventABI, exist := contractABI.Events[event]
	if !exist {
		return abi.Event{}, fmt.Errorf("Could not find event %v in contract %v ABI", event, contract)
	}
	return eventABI, nil
}

// GetMethodBySig returns the method corresponding to input signature
// The input signature should be in hex format (matching the regex patterns "0x[0-9a-f]{8}" or "[0-9a-f]{8}")
func (r *ContractABIRegistry) GetMethodBySig(sig string) (abi.Method, error) {
	s, err := cleanSig(sig, "method")
	if err != nil {
		return abi.Method{}, err
	}
	methodABI, exist := r.abiMethodBySig[s]
	if !exist {
		return abi.Method{}, fmt.Errorf("Could not find method signature %v in the registry", sig)
	}
	return methodABI, nil
}

// GetEventBySig returns the event corresponding to input signature
// The input signature should be in hex format (matching the regex patterns "0x[0-9a-f]{16}" or "[0-9a-f]{16}")
func (r *ContractABIRegistry) GetEventBySig(sig string) (abi.Event, error) {
	s, err := cleanSig(sig, "event")
	if err != nil {
		return abi.Event{}, err
	}
	eventABI, exist := r.abiEventBySig[s]
	if !exist {
		return abi.Event{}, fmt.Errorf("Could not find event signature %v in the registry", sig)
	}
	return eventABI, nil
}

// NewContractABIRegistry initialise a newly created ContractABIRegistry
func NewContractABIRegistry() *ContractABIRegistry {
	r := &ContractABIRegistry{}
	r.abis = make(map[string]*abi.ABI)
	r.abiMethodBySig = make(map[string]abi.Method)
	r.abiEventBySig = make(map[string]abi.Event)
	return r
}
