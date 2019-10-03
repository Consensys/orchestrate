package redis

import (
	"context"
	"sort"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	remote "github.com/gomodule/redigo/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/abi"
)

var (
	defaultCodeHash = ethcommon.Hash{}
)

// ContractRegistry is a Redis based implementation of the interface pkg.git/services/contract-registry.RegistryServer
type ContractRegistry struct {
	pool *remote.Pool
}

// Conn dials remote redis and returns a new connection
func (r *ContractRegistry) Conn() *Conn {
	return &Conn{Conn: r.pool.Get()}
}

// NewRegistry creates a ContractRegistry
func NewRegistry(pool *remote.Pool) *ContractRegistry {
	return &ContractRegistry{
		pool: pool,
	}
}

// RegisterContract registers a contract
func (r *ContractRegistry) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest) (*svc.RegisterContractResponse, error) {

	conn := r.Conn()
	defer conn.Close()

	contract := req.GetContract()

	bytecode, deployedBytecode, abiRaw, err := common.CheckExtractArtifacts(contract)
	if err != nil {
		return nil, err
	}

	name, tag, err := common.CheckExtractNameTag(contract.Id)
	if err != nil {
		return nil, err
	}

	// Attempt deserializing the ABI
	contractAbi, err := contract.ToABI()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Attempts deserializing methods and events
	methodJSONs, eventJSONs, err := common.ParseJSONABI(abiRaw)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	byteCodeHash := crypto.Keccak256Hash(bytecode)
	deployedByteCodeHash := crypto.Keccak256Hash(deployedBytecode)

	err = Catalog.PushIfNotExist(conn, name)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	err = Tags.PushIfNotExist(conn, name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	err = ByteCodeHash.Set(conn, name, tag, byteCodeHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if err = Artifact.Set(conn, byteCodeHash,
		&abi.Contract{
			Abi:              abiRaw,
			Bytecode:         bytecode,
			DeployedBytecode: deployedBytecode,
		},
	); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if err := Methods.Registers(conn,
		deployedByteCodeHash,
		defaultCodeHash,
		contractAbi.Methods,
		methodJSONs,
	); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if err := Events.Registers(conn,
		deployedByteCodeHash,
		defaultCodeHash,
		contractAbi.Events,
		eventJSONs,
	); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.RegisterContractResponse{}, nil
}

// DeregisterContract remove the name + tag association to a contract artifact (abi, bytecode, deployedBytecode). Artifacts are not deleted.
func (r *ContractRegistry) DeregisterContract(ctx context.Context, req *svc.DeregisterContractRequest) (*svc.DeregisterContractResponse, error) {
	return nil, errors.FeatureNotSupportedError("Registry does not support Deregistration yet")
}

// DeleteArtifact remove the name + tag association to a contract artifact (abi, bytecode, deployedBytecode). Artifacts are not deleted.
func (r *ContractRegistry) DeleteArtifact(ctx context.Context, req *svc.DeleteArtifactRequest) (*svc.DeleteArtifactResponse, error) {
	return nil, errors.FeatureNotSupportedError("Registry does not support Deregistration yet")
}

// getContract looks up an *abi.Contract object stored in Redis
func (r *ContractRegistry) getContract(name, tag string) (*abi.Contract, error) {

	conn := r.Conn()
	defer conn.Close()

	byteCodeHash, ok, err := ByteCodeHash.Get(conn, name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	} else if !ok {
		return nil, errors.NotFoundError("No contract bytecode hash found for given name and tags %v:%v", name, tag).SetComponent(component)
	}

	artifact, ok, err := Artifact.Get(conn, byteCodeHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	} else if !ok {
		return nil, errors.NotFoundError("No artifact found for found bytecode hash %v. This means that the state of the database is inconsistent", byteCodeHash).SetComponent(component)
	}

	return artifact, nil
}

// GetContract retrieves the whole contract object
func (r *ContractRegistry) GetContract(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractResponse, error) {

	contractID := req.GetContractId()
	if contractID == nil {
		return nil, errors.InvalidArgError("No contract ID found in request").ExtendComponent(component)
	}

	name, tag, err := common.CheckExtractNameTag(contractID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	contract, err := r.getContract(name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractResponse{
		Contract: contract,
	}, nil
}

// GetContractABI retrieve contract ABI
func (r *ContractRegistry) GetContractABI(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractABIResponse, error) {

	contractID := req.GetContractId()
	if contractID == nil {
		return nil, errors.InvalidArgError("No contract ID found in request").ExtendComponent(component)
	}

	name, tag, err := common.CheckExtractNameTag(contractID)
	if err != nil {
		return nil, err
	}

	contract, err := r.getContract(name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractABIResponse{
		Abi: contract.Abi,
	}, nil
}

// GetContractBytecode returns the bytecode
func (r *ContractRegistry) GetContractBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractBytecodeResponse, error) {

	contractID := req.GetContractId()
	if contractID == nil {
		return nil, errors.InvalidArgError("No contract ID found in request").ExtendComponent(component)
	}

	name, tag, err := common.CheckExtractNameTag(contractID)
	if err != nil {
		return nil, err
	}

	contract, err := r.getContract(name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractBytecodeResponse{
		Bytecode: contract.Bytecode,
	}, nil
}

// GetContractDeployedBytecode returns the deployed bytecode
func (r *ContractRegistry) GetContractDeployedBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractDeployedBytecodeResponse, error) {

	contractID := req.GetContractId()
	if contractID == nil {
		return nil, errors.InvalidArgError("No contract ID found in request").ExtendComponent(component)
	}

	name, tag, err := common.CheckExtractNameTag(contractID)
	if err != nil {
		return nil, err
	}

	contract, err := r.getContract(name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractDeployedBytecodeResponse{
		DeployedBytecode: contract.DeployedBytecode,
	}, nil
}

// GetMethodsBySelector retrieve methods using 4 bytes unique selector
func (r *ContractRegistry) GetMethodsBySelector(ctx context.Context, req *svc.GetMethodsBySelectorRequest) (*svc.GetMethodsBySelectorResponse, error) {

	conn := r.Conn()
	defer conn.Close()

	selectorBytes := req.GetSelector()
	if len(selectorBytes) == 0 {
		return nil, errors.InvalidArgError("No selector found in request").ExtendComponent(component)
	}

	selector := common.SigHashToSelector(req.GetSelector())

	// Flag used to detect if we need to query with the default code hash or not
	codeFound := false
	deployedByteCodeHash := defaultCodeHash

	// common.Check if the address and chainID have been provided in the request
	accountChain, address, err := common.CheckExtractChainAddress(req.GetAccountInstance())
	addressChainProvided := err == nil

	if addressChainProvided {
		// If a contract has been registered at given chainID:address, override the default value we gave earlier
		deployedByteCodeHash, codeFound, err = DeployedByteCodeHash.Get(conn, accountChain, address)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}
	}

	method, methodFound, err := Methods.Get(conn, deployedByteCodeHash, selector)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	} else if !methodFound {
		return nil, errors.NotFoundError("Method not found for given selector").SetComponent(component)
	}

	response := &svc.GetMethodsBySelectorResponse{}

	switch {
	case codeFound:
		response.Method = method[0]
	default:
		response.DefaultMethods = method
	}

	return response, nil
}

// GetEventsBySigHash retrieve events using hash of signature
func (r *ContractRegistry) GetEventsBySigHash(ctx context.Context, req *svc.GetEventsBySigHashRequest) (*svc.GetEventsBySigHashResponse, error) {
	conn := r.Conn()
	defer conn.Close()

	sigHashBytes := req.GetSigHash()
	// Nil value of uint32 is 0, hence the uint cast is error-free
	indexedCount := uint(req.GetIndexedInputCount())

	if len(sigHashBytes) == 0 {
		return nil, errors.InvalidArgError("No selector found in request").ExtendComponent(component)
	}

	sigHash := ethcommon.BytesToHash(sigHashBytes)

	// Flag used to detect if we need to query with the default code hash or not
	codeFound := false
	deployedByteCodeHash := defaultCodeHash

	// common.Check if the address and chainID have been provided in the request
	accountChain, address, err := common.CheckExtractChainAddress(req.GetAccountInstance())
	addressChainProvided := err == nil

	if addressChainProvided {
		deployedByteCodeHash, codeFound, err = DeployedByteCodeHash.Get(conn, accountChain, address)
		if err != nil {
			return nil, err
		}
	}

	// Case where the connection to redis is codeHashFound, but the hash is not found
	if !codeFound {
		// Use the defaultCodeHash instead, and try to look at the default contract-registry
		deployedByteCodeHash = defaultCodeHash
	}

	event, eventFound, err := Events.Get(conn, deployedByteCodeHash, sigHash, indexedCount)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	} else if !eventFound {
		return nil, errors.NotFoundError("event not found").SetComponent(component)
	}

	response := &svc.GetEventsBySigHashResponse{}

	switch {
	case codeFound:
		response.Event = event[0]
	default:
		response.DefaultEvents = event
	}

	return response, nil
}

// GetCatalog returns a list of all registered contracts.
func (r *ContractRegistry) GetCatalog(ctx context.Context, req *svc.GetCatalogRequest) (*svc.GetCatalogResponse, error) {
	conn := r.Conn()
	defer conn.Close()

	// By convention, the catalog always exists
	catalog, _, err := Catalog.Get(conn)
	if err != nil {
		return nil, err
	}

	sort.Strings(catalog)
	return &svc.GetCatalogResponse{
		Names: catalog,
	}, nil
}

// GetTags returns a list of all tags available for a contract name.
func (r *ContractRegistry) GetTags(ctx context.Context, req *svc.GetTagsRequest) (*svc.GetTagsResponse, error) {
	conn := r.Conn()
	defer conn.Close()

	name := req.GetName()
	if name == "" {
		return nil, errors.InvalidArgError("Name provided was empty").ExtendComponent(component)
	}

	tags, ok, err := Tags.Get(conn, name)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if !ok {
		return nil, errors.NotFoundError("No Tags found for requested contract name").ExtendComponent(component)
	}

	sort.Strings(tags)
	return &svc.GetTagsResponse{
		Tags: tags,
	}, nil
}

// SetAccountCodeHash set the codehash of a contract address for a given chain
func (r *ContractRegistry) SetAccountCodeHash(ctx context.Context, req *svc.SetAccountCodeHashRequest) (*svc.SetAccountCodeHashResponse, error) {
	conn := r.Conn()
	defer conn.Close()

	hashBytes := req.GetCodeHash()
	if len(hashBytes) == 0 {
		return nil, errors.InvalidArgError("No deployed contract bytecode hash found in request").ExtendComponent(component)
	}

	deployedByteCodeHash := ethcommon.BytesToHash(hashBytes)

	accountChain, address, err := common.CheckExtractChainAddress(req.GetAccountInstance())
	if err != nil {
		return nil, err
	}

	err = DeployedByteCodeHash.Set(conn, accountChain, address, deployedByteCodeHash)
	if err != nil {
		return nil, err
	}

	return &svc.SetAccountCodeHashResponse{}, nil
}
