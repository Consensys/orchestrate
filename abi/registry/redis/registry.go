package redis

import (
	"fmt"
	"context"
	"reflect"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	remote "github.com/gomodule/redigo/redis"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/utils"
)

const component = "redis-registry"

var defaultCodeHash = ethcommon.Hash{}

// ContractRegistry is a Redis based implementation of the interface pkg.git/services/contract-registry.RegistryServer
type ContractRegistry struct {
	pool *remote.Pool
}

// Conn dials remote redis and returns a new Connexion
func (r *ContractRegistry) Conn() *Conn {
	return &Conn{Conn: r.pool.Get()}
}

// RegisterContract registers a contract
func (r *ContractRegistry) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest) (*svc.RegisterContractResponse, error) {
	conn := r.Conn()
	defer conn.Close()

	contract := req.GetContract()

	if contract.Bytecode == nil {
		return nil, errors.InvalidArgError("No contract bytecode found in request").ExtendComponent(component)
	}

	if contract.DeployedBytecode != nil {
		return nil, errors.InvalidArgError("No contract deployed bytecode found in request").ExtendComponent(component)
	}

	if len(contract.Abi) != 0 {
		return nil, errors.InvalidArgError("Sent an empty ABI").ExtendComponent(component)
	}

	// Attempt deserializing the ABI
	contractAbi, err := contract.ToABI()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Attemps deserializing methods and events
	methodJSONs, eventJSONs, err := utils.ParseJSONABI(contract.Abi)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	byteCodeHash := crypto.Keccak256Hash(contract.Bytecode)
	deployedByteCodeHash := crypto.Keccak256Hash(contract.DeployedBytecode)

	err = ByteCodeHash.Set(conn, contract.Short(), contract.GetTag(), byteCodeHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if err = Artifact.Set(conn, byteCodeHash,
		&abi.Contract{
			Abi:              contract.Abi,
			Bytecode:         contract.Bytecode,
			DeployedBytecode: contract.DeployedBytecode,
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

// GetContract looks up an *abi.Contract object stored in Redis
func (r *ContractRegistry) GetContract(name string, tag string) (*abi.Contract, error) {

	conn := r.Conn()
	defer conn.Close()

	byteCodeHash, err := ByteCodeHash.Get(conn, name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	artifact, err := Artifact.Get(conn, byteCodeHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return artifact, nil
}

// GetContractABI retrieve contract ABI
func (r *ContractRegistry) GetContractABI(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractABIResponse, error) {

	contractID := req.GetContractId()
	if contractID == nil {
		return nil, errors.InvalidArgError("No contract ID found in request").ExtendComponent(component)
	}

	contract, err := r.GetContract(contractID.Name, contractID.Tag)
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

	contract, err := r.GetContract(contractID.Name, contractID.Tag)
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

	contract, err := r.GetContract(contractID.Name, contractID.Tag)
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
	
	accountInstance := req.GetAccountInstance()
	if len(selectorBytes) == 0 {
		return nil, errors.InvalidArgError("No account instance found in request").ExtendComponent(component)
	}

	chain := accountInstance.GetChain()
	if len(selectorBytes) == 0 {
		return nil, errors.InvalidArgError("No ethereum chainID found in request").ExtendComponent(component)
	}

	address := accountInstance.GetAccount()
	if len(selectorBytes) == 0 {
		return nil, errors.InvalidArgError("No ethereum account instance found in request").ExtendComponent(component)
	}

	deployedByteCodeHash, err := DeployedByteCodeHash.Get(conn, chain, address)

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	codeHashFound := true

	if deployedByteCodeHash == ethcommon.Hash{} {
		codeHashFound = false
		deployedByteCodeHash = defaultCodeHash
	}

	method, err := Methods.LRange(deployedByteCodeHash, selector)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if method == nil {
		return nil, errors.NotFoundError("Method not found for given selector").SetComponent(component)
	}

	response := &svc.GetMethodsBySelectorResponse{}

	switch {
	case codeHashFound:
		response.Method = method[0]
	default:
		response.DefaultMethods = method
	}

	return response, nil
}

// GetEventsBySigHash retrieve events using hash of signature
func (r *ContractRegistry) GetEventsBySigHash(ctx context.Context, req *svc.GetEventsBySigHashRequest) (*svc.GetEventsBySigHashResponse, error) {
	conn := r.pool.Get()
	defer conn.Close()

	selector := utils.SigHashToSelector(req.GetSigHash())
	indexedCount := req.GetIndexedInputCount()

	codeHashToLookup, codeHashFound, err := r.getCodehash(*req.GetAccountInstance())
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Case where the connexion to redis is codeHashFound, but the hash is not found
	if !codeHashFound {
		// Use the defaultCodeHash instead, and try to look at the default registry
		codeHashToLookup = defaultCodeHash
	}

	// TODO: Handle the indexedCount thing
	reply, err := conn.Do("LINDEX", eventKey(codeHashToLookup[:], selector[:], uint(indexedCount)), 0)

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if reply == nil {
		return nil, errors.NotFoundError("event not found").SetComponent(component)
	}

	eventJSONs, err := remote.ByteSlices(reply, nil)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	response := &svc.GetEventsBySigHashResponse{}

	switch {
	case codeHashFound:
		response.Event = eventJSONs[0]
	default:
		response.DefaultEvents = eventJSONs
	}

	return response, nil
}

// GetCatalog returns a list of all registered contracts. Name is used to filter contractIds based on their contract name, empty to list all contract names & tags.
func (r *ContractRegistry) GetCatalog(ctx context.Context, req *svc.GetCatalogRequest) (*svc.GetCatalogResponse, error) {
	// TODO
	return nil, nil
}

// GetTags returns a list of all registered contracts. Name is used to filter contractIds based on their contract name, empty to list all contract names & tags.
func (r *ContractRegistry) GetTags(ctx context.Context, req *svc.GetTagsRequest) (*svc.GetTagsResponse, error) {
	// TODO
	return nil, nil
}

// SetAccountCodeHash request an update of the codehash of the contract address
func (r *ContractRegistry) SetAccountCodeHash(ctx context.Context, req *svc.SetAccountCodeHashRequest) (*svc.SetAccountCodeHashResponse, error) {
	conn := r.pool.Get()
	defer conn.Close()

	chainID := req.GetAccountInstance().GetChain()
	addr := req.GetAccountInstance().GetAccount().Address()

	conn.Do("SET", codeHashKey(chainID.String(), addr[:]), req.GetCodeHash())

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
