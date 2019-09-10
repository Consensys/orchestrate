package redis

import (
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

// Conn is a wrapper around a remote.Conn that handles internal errors
type Conn struct{ remote.Conn }

// Dial returns a connection to redis
func Dial(network, address string, options ...remote.DialOption) (remote.Conn, error) {
	conn, err := remote.Dial(network, address, options...)
	if err != nil {
		return conn, errors.ConnectionError(err.Error())
	}
	return Conn{conn}, nil
}

// Do sends a commands to the remote Redis instance
func (conn Conn) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := conn.Conn.Do(commandName, args...)
	if err != nil {
		return reply, errors.ConnectionError(err.Error())
	}
	return reply, nil
}

// RegisterContract registers a contract
func (r *ContractRegistry) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest) (*svc.RegisterContractResponse, error) {
	conn := r.pool.Get()
	defer conn.Close()

	contract := req.GetContract()

	if contract.Bytecode != nil {
		bytecodeHash := crypto.Keccak256Hash(contract.Bytecode)
		_, err := conn.Do("SET", contract.Short(), bytecodeHash)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		marshalledContract, _ := proto.Marshal(&abi.Contract{
			Abi:              contract.Abi,
			Bytecode:         contract.Bytecode,
			DeployedBytecode: contract.DeployedBytecode,
		})

		_, err = conn.Do("SET", bytecodeHash, marshalledContract)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}
	}

	if len(contract.Abi) != 0 {
		// Preformat the keys and values that we are going to
		codeHash := crypto.Keccak256Hash(contract.DeployedBytecode)
		contractAbi, err := contract.ToABI()
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		methodJSONs, eventJSONs, err := utils.ParseJSONABI(contract.Abi)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		for _, method := range contractAbi.Methods {

			sel := utils.SigHashToSelector(method.Id())

			if contract.DeployedBytecode != nil {
				// Registers the methods list
				_, err = conn.Do("RPUSH", methodKey(codeHash[:], sel), methodJSONs[method.Name])
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}
			}

			found := false

			// Attempts to find the registered method
			reply, err := conn.Do("LRANGE", methodKey(defaultCodeHash[:], sel), 0, -1)
			if err != nil {
				return nil, errors.FromError(err).ExtendComponent(component)
			}

			if reply != nil {
				replySlice, err := remote.ByteSlices(reply, nil)
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}

				for _, registeredMethod := range replySlice {
					if reflect.DeepEqual(registeredMethod, methodJSONs[method.Name]) {
						found = true
						break
					}
				}
			}

			// If not found, register it
			if !found {
				_, err := conn.Do("RPUSH", methodKey(defaultCodeHash[:], sel), methodJSONs[method.Name])
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}
			}
		}

		for _, event := range contractAbi.Events {

			eventID := event.Id()
			indexedCount := getIndexedCount(event)

			if contract.DeployedBytecode != nil {
				_, err := conn.Do("RPUSH", eventKey(codeHash[:], eventID[:], indexedCount), eventJSONs[event.Name])
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}
			}

			found := false

			// Attempts to find the registered event
			reply, err := conn.Do("LRANGE", eventKey(defaultCodeHash[:], eventID[:], indexedCount), 0, -1)
			if err != nil {
				return nil, errors.FromError(err).ExtendComponent(component)
			}

			if reply != nil {
				replySlice, err := remote.ByteSlices(reply, nil)
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}

				for _, registeredEvent := range replySlice {
					if reflect.DeepEqual(registeredEvent, eventJSONs[event.Name]) {
						found = true
						break
					}
				}
			}

			// If not found, register it
			if !found {
				_, err := conn.Do("RPUSH", eventKey(defaultCodeHash[:], eventID[:], indexedCount), eventJSONs[event.Name])
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}
			}
		}
	}

	return &svc.RegisterContractResponse{}, nil
}

// GetContract looks up an *abi.Contract object stored in Redis
func (r *ContractRegistry) GetContract(contractName string) (*abi.Contract, error) {

	conn := r.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("GET", contractName)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if reply == nil {
		return nil, errors.NotFoundError("Could not find bytecode hash for contract name").SetComponent(component)
	}

	bytecodeHash, err := remote.Bytes(reply, nil)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	reply, err = conn.Do("GET", bytecodeHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if reply == nil {
		return nil, errors.NotFoundError("Contract ABI not found").SetComponent(component)
	}

	contractBytes, err := remote.Bytes(reply, nil)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	contract := &abi.Contract{}
	err = proto.Unmarshal(contractBytes, contract)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return contract, nil
}

// GetContractABI retrieve contract ABI
func (r *ContractRegistry) GetContractABI(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractABIResponse, error) {

	contract, err := r.GetContract(req.GetContractId().Short())
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractABIResponse{
		Abi: contract.Abi,
	}, nil
}

// GetContractBytecode returns the bytecode
func (r *ContractRegistry) GetContractBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractBytecodeResponse, error) {

	contract, err := r.GetContract(req.GetContractId().Short())
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractBytecodeResponse{
		Bytecode: contract.Bytecode,
	}, nil
}

// GetContractDeployedBytecode returns the deployed bytecode
func (r *ContractRegistry) GetContractDeployedBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractDeployedBytecodeResponse, error) {

	contract, err := r.GetContract(req.GetContractId().Short())
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractDeployedBytecodeResponse{
		DeployedBytecode: contract.DeployedBytecode,
	}, nil
}

// getCodehash retrieve codehash of a contract instance
func (r *ContractRegistry) getCodehash(contract common.AccountInstance) (ethcommon.Hash, bool, error) {

	conn := r.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("GET",
		codeHashKey(contract.GetChain().String(), contract.GetAccount().GetRaw()))

	if err != nil {
		return ethcommon.Hash{}, false, errors.FromError(err).ExtendComponent(component)
	}

	if reply == nil {
		return ethcommon.Hash{}, false, nil
	}

	codeHash, err := remote.Bytes(reply, nil)
	if err != nil {
		return ethcommon.Hash{}, false, errors.FromError(err).ExtendComponent(component)
	}

	return ethcommon.BytesToHash(codeHash), true, nil
}

// GetMethodsBySelector retrieve methods using 4 bytes unique selector
func (r *ContractRegistry) GetMethodsBySelector(ctx context.Context, req *svc.GetMethodsBySelectorRequest) (*svc.GetMethodsBySelectorResponse, error) {

	conn := r.pool.Get()
	defer conn.Close()

	var selector [4]byte
	copy(selector[:], req.GetSelector())

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
	reply, err := conn.Do("LRANGE", methodKey(codeHashToLookup[:], selector), 0, -1)
	
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	
	if reply == nil {
		return nil, errors.NotFoundError("method not found").SetComponent(component)		
	}

	methodJSONs, err := remote.ByteSlices(reply, nil)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	response := &svc.GetMethodsBySelectorResponse{}

	switch {
	case codeHashFound:
		response.Method = methodJSONs[0]
	default:
		response.DefaultMethods = methodJSONs
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
