package memory

import (
	"context"
	"reflect"
	"sort"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/services/contract-registry"
	rcommon "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
)

// ContractRegistry stores contract ABI and bytecode in memory
type ContractRegistry struct {
	// Contract registry/name#tag to bytecode hash
	contractHashes map[string]map[string]ethcommon.Hash

	// Bytecode hash to artifacts (ABI, bytecode and deployed bytecode)
	artifacts map[ethcommon.Hash]*rcommon.Artifact

	// Address to Codehash (deployed bytecode hash) map
	codehashes map[string]map[ethcommon.Address]ethcommon.Hash

	// Codehash to Selector to method ABIs
	methods map[ethcommon.Hash]map[[4]byte][][]byte

	// Codehash to SigHash to event ABIs
	events map[ethcommon.Hash]map[ethcommon.Hash]map[uint][][]byte
}

var defaultCodehash = ethcommon.Hash{}

// NewRegistry creates a ContractRegistry
func NewRegistry() *ContractRegistry {
	r := &ContractRegistry{
		contractHashes: make(map[string]map[string]ethcommon.Hash),
		artifacts:      make(map[ethcommon.Hash]*rcommon.Artifact),
		codehashes:     make(map[string]map[ethcommon.Address]ethcommon.Hash),
		methods:        make(map[ethcommon.Hash]map[[4]byte][][]byte),
		events:         make(map[ethcommon.Hash]map[ethcommon.Hash]map[uint][][]byte),
	}
	r.methods[defaultCodehash] = make(map[[4]byte][][]byte)
	r.events[defaultCodehash] = make(map[ethcommon.Hash]map[uint][][]byte)
	return r
}

// RegisterContract register a contract including ABI, bytecode and deployed bytecode
func (r *ContractRegistry) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest) (*svc.RegisterContractResponse, error) {
	contract := req.GetContract()

	bytecode, deployedBytecode, abiRaw, err := rcommon.CheckExtractArtifacts(contract)
	if err != nil {
		return nil, err
	}

	name, tag, err := rcommon.CheckExtractNameTag(contract.GetId())
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if bytecode != nil {
		bytecodeHash := crypto.Keccak256Hash(bytecode)

		if r.contractHashes[name] == nil {
			r.contractHashes[name] = make(map[string]ethcommon.Hash)
		}

		r.contractHashes[name][tag] = bytecodeHash
		r.contractHashes[name]["latest"] = bytecodeHash

		r.artifacts[bytecodeHash] = &rcommon.Artifact{
			Abi:              abiRaw,
			Bytecode:         bytecode,
			DeployedBytecode: deployedBytecode,
		}
	}

	if len(abiRaw) != 0 {
		codeHash := crypto.Keccak256Hash(deployedBytecode)
		contractAbi, err := contract.ToABI()
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}
		methodJSONs, eventJSONs, err := rcommon.ParseJSONABI(abiRaw)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		for _, m := range contractAbi.Methods {
			// Register methods for this bytecode
			method := m
			sel := rcommon.SigHashToSelector(method.ID())
			if deployedBytecode != nil {
				// Init map
				if r.methods[codeHash] == nil {
					r.methods[codeHash] = make(map[[4]byte][][]byte)
				}

				r.methods[codeHash][sel] = [][]byte{methodJSONs[method.Name]}
			}

			// Register in default methods if not present
			found := false
			for _, registeredMethod := range r.methods[defaultCodehash][sel] {
				if reflect.DeepEqual(registeredMethod, methodJSONs[method.Name]) {
					found = true
				}
			}
			if !found {
				r.methods[defaultCodehash][sel] = append(
					r.methods[defaultCodehash][sel],
					methodJSONs[method.Name],
				)
			}
		}

		for _, e := range contractAbi.Events {
			event := e
			indexedCount := rcommon.GetIndexedCount(event)

			// Register events for this bytecode
			if deployedBytecode != nil {
				// Init map
				if r.events[codeHash] == nil {
					r.events[codeHash] = make(map[ethcommon.Hash]map[uint][][]byte)
				}
				// Init map
				if r.events[codeHash][event.ID()] == nil {
					r.events[codeHash][event.ID()] = make(map[uint][][]byte)
				}

				r.events[codeHash][event.ID()][indexedCount] = [][]byte{eventJSONs[event.Name]}
			}

			// Init map
			if r.events[defaultCodehash][event.ID()] == nil {
				r.events[defaultCodehash][event.ID()] = make(map[uint][][]byte)
			}
			// Register in default events if not present
			found := false
			for _, registeredEvent := range r.events[defaultCodehash][event.ID()][indexedCount] {
				if reflect.DeepEqual(registeredEvent, eventJSONs[event.Name]) {
					found = true
				}
			}
			if !found {
				r.events[defaultCodehash][event.ID()][indexedCount] = append(
					r.events[defaultCodehash][event.ID()][indexedCount],
					eventJSONs[event.Name],
				)
			}
		}
	}

	return &svc.RegisterContractResponse{}, nil
}

// DeregisterContract remove the name + tag association to a contract artifact (abi, bytecode, deployedBytecode). Artifacts are not deleted.
func (r *ContractRegistry) DeregisterContract(ctx context.Context, req *svc.DeregisterContractRequest) (*svc.DeregisterContractResponse, error) {
	delete(r.contractHashes, req.GetContractId().Short())
	return &svc.DeregisterContractResponse{}, nil
}

// DeleteArtifact remove an artifacts based on its BytecodeHash.
func (r *ContractRegistry) DeleteArtifact(ctx context.Context, req *svc.DeleteArtifactRequest) (*svc.DeleteArtifactResponse, error) {
	delete(r.artifacts, ethcommon.BytesToHash(req.GetBytecodeHash()))
	return &svc.DeleteArtifactResponse{}, nil
}

func (r *ContractRegistry) getArtifact(id *abi.ContractId) (a *rcommon.Artifact, ok bool) {
	name, tag, err := rcommon.CheckExtractNameTag(id)
	if err != nil {
		return nil, false
	}

	a, ok = r.artifacts[r.contractHashes[name][tag]]
	return a, ok
}

func (r *ContractRegistry) getContract(id *abi.ContractId) (c *abi.Contract, ok bool) {
	name, tag, err := rcommon.CheckExtractNameTag(id)
	if err != nil {
		return nil, false
	}

	a, ok := r.artifacts[r.contractHashes[name][tag]]
	if !ok {
		return nil, ok
	}

	return &abi.Contract{
		Id:               id,
		Abi:              a.Abi,
		Bytecode:         a.Bytecode,
		DeployedBytecode: a.DeployedBytecode,
	}, ok
}

// GetContract loads a contract
func (r *ContractRegistry) GetContract(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractResponse, error) {
	id := req.GetContractId()
	_, _, err := rcommon.CheckExtractNameTag(id)
	if err != nil {
		return nil, err
	}

	c, ok := r.getContract(id)
	if !ok {
		return nil, errors.StorageError("could not load contract (%v)", err).ExtendComponent(component)
	}

	return &svc.GetContractResponse{
		Contract: c,
	}, nil
}

// GetContractABI loads contract ABI
func (r *ContractRegistry) GetContractABI(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractABIResponse, error) {
	id := req.GetContractId()
	_, _, err := rcommon.CheckExtractNameTag(id)
	if err != nil {
		return nil, err
	}

	a, ok := r.getArtifact(id)
	if !ok {
		return nil, errors.NotFoundError("contract ABI not found").SetComponent(component)
	}

	return &svc.GetContractABIResponse{
		Abi: a.Abi,
	}, nil
}

// GetContractBytecode loads contract bytecode
func (r *ContractRegistry) GetContractBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractBytecodeResponse, error) {
	id := req.GetContractId()
	_, _, err := rcommon.CheckExtractNameTag(id)
	if err != nil {
		return nil, err
	}

	a, ok := r.getArtifact(id)
	if !ok {
		return nil, errors.NotFoundError("contract bytecode not found").SetComponent(component)
	}
	return &svc.GetContractBytecodeResponse{
		Bytecode: a.Bytecode,
	}, nil
}

// GetContractDeployedBytecode loads contract deployed bytecode
func (r *ContractRegistry) GetContractDeployedBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractDeployedBytecodeResponse, error) {
	id := req.GetContractId()
	_, _, err := rcommon.CheckExtractNameTag(id)
	if err != nil {
		return nil, err
	}

	a, ok := r.getArtifact(id)
	if !ok {
		return nil, errors.NotFoundError("contract deployed bytecode not found").SetComponent(component)
	}
	return &svc.GetContractDeployedBytecodeResponse{
		DeployedBytecode: a.DeployedBytecode,
	}, nil
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
	sel := rcommon.SigHashToSelector(req.GetSelector())

	// Search in specific method storage
	codehash, err := r.getCodehash(*req.GetAccountInstance())
	if err == nil {
		methods, ok := r.methods[codehash][sel]
		if ok && len(methods) == 1 {
			return &svc.GetMethodsBySelectorResponse{
				Method:         methods[0],
				DefaultMethods: nil,
			}, nil
		}
	}

	// Search in default methods
	defaultMethods, ok := r.methods[defaultCodehash][sel]
	if ok {
		return &svc.GetMethodsBySelectorResponse{
			Method:         nil,
			DefaultMethods: defaultMethods,
		}, nil
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
				return &svc.GetEventsBySigHashResponse{
					Event:         matchingEvents[0],
					DefaultEvents: nil,
				}, nil
			}
		}
	}

	// Search in default events
	if defaultEvents, ok := r.events[defaultCodehash][sigHash][indexedInputCount]; ok {
		return &svc.GetEventsBySigHashResponse{
			Event:         nil,
			DefaultEvents: defaultEvents,
		}, nil
	}

	return nil, errors.NotFoundError("events not found").SetComponent(component)
}

// GetCatalog returns a list of all registered contracts.
func (r *ContractRegistry) GetCatalog(ctx context.Context, req *svc.GetCatalogRequest) (*svc.GetCatalogResponse, error) {
	resp := &svc.GetCatalogResponse{}
	for name := range r.contractHashes {
		resp.Names = append(resp.Names, name)
	}
	sort.Strings(resp.Names)
	return resp, nil
}

// GetTags returns a list of all tags available for a contract name.
func (r *ContractRegistry) GetTags(ctx context.Context, req *svc.GetTagsRequest) (*svc.GetTagsResponse, error) {
	if _, ok := r.contractHashes[req.GetName()]; !ok {
		return nil, errors.NotFoundError("No Tags found for requested contract name").ExtendComponent(component)
	}

	resp := &svc.GetTagsResponse{}
	for tag := range r.contractHashes[req.GetName()] {
		resp.Tags = append(resp.Tags, tag)
	}
	sort.Strings(resp.Tags)
	return resp, nil
}

// SetAccountCodeHash set the codehash of a contract address for a given chain
func (r *ContractRegistry) SetAccountCodeHash(ctx context.Context, req *svc.SetAccountCodeHashRequest) (*svc.SetAccountCodeHashResponse, error) {
	chainID, addr := req.GetAccountInstance().GetChain(), req.GetAccountInstance().GetAccount().Address()

	// Codehash already stored for this contract instance
	if _, ok := r.codehashes[chainID.String()][addr]; ok {
		return &svc.SetAccountCodeHashResponse{}, nil
	}

	chainStr := chainID.String()
	if _, ok := r.codehashes[chainStr]; !ok {
		r.codehashes[chainStr] = make(map[ethcommon.Address]ethcommon.Hash)
	}

	r.codehashes[chainStr][addr] = ethcommon.BytesToHash(req.GetCodeHash())

	return &svc.SetAccountCodeHashResponse{}, nil
}
