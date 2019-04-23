package static

import (
	"fmt"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

// Registry stores contract ABI and bytecode in memory
type Registry struct {
	abis      map[string]*ethabi.ABI
	bytecodes map[string][]byte

	abiMethodBySig map[string]ethabi.Method
	abiEventBySig  map[string]ethabi.Event
}

// NewRegistry creates a New Registry
func NewRegistry() *Registry {
	return &Registry{
		abis:           make(map[string]*ethabi.ABI),
		bytecodes:      make(map[string][]byte),
		abiMethodBySig: make(map[string]ethabi.Method),
		abiEventBySig:  make(map[string]ethabi.Event),
	}
}

// RegisterContract allow to add a contract and its corresponding ABI to the registry
func (r *Registry) RegisterContract(contract *abi.Contract) error {
	contractAbi, err := contract.ToABI()
	if err != nil {
		return err
	}

	// Register ABI and bytecode
	r.abis[contract.Short()] = contractAbi
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

func (r *Registry) getContract(name string) (*ethabi.ABI, error) {
	contractAbi, ok := r.abis[name]
	if !ok {
		return nil, fmt.Errorf("unknown contract %q", name)
	}
	return contractAbi, nil
}

// GetMethodByID returns the abi for a given method of a contract
// id should match the following pattern "<MethodName>@<ContracName>"
func (r *Registry) GetMethodByID(id string) (ethabi.Method, error) {
	// Computing call ensure ID has been properly formated
	call, err := common.StringToCall(id)
	if err != nil {
		return ethabi.Method{}, err
	}

	// Retrieve contract ABI from registry
	contractAbi, err := r.getContract(call.GetContract().Short())
	if err != nil {
		return ethabi.Method{}, err
	}

	// If call is a deployment we return constructor
	if call.IsDeploy() {
		return contractAbi.Constructor, nil
	}

	method, ok := contractAbi.Methods[call.GetMethod().GetName()]
	if !ok {
		return ethabi.Method{}, fmt.Errorf("contract %q has no method %q", call.GetContract().Short(), call.GetMethod().GetName())
	}

	return method, nil
}

// GetEventByID returns the abi for a given event of a contract
// id should match the following pattern "<EventName>@<ContracName>"
func (r *Registry) GetEventByID(id string) (ethabi.Event, error) {
	// Computing call ensure ID has been properly formated
	call, err := common.StringToCall(id)
	if err != nil {
		return ethabi.Event{}, err
	}

	contractAbi, err := r.getContract(call.GetContract().Short())
	if err != nil {
		return ethabi.Event{}, err
	}

	event, exist := contractAbi.Events[call.GetMethod().GetName()]
	if !exist {
		return ethabi.Event{}, fmt.Errorf("contract %q has no event %q", call.GetContract().Short(), call.GetMethod().GetName())
	}

	return event, nil
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// GetMethodBySig returns the method corresponding to input signature
// The input signature should be in hex format (matching the regex patterns "0x[0-9a-f]{8}" or "[0-9a-f]{8}")
func (r *Registry) GetMethodBySig(sig string) (ethabi.Method, error) {
	if !has0xPrefix(sig) {
		sig = fmt.Sprintf("0x%v", sig)
	}

	bytesig, err := hexutil.Decode(sig)
	if err != nil {
		return ethabi.Method{}, err
	}

	method, ok := r.abiMethodBySig[hexutil.Encode(bytesig)]
	if !ok {
		return ethabi.Method{}, fmt.Errorf("no method with signature %v", sig)
	}
	return method, nil
}

// GetEventBySig returns the event corresponding to input signature
// The input signature should be in hex format (matching the regex patterns "0x[0-9a-f]{16}" or "[0-9a-f]{16}")
func (r *Registry) GetEventBySig(topic string) (ethabi.Event, error) {
	if !has0xPrefix(topic) {
		topic = fmt.Sprintf("0x%v", topic)
	}

	bytetopic, err := hexutil.Decode(topic)
	if err != nil {
		return ethabi.Event{}, err
	}

	event, exist := r.abiEventBySig[ethcommon.BytesToHash(bytetopic).Hex()]
	if !exist {
		return ethabi.Event{}, fmt.Errorf("no event with topic %v", topic)
	}
	return event, nil
}

// GetBytecodeByID returns the bytecode of the contract
func (r *Registry) GetBytecodeByID(id string) (code []byte, err error) {
	// Computing call ensure ID has been properly formated
	call, err := common.StringToCall(id)
	if err != nil {
		return []byte{}, err
	}

	res, err := r.getBytecode(call.GetContract().Short())
	if err != nil {
		return []byte{}, err
	}

	return res, nil
}

// getBytecode is a low-level getter for the bytecode
func (r *Registry) getBytecode(name string) ([]byte, error) {
	code, ok := r.bytecodes[name]
	if !ok {
		return nil, fmt.Errorf("unknown contract %q", name)
	}
	return code, nil
}
