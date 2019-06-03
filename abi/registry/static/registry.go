package static

import (
	"context"
	"fmt"
	"reflect"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

// Registry stores contract ABI and bytecode in memory
type Registry struct {
	ethClient ethclient.ChainStateReader
	// Contract registry/name#tag to bytecode hash
	contractHash map[string]ethCommon.Hash
	// Bytecode hash to ABI, bytecode and deployed bytecode
	contracts map[ethCommon.Hash]*abi.Contract

	// Address to Codehash (deployed bytecode hash) map
	addressCodehash map[string]map[ethCommon.Address]ethCommon.Hash

	// Codehash to Selector to ABIs
	methods map[ethCommon.Hash]map[[4]byte][]*ethAbi.Method
	events  map[ethCommon.Hash]map[ethCommon.Hash]map[uint][]*ethAbi.Event
}

var defaultCodehash = ethCommon.Hash{}

// NewRegistry creates a New Registry
func NewRegistry(client ethclient.ChainStateReader) *Registry {
	r := &Registry{
		ethClient:       client,
		contractHash:    make(map[string]ethCommon.Hash),
		contracts:       make(map[ethCommon.Hash]*abi.Contract),
		addressCodehash: make(map[string]map[ethCommon.Address]ethCommon.Hash),
		methods:         make(map[ethCommon.Hash]map[[4]byte][]*ethAbi.Method),
		events:          make(map[ethCommon.Hash]map[ethCommon.Hash]map[uint][]*ethAbi.Event),
	}
	r.methods[defaultCodehash] = make(map[[4]byte][]*ethAbi.Method)
	r.events[defaultCodehash] = make(map[ethCommon.Hash]map[uint][]*ethAbi.Event)
	return r
}

// RegisterContract allow to add a contract and its corresponding ABI to the registry
func (r *Registry) RegisterContract(contract *abi.Contract) error {
	if contract.Bytecode != nil {
		bytecodeHash := crypto.Keccak256Hash(contract.Bytecode)
		r.contractHash[contract.Short()] = bytecodeHash

		r.contracts[bytecodeHash] = &abi.Contract{
			Abi:              contract.Abi,
			Bytecode:         contract.Bytecode,
			DeployedBytecode: contract.DeployedBytecode,
		}
	}

	codeHash := crypto.Keccak256Hash(contract.DeployedBytecode)
	contractAbi, err := contract.ToABI()
	if err != nil {
		return fmt.Errorf("registry: could not register contract, wrong ABI format %v", err)
	}

	for _, method := range contractAbi.Methods {
		var id [4]byte
		copy(id[:], method.Id())
		if contract.DeployedBytecode != nil {
			// Init map
			if r.methods[codeHash] == nil {
				r.methods[codeHash] = make(map[[4]byte][]*ethAbi.Method)
			}

			r.methods[codeHash][id] = []*ethAbi.Method{&method}
		}

		// Register in default methods if not present
		found := false
		for _, m := range r.methods[defaultCodehash][id] {
			if reflect.DeepEqual(m, &method) {
				found = true
			}
		}
		if !found {
			r.methods[defaultCodehash][id] = append(r.methods[defaultCodehash][id], &method)
		}
	}

	for _, event := range contractAbi.Events {
		indexedCount := getIndexedCount(event)

		if contract.DeployedBytecode != nil {
			// Init map
			if r.events[codeHash] == nil {
				r.events[codeHash] = make(map[ethCommon.Hash]map[uint][]*ethAbi.Event)
			}
			// Init map
			if r.events[codeHash][event.Id()] == nil {
				r.events[codeHash][event.Id()] = make(map[uint][]*ethAbi.Event)
			}

			r.events[codeHash][event.Id()][indexedCount] = []*ethAbi.Event{&event}
		}

		// Init map
		if r.events[defaultCodehash][event.Id()] == nil {
			r.events[defaultCodehash][event.Id()] = make(map[uint][]*ethAbi.Event)
		}
		// Register in default events if not present
		found := false
		for _, e := range r.events[defaultCodehash][event.Id()][indexedCount] {
			if reflect.DeepEqual(e, &event) {
				found = true
			}
		}
		if !found {
			r.events[defaultCodehash][event.Id()][indexedCount] = append(
				r.events[defaultCodehash][event.Id()][indexedCount],
				&event,
			)
		}
	}

	return nil
}

// Retrieve contract ABI
func (r *Registry) GetContractABI(contract *abi.Contract) ([]byte, error) {
	bytecodeHash := r.contractHash[contract.Short()]
	c, ok := r.contracts[bytecodeHash]
	if !ok {
		return nil, fmt.Errorf("registry: could not find contract")
	}
	return c.Abi, nil
}

// Returns the bytecode
func (r *Registry) GetContractBytecode(contract *abi.Contract) ([]byte, error) {
	bytecodeHash := r.contractHash[contract.Short()]
	c, ok := r.contracts[bytecodeHash]
	if !ok {
		return nil, fmt.Errorf("registry: could not find contract")
	}
	return c.Bytecode, nil
}

// Returns the deployed bytecode
func (r *Registry) GetContractDeployedBytecode(contract *abi.Contract) ([]byte, error) {
	bytecodeHash := r.contractHash[contract.Short()]
	c, ok := r.contracts[bytecodeHash]
	if !ok {
		return nil, fmt.Errorf("registry: could not find contract")
	}
	return c.DeployedBytecode, nil
}

// getIndexedCount returns the count of indexed inputs in the event
func getIndexedCount(event ethAbi.Event) uint {
	var indexedInputCount uint
	for i := range event.Inputs {
		if event.Inputs[i].Indexed {
			indexedInputCount++
		}
	}
	return indexedInputCount
}

// Get the codehash of a contract instance
func (r *Registry) getCodehash(contract common.AccountInstance) (ethCommon.Hash, error) {
	codehashToAddressMap, ok := r.addressCodehash[contract.GetChain().String()]
	if !ok {
		return ethCommon.Hash{}, fmt.Errorf("registry: could not find contract: bad chainid")
	}
	address, err := contract.GetAccount().Address()
	if err != nil {
		return ethCommon.Hash{}, fmt.Errorf("registry: could not find contract: %v", err)
	}
	codehash, ok := codehashToAddressMap[address]
	if !ok {
		return ethCommon.Hash{}, fmt.Errorf("registry: could not find contract: bad address")
	}
	return codehash, nil
}

// Retrieve method using 4 bytes unique selector and the address of the contract
func (r *Registry) GetMethodsBySelector(selector [4]byte, contract common.AccountInstance) (method *ethAbi.Method, defaultMethods []*ethAbi.Method, err error) {
	// Search in specific method storage
	contractCodehash, err := r.getCodehash(contract)
	if err == nil {
		contractMethods, ok := r.methods[contractCodehash][selector]
		if ok && len(contractMethods) == 1 {
			return contractMethods[0], nil, nil
		}
	}

	// Search in default methods
	defaultMethods, ok := r.methods[defaultCodehash][selector]
	if ok {
		return nil, defaultMethods, nil
	}

	return nil, nil, fmt.Errorf("registry: could not find corresponding method ABIs")
}

// Retrieve event using 4 bytes unique selector
func (r *Registry) GetEventsBySelector(selector ethCommon.Hash, contract common.AccountInstance, indexedInputCount uint) (event *ethAbi.Event, defaultEvents []*ethAbi.Event, err error) {
	// Search in specific event storage
	contractCodehash, err := r.getCodehash(contract)
	if err == nil {
		contractEvents, ok := r.events[contractCodehash][selector]
		if ok {
			matchingContractEvents, ok := contractEvents[indexedInputCount]
			if ok && len(matchingContractEvents) == 1 {
				return matchingContractEvents[0], nil, nil
			}
		}
	}

	// Search in default events
	if defaultEvents, ok := r.events[defaultCodehash][selector][indexedInputCount]; ok {
		return nil, defaultEvents, nil
	}

	return nil, nil, fmt.Errorf("registry: no event match found, no default and can't find contract: %v", err)
}

// Request an update of the codehash of the contract address
func (r *Registry) RequestAddressUpdate(contract common.AccountInstance) error {
	addr, err := contract.GetAccount().Address()
	if err != nil {
		return fmt.Errorf("registry: could not update address: address invalid: %v", err)
	}

	// Codehash already stored for this contract instance
	if _, ok := r.addressCodehash[contract.GetChain().String()][addr]; ok {
		return nil
	}

	// Codehash not stored, trying to retrieve it from chain
	code, err := r.ethClient.CodeAt(context.Background(), contract.GetChain().ID(), addr, nil)
	if err != nil {
		return fmt.Errorf("registry: could not update address: client error: %v", err)
	}
	codehash := crypto.Keccak256Hash(code)
	chain := contract.GetChain().String()
	if r.addressCodehash[chain] == nil {
		r.addressCodehash[chain] = make(map[ethCommon.Address]ethCommon.Hash)
	}
	r.addressCodehash[chain][addr] = codehash

	return nil
}
