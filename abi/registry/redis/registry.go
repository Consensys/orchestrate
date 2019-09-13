package redis

import (
	"context"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	remote "github.com/gomodule/redigo/redis"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/utils"
)

var (
	defaultTag      = "latest"
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
func NewRegistry() *ContractRegistry {
	r := &ContractRegistry{}
	return r
}

// RegisterContract registers a contract
func (r *ContractRegistry) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest) (*svc.RegisterContractResponse, error) {

	conn := r.Conn()
	defer conn.Close()

	contract := req.GetContract()

	bytecode, deployedBytecode, abiRaw, err := checkExtractArtifacts(contract)
	if err != nil {
		return nil, err
	}

	name, tag, err := checkExtractNameTag(contract)
	if err != nil {
		return nil, err
	}

	// Attempt deserializing the ABI
	contractAbi, err := contract.ToABI()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Attempts deserializing methods and events
	methodJSONs, eventJSONs, err := utils.ParseJSONABI(abiRaw)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	byteCodeHash := crypto.Keccak256Hash(bytecode)
	deployedByteCodeHash := crypto.Keccak256Hash(deployedBytecode)

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

// DeregisterContract remove the name + tag association to a contract artifact (abi, bytecode, deployedBytecode). Artifacts are not deleted.
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
		return nil, errors.NotFoundError("No contract found for given name and tags").SetComponent(component)
	}

	artifact, ok, err := Artifact.Get(conn, byteCodeHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	} else if !ok {
		return nil, errors.NotFoundError("No contract found for given name and tags").SetComponent(component)
	}

	return artifact, nil
}

// GetContractABI retrieve contract ABI
func (r *ContractRegistry) GetContract(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractResponse, error) {

	contractID := req.GetContractId()
	if contractID == nil {
		return nil, errors.InvalidArgError("No contract ID found in request").ExtendComponent(component)
	}

	contract, err := r.getContract(contractID.Name, contractID.Tag)
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

	contract, err := r.getContract(contractID.Name, contractID.Tag)
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

	contract, err := r.getContract(contractID.Name, contractID.Tag)
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

	contract, err := r.getContract(contractID.Name, contractID.Tag)
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

	selector := utils.SigHashToSelector(req.GetSelector())

	accountChain, address, err := checkExtractChainAddress(req.GetAccountInstance())
	if err != nil {
		return nil, err
	}

	deployedByteCodeHash, codeFound, err := DeployedByteCodeHash.Get(conn, accountChain, address)

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if !codeFound {
		deployedByteCodeHash = defaultCodeHash
	}

	method, methodFound, err := Methods.Get(conn, deployedByteCodeHash, selector)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if !methodFound {
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

	accountChain, address, err := checkExtractChainAddress(req.GetAccountInstance())
	if err != nil {
		return nil, err
	}

	deployedByteCodeHash, codeFound, err := DeployedByteCodeHash.Get(conn, accountChain, address)
	if err != nil {
		return nil, err
	}

	// Case where the connection to redis is codeHashFound, but the hash is not found
	if !codeFound {
		// Use the defaultCodeHash instead, and try to look at the default registry
		deployedByteCodeHash = defaultCodeHash
	}

	event, ok, err := Events.Get(conn, deployedByteCodeHash, sigHash, indexedCount)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if !ok {
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

// GetCatalog returns a list of all registered contracts. Name is used to filter contractIds based on their contract name, empty to list all contract names & tags.
func (r *ContractRegistry) GetCatalog(ctx context.Context, req *svc.GetCatalogRequest) (*svc.GetCatalogResponse, error) {
	conn := r.Conn()
	defer conn.Close()

	// By convention, the catalog always exists
	catalog, _, err := Catalog.Get(conn)
	if err != nil {
		return nil, err
	}

	return &svc.GetCatalogResponse{
		Names: catalog,
	}, nil
}

// GetTags returns a list of all registered contracts. Name is used to filter contractIds based on their contract name, empty to list all contract names & tags.
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

	return &svc.GetTagsResponse{
		Tags: tags,
	}, nil
}

// SetAccountCodeHash request an update of the codehash of the contract address
func (r *ContractRegistry) SetAccountCodeHash(ctx context.Context, req *svc.SetAccountCodeHashRequest) (*svc.SetAccountCodeHashResponse, error) {
	conn := r.Conn()
	defer conn.Close()

	hashBytes := req.GetCodeHash()
	if len(hashBytes) == 0 {
		return nil, errors.InvalidArgError("No deployed contract bytecode hash found in request").ExtendComponent(component)
	}

	deployedByteCodeHash := ethcommon.BytesToHash(hashBytes)

	accountChain, address, err := checkExtractChainAddress(req.GetAccountInstance())
	if err != nil {
		return nil, err
	}

	err = DeployedByteCodeHash.Set(conn, accountChain, address, deployedByteCodeHash)
	if err != nil {
		return nil, err
	}

	return &svc.SetAccountCodeHashResponse{}, nil
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

func checkExtractChainAddress(accountInstance *common.AccountInstance) (*chain.Chain, *ethereum.Account, error) {
	if accountInstance == nil {
		return nil, nil, errors.InvalidArgError("No account instance found in request").ExtendComponent(component)
	}

	accountChain := accountInstance.GetChain()
	if accountChain == nil {
		return nil, nil, errors.InvalidArgError("No ethereum chainID found in request").ExtendComponent(component)
	}

	address := accountInstance.GetAccount()
	if address == nil {
		return nil, nil, errors.InvalidArgError("No ethereum account instance found in request").ExtendComponent(component)
	}

	return accountChain, address, nil
}

func checkExtractArtifacts(contract *abi.Contract) (bytecode, deployedBytecode, abiBytes []byte, err error) {
	if contract.Bytecode == nil {
		return []byte{}, []byte{}, []byte{}, errors.InvalidArgError("No contract bytecode provided in request").ExtendComponent(component)
	}

	if contract.DeployedBytecode != nil {
		return []byte{}, []byte{}, []byte{}, errors.InvalidArgError("No contract deployed bytecode provided in request").ExtendComponent(component)
	}

	if len(contract.Abi) != 0 {
		return []byte{}, []byte{}, []byte{}, errors.InvalidArgError("No abi provided in request").ExtendComponent(component)
	}

	return contract.Bytecode, contract.DeployedBytecode, contract.Abi, nil
}

func checkExtractNameTag(contract *abi.Contract) (name, tag string, err error) {
	name = contract.GetName()
	if name == "" {
		return "", "", errors.InvalidArgError("No abi provided in request").ExtendComponent(component)
	}

	// Set Tag to latest if it was not set in the request
	tag = contract.GetTag()
	if tag == "" {
		tag = defaultTag
	}

	return name, tag, nil
}
