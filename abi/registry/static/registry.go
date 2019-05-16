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

	abiMethods map[string]ethabi.Method
	abiEvents  map[string]ethabi.Event
}

// NewRegistry creates a New Registry
func NewRegistry() *Registry {
	return &Registry{
		abis:       make(map[string]*ethabi.ABI),
		bytecodes:  make(map[string][]byte),
		abiMethods: make(map[string]ethabi.Method),
		abiEvents:  make(map[string]ethabi.Event),
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
		r.abiMethods[hexutil.Encode(method.Id())] = method
	}

	for _, event := range r.abis[contract.Short()].Events {
		r.abiEvents[event.Id().Hex()] = event
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

// GetMethodBySig returns the abi for a given method of a contract
// sig should match the following pattern "func(type1,type2)"
func (r *Registry) GetMethodBySig(contract, sig string) (*ethabi.Method, error) {
	// Computing call ensure sig has been properly formated
	call, err := common.SignatureToCall(sig)
	if err != nil {
		return nil, err
	}

	// Retrieve contract ABI from registry
	contractAbi, err := r.getContract(contract)
	if err != nil {
		return nil, err
	}

	// If call is a deployment we return constructor
	if call.IsConstructor() {
		return &contractAbi.Constructor, nil
	}

	method, ok := contractAbi.Methods[call.GetMethod().GetName()]
	if !ok {
		return nil, fmt.Errorf("contract %q has no method %q", contract, call.GetMethod().GetName())
	}

	return &method, nil
}

// GetEventBySig returns the abi for a given event of a contract
// sig should match the following pattern "event(type1,type2)"
func (r *Registry) GetEventBySig(contract, sig string) (*ethabi.Event, error) {
	// Computing call ensure sig has been properly formated
	call, err := common.SignatureToCall(sig)
	if err != nil {
		return nil, err
	}

	contractAbi, err := r.getContract(contract)
	if err != nil {
		return nil, err
	}

	event := &ethabi.Event{}
	var ok bool
	*event, ok = contractAbi.Events[call.GetMethod().GetName()]
	if !ok {
		return nil, fmt.Errorf("contract %q has no event %q", call.GetContract().Short(), call.GetMethod().GetName())
	}

	return event, nil
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// GetMethodBySelector returns the method corresponding to input signature
// The input signature should be in hex format (matching the regex patterns "0x[0-9a-f]{8}" or "[0-9a-f]{8}")
func (r *Registry) GetMethodBySelector(selector string) (*ethabi.Method, error) {
	if !has0xPrefix(selector) {
		selector = fmt.Sprintf("0x%v", selector)
	}

	bytesig, err := hexutil.Decode(selector)
	if err != nil {
		return nil, err
	}

	method, ok := r.abiMethods[hexutil.Encode(bytesig)]
	if !ok {
		return nil, fmt.Errorf("no method with signature %v", selector)
	}
	return &method, nil
}

// GetEventBySelector returns the event corresponding to input signature
// The input signature should be in hex format (matching the regex patterns "0x[0-9a-f]{16}" or "[0-9a-f]{16}")
func (r *Registry) GetEventBySelector(selector string) (*ethabi.Event, error) {
	if !has0xPrefix(selector) {
		selector = fmt.Sprintf("0x%v", selector)
	}

	byteselector, err := hexutil.Decode(selector)
	if err != nil {
		return &ethabi.Event{}, err
	}

	event, ok := r.abiEvents[ethcommon.BytesToHash(byteselector).Hex()]
	if !ok {
		return &ethabi.Event{}, fmt.Errorf("no event with topic %v", selector)
	}
	return &event, nil
}

// GetBytecodeByID returns the bytecode of the contract
func (r *Registry) GetBytecodeByID(id string) (code []byte, err error) {
	// Computing call ensure ID has been properly formated
	call, err := common.SignatureToCall(id)
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
