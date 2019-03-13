package abi

import (
	"fmt"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

// StaticRegistry stores contract ABI and bytecode in memory
type StaticRegistry struct {
	abis      map[string]*ethabi.ABI
	bytecodes map[string][]byte

	abiMethodBySig map[string]ethabi.Method
	abiEventBySig  map[string]ethabi.Event
}

// NewStaticRegistry initialise a newly created ContractABIRegistry
func NewStaticRegistry() *StaticRegistry {
	return &StaticRegistry{
		abis:           make(map[string]*ethabi.ABI),
		bytecodes:      make(map[string][]byte),
		abiMethodBySig: make(map[string]ethabi.Method),
		abiEventBySig:  make(map[string]ethabi.Event),
	}
}

// RegisterContract allow to add a contract and its corresponding ABI to the registry
func (r *StaticRegistry) RegisterContract(contract *abi.Contract) error {
	abi, err := contract.ToABI()
	if err != nil {
		return err
	}

	// Register ABI and bytecode
	r.abis[contract.Short()] = abi
	r.bytecodes[contract.Short()] = contract.Bytecode

	// TODO differentiate registering vs updating
	for _, method := range r.abis[contract.Short()].Methods {
		r.abiMethodBySig[hexutil.Encode(method.Id())] = method
	}

	for _, event := range r.abis[contract.Short()].Events {
		r.abiEventBySig[event.Id().Hex()] = event
	}

	return nil
}

func (r *StaticRegistry) getContract(name string) (*ethabi.ABI, error) {
	abi, ok := r.abis[name]
	if !ok {
		return nil, fmt.Errorf("Unknown contract %q", name)
	}
	return abi, nil
}

// GetMethodByID returns the abi for a given method of a contract
// id should match the following pattern "<MethodName>@<ContracName>"
func (r *StaticRegistry) GetMethodByID(id string) (ethabi.Method, error) {
	// Computing call ensure ID has been properly formated
	call, err := common.StringToCall(id)
	if err != nil {
		return ethabi.Method{}, err
	}

	// Retrieve contract ABI from registry
	abi, err := r.getContract(call.GetContract().Short())
	if err != nil {
		return ethabi.Method{}, err
	}

	// If call is a deployment we return constructor
	if call.IsDeploy() {
		return abi.Constructor, nil
	}

	method, ok := abi.Methods[call.GetMethod().GetName()]
	if !ok {
		return ethabi.Method{}, fmt.Errorf("Contract %q has no method %q", call.GetContract().Short(), call.GetMethod().GetName())
	}

	return method, nil
}

// GetEventByID returns the abi for a given event of a contract
// id should match the following pattern "<EventName>@<ContracName>"
func (r *StaticRegistry) GetEventByID(id string) (ethabi.Event, error) {
	// Computing call ensure ID has been properly formated
	call, err := common.StringToCall(id)
	if err != nil {
		return ethabi.Event{}, err
	}

	abi, err := r.getContract(call.GetContract().Short())
	if err != nil {
		return ethabi.Event{}, err
	}

	event, exist := abi.Events[call.GetMethod().GetName()]
	if !exist {
		return ethabi.Event{}, fmt.Errorf("Contract %q has no event %q", call.GetContract().Short(), call.GetMethod().GetName())
	}

	return event, nil
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// GetMethodBySig returns the method corresponding to input signature
// The input signature should be in hex format (matching the regex patterns "0x[0-9a-f]{8}" or "[0-9a-f]{8}")
func (r *StaticRegistry) GetMethodBySig(sig string) (ethabi.Method, error) {
	if !has0xPrefix(sig) {
		sig = fmt.Sprintf("0x%v", sig)
	}

	bytesig, err := hexutil.Decode(sig)
	if err != nil {
		return ethabi.Method{}, err
	}

	method, ok := r.abiMethodBySig[hexutil.Encode(bytesig)]
	if !ok {
		return ethabi.Method{}, fmt.Errorf("No method with signature %v", sig)
	}
	return method, nil
}

// GetEventBySig returns the event corresponding to input signature
// The input signature should be in hex format (matching the regex patterns "0x[0-9a-f]{16}" or "[0-9a-f]{16}")
func (r *StaticRegistry) GetEventBySig(topic string) (ethabi.Event, error) {
	if !has0xPrefix(topic) {
		topic = fmt.Sprintf("0x%v", topic)
	}

	bytetopic, err := hexutil.Decode(topic)
	if err != nil {
		return ethabi.Event{}, err
	}

	event, exist := r.abiEventBySig[ethcommon.BytesToHash(bytetopic).Hex()]
	if !exist {
		return ethabi.Event{}, fmt.Errorf("No event with topic %v", topic)
	}
	return event, nil
}
