package mock

import (
	"context"
	"reflect"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

// ContractRegistry stores contract ABI and bytecode in memory
type ContractRegistry struct {
	ec ethclient.ChainStateReader
	// Contract registry/name#tag to bytecode hash
	contractHashes map[string]ethcommon.Hash

	// Bytecode hash to ABI, bytecode and deployed bytecode
	contracts map[ethcommon.Hash]*abi.Contract

	// Address to Codehash (deployed bytecode hash) map
	codehashes map[string]map[ethcommon.Address]ethcommon.Hash

	// Codehash to Selector to method ABIs
	methods map[ethcommon.Hash]map[[4]byte][]*ethabi.Method

	// Codehash to SigHash to event ABIs
	events map[ethcommon.Hash]map[ethcommon.Hash]map[uint][]*ethabi.Event
}

var defaultCodehash = ethcommon.Hash{}

// NewRegistry creates a ContractRegistry
func NewRegistry(client ethclient.ChainStateReader) *ContractRegistry {
	r := &ContractRegistry{
		ec:             client,
		contractHashes: make(map[string]ethcommon.Hash),
		contracts:      make(map[ethcommon.Hash]*abi.Contract),
		codehashes:     make(map[string]map[ethcommon.Address]ethcommon.Hash),
		methods:        make(map[ethcommon.Hash]map[[4]byte][]*ethabi.Method),
		events:         make(map[ethcommon.Hash]map[ethcommon.Hash]map[uint][]*ethabi.Event),
	}
	r.methods[defaultCodehash] = make(map[[4]byte][]*ethabi.Method)
	r.events[defaultCodehash] = make(map[ethcommon.Hash]map[uint][]*ethabi.Event)
	return r
}

// RegisterContract register a contract including ABI, bytecode and deployed bytecode
func (r *ContractRegistry) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest) (*svc.RegisterContractResponse, error) {
	contract := req.GetContract()

	if contract.Bytecode != nil {
		bytecodeHash := crypto.Keccak256Hash(contract.Bytecode)
		r.contractHashes[contract.Short()] = bytecodeHash

		r.contracts[bytecodeHash] = &abi.Contract{
			Abi:              contract.Abi,
			Bytecode:         contract.Bytecode,
			DeployedBytecode: contract.DeployedBytecode,
		}
	}

	codeHash := crypto.Keccak256Hash(contract.DeployedBytecode)
	contractAbi, err := contract.ToABI()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	for _, m := range contractAbi.Methods {
		method := m
		var id [4]byte
		copy(id[:], method.Id())
		if contract.DeployedBytecode != nil {
			// Init map
			if r.methods[codeHash] == nil {
				r.methods[codeHash] = make(map[[4]byte][]*ethabi.Method)
			}

			r.methods[codeHash][id] = []*ethabi.Method{&method}
		}

		// Register in default methods if not present
		found := false
		for _, m := range r.methods[defaultCodehash][id] {
			if reflect.DeepEqual(*m, method) {
				found = true
			}
		}
		if !found {
			r.methods[defaultCodehash][id] = append(r.methods[defaultCodehash][id], &method)
		}
	}

	for _, e := range contractAbi.Events {
		event := e
		indexedCount := getIndexedCount(event)

		if contract.DeployedBytecode != nil {
			// Init map
			if r.events[codeHash] == nil {
				r.events[codeHash] = make(map[ethcommon.Hash]map[uint][]*ethabi.Event)
			}
			// Init map
			if r.events[codeHash][event.Id()] == nil {
				r.events[codeHash][event.Id()] = make(map[uint][]*ethabi.Event)
			}

			r.events[codeHash][event.Id()][indexedCount] = []*ethabi.Event{&event}
		}

		// Init map
		if r.events[defaultCodehash][event.Id()] == nil {
			r.events[defaultCodehash][event.Id()] = make(map[uint][]*ethabi.Event)
		}
		// Register in default events if not present
		found := false
		for _, e := range r.events[defaultCodehash][event.Id()][indexedCount] {
			if reflect.DeepEqual(*e, event) {
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

	return &svc.RegisterContractResponse{}, nil
}

func (r *ContractRegistry) getContract(c *abi.Contract) (contract *abi.Contract, ok bool) {
	contract, ok = r.contracts[r.contractHashes[c.Short()]]
	return
}

// GetContractABI loads contract ABI
func (r *ContractRegistry) GetContractABI(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractABIResponse, error) {
	c, ok := r.getContract(req.GetContract())
	if !ok {
		return nil, errors.NotFoundError("contract ABI not found").SetComponent(component)
	}
	return &svc.GetContractABIResponse{
		Abi: c.Abi,
	}, nil
}

// GetContractBytecode loads contract bytecode
func (r *ContractRegistry) GetContractBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractBytecodeResponse, error) {
	c, ok := r.getContract(req.GetContract())
	if !ok {
		return nil, errors.NotFoundError("contract bytecode not found").SetComponent(component)
	}
	return &svc.GetContractBytecodeResponse{
		Bytecode: c.Bytecode,
	}, nil
}

// GetContractDeployedBytecode loads contract deployed bytecode
func (r *ContractRegistry) GetContractDeployedBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractDeployedBytecodeResponse, error) {
	c, ok := r.getContract(req.GetContract())
	if !ok {
		return nil, errors.NotFoundError("contract deployed bytecode not found").SetComponent(component)
	}
	return &svc.GetContractDeployedBytecodeResponse{
		DeployedBytecode: c.DeployedBytecode,
	}, nil
}

// getIndexedCount returns the count of indexed inputs in the event
func getIndexedCount(event ethabi.Event) uint {
	var indexedInputCount uint
	for i := range event.Inputs {
		if event.Inputs[i].Indexed {
			indexedInputCount++
		}
	}
	return indexedInputCount
}

// getCodehash retrieve codehash of a contract instance
func (r *ContractRegistry) getCodehash(contract common.AccountInstance) (ethcommon.Hash, error) {
	codehashes, ok := r.codehashes[contract.GetChain().String()]
	if ok {
		codehash, ok := codehashes[contract.GetAccount().Address()]
		if ok {
			return codehash, nil
		}
	}

	instance, _ := contract.Short()
	return ethcommon.Hash{},
		errors.NotFoundError(
			"contract instance %q not registered", instance,
		).SetComponent(component)
}

// GetMethodsBySelector load method using 4 bytes unique selector and the address of the contract
func (r *ContractRegistry) GetMethodsBySelector(ctx context.Context, req *svc.GetMethodsBySelectorRequest) (*svc.GetMethodsBySelectorResponse, error) {
	var sel [4]byte
	copy(sel[:], req.GetSelector())

	// Search in specific method storage
	codehash, err := r.getCodehash(*req.GetAccountInstance())
	if err == nil {
		methods, ok := r.methods[codehash][sel]
		if ok && len(methods) == 1 {
			pMethod, err := json.Marshal(methods[0])
			if err != nil {
				return nil, errors.FromError(err).SetComponent(component)
			}
			return &svc.GetMethodsBySelectorResponse{
				Method:         pMethod,
				DefaultMethods: nil,
			}, nil
		}
	}

	// Search in default methods
	defaultMethods, ok := r.methods[defaultCodehash][sel]
	if ok {
		pDefaultMethods := make([][]byte, len(defaultMethods))
		for i := range defaultMethods {
			pDefaultMethods[i], err = json.Marshal(defaultMethods[i])
			if err != nil {
				return nil, errors.FromError(err).SetComponent(component)
			}

			return &svc.GetMethodsBySelectorResponse{
				Method:         nil,
				DefaultMethods: pDefaultMethods,
			}, nil
		}
	}

	return nil, errors.NotFoundError("method not found").SetComponent(component)
}

// GetEventsBySigHash load event using event signature hash
func (r *ContractRegistry) GetEventsBySigHash(ctx context.Context, req *svc.GetEventsBySigHashRequest) (*svc.GetEventsBySigHashResponse, error) {
	// Search in specific event storage
	codehash, err := r.getCodehash(*req.GetAccountInstance())
	sigHash := ethcommon.BytesToHash(req.SigHash)
	indexedInputCount := uint(req.IndexedInputCount)

	if err == nil {
		events, ok := r.events[codehash][sigHash]
		if ok {
			matchingEvents, ok := events[indexedInputCount]
			if ok && len(matchingEvents) == 1 {
				pEvent, err := json.Marshal(matchingEvents[0])
				if err != nil {
					return nil, errors.FromError(err).SetComponent(component)
				}
				return &svc.GetEventsBySigHashResponse{
					Event:         pEvent,
					DefaultEvents: nil,
				}, nil
			}
		}
	}

	// Search in default events
	if defaultEvents, ok := r.events[defaultCodehash][sigHash][indexedInputCount]; ok {
		pDefaultEvents := make([][]byte, len(defaultEvents))
		for i := range defaultEvents {
			pDefaultEvents[i], err = json.Marshal(defaultEvents[i])
			if err != nil {
				return nil, errors.FromError(err).SetComponent(component)
			}
		}
		return &svc.GetEventsBySigHashResponse{
			Event:         nil,
			DefaultEvents: pDefaultEvents,
		}, nil
	}

	return nil, errors.NotFoundError("events not found").SetComponent(component)
}

// Request an update of the codehash of the contract address
func (r *ContractRegistry) RequestAddressUpdate(ctx context.Context, req *svc.AddressUpdateRequest) (*svc.AddressUpdateResponse, error) {
	chainID, addr := req.GetAccountInstance().GetChain(), req.GetAccountInstance().GetAccount().Address()

	// Codehash already stored for this contract instance
	if _, ok := r.codehashes[chainID.String()][addr]; ok {
		return &svc.AddressUpdateResponse{}, nil
	}

	// Codehash not stored, trying to retrieve it from chain
	code, err := r.ec.CodeAt(context.Background(), chainID.ID(), addr, nil)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	chainStr := chainID.String()
	if _, ok := r.codehashes[chainStr]; !ok {
		r.codehashes[chainStr] = make(map[ethcommon.Address]ethcommon.Hash)
	}

	r.codehashes[chainStr][addr] = crypto.Keccak256Hash(code)

	return &svc.AddressUpdateResponse{}, nil
}
