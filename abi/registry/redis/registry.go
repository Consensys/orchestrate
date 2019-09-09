package redis

import (
	"reflect"
	"context"
	"github.com/golang/protobuf/proto"

	"github.com/ethereum/go-ethereum/crypto"
	remote	"github.com/gomodule/redigo/redis"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/utils"
)

const component = "redis-registry"

var defaultCodehash = ethcommon.Hash{}

// ContractRegistry is a Redis based implementation of the interface pkg.git/services/contract-registry.RegistryServer
type ContractRegistry struct {
	pool *remote.Pool
}

// Conn is a wrapper around a remote.Conn that handles internal errors
type Conn struct { remote.Conn }

// Dial returns a connexion to redis
func Dial(network, address string, options ...remote.DialOption) (remote.Conn, error) {
	conn, err := remote.Dial(network, address, options...)
	if err != nil {
		return conn, errors.ConnectionError(err.Error())
	}
	return Conn{ conn }, nil
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
		conn.Do("SET", contract.Short(), bytecodeHash)

		marshalledContract, _ := proto.Marshal(&abi.Contract {
			Abi:              contract.Abi,
			Bytecode:         contract.Bytecode,
			DeployedBytecode: contract.DeployedBytecode,
		})

		conn.Do("SET", bytecodeHash, marshalledContract)
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

			sel:= utils.SigHashToSelector(method.Id())

			if contract.DeployedBytecode != nil {
				// Registers the methods list
				_, err = conn.Do("RPUSH", methodKey(codeHash[:], sel), [][]byte{methodJSONs[method.Name]})
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}	
			} 

			found := false
			
			// Attemps to find the registered method
			reply, err := conn.Do("LRANGE", methodKey(defaultCodehash[:], sel), 0 , -1)
			if err != nil {
				return nil, errors.FromError(err).ExtendComponent(component)
			}

			if reply != nil {
				replySlice, err := remote.ByteSlices(reply, nil)
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}

				searchMethodLoop:
				for _, registeredMethod := range replySlice {
					if reflect.DeepEqual(registeredMethod, methodJSONs[method.Name]) {
						found = true
						break searchMethodLoop
					}
				}
			}

			// If not found, register it
			if !found {
				_, err := conn.Do("RPUSH", methodKey(defaultCodehash[:], sel), [][]byte{methodJSONs[method.Name]})
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}
			}
		}

		for _, event := range contractAbi.Events {

			eventID := event.Id()
			indexedCount := getIndexedCount(event)

			if contract.DeployedBytecode != nil {
				_, err := conn.Do("RPUSH", eventKey(codeHash[:], eventID[:], indexedCount), [][]byte{eventJSONs[event.Name]})
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}
			}

			found := false
			
			// Attemps to find the registered event
			reply, err := conn.Do("LRANGE", eventKey(defaultCodehash[:], eventID[:], indexedCount), 0 , -1)
			if err != nil {
				return nil, errors.FromError(err).ExtendComponent(component)
			}

			if reply != nil {
				replySlice, err := remote.ByteSlices(reply, nil)
				if err != nil {
					return nil, errors.FromError(err).ExtendComponent(component)
				}

				searchEventLoop:
				for _, registeredEvent := range replySlice {
					if reflect.DeepEqual(registeredEvent, eventJSONs[event.Name]) {
						found = true
						break searchEventLoop
					}
				}
			}

			// If not found, register it
			if !found {
				_, err := conn.Do("RPUSH", eventKey(defaultCodehash[:], eventID[:], indexedCount), [][]byte{eventJSONs[event.Name]})
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

	contract, err := r.GetContract(req.GetContract().Short())
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractABIResponse{
		Abi: contract.Abi,
	}, nil
}

// GetContractBytecode returns the bytecode
func (r *ContractRegistry) GetContractBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractBytecodeResponse, error) {

	contract, err := r.GetContract(req.GetContract().Short())
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractBytecodeResponse {
		Bytecode: contract.Bytecode,
	}, nil
}

// GetContractDeployedBytecode returns the deployed bytecode
func (r *ContractRegistry) GetContractDeployedBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractDeployedBytecodeResponse, error) {

	contract, err := r.GetContract(req.GetContract().Short())
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.GetContractDeployedBytecodeResponse {
		DeployedBytecode: contract.DeployedBytecode,
	}, nil
}

// getCodehash retrieve codehash of a contract instance
func (r *ContractRegistry) getCodehash(contract common.AccountInstance) (ethcommon.Hash, error) {
	
	contractName, err := contract.Short()
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}

	conn := r.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("GET", 
		codeHashKey(contract.GetChain().String(), contract.GetAccount().GetRaw()))

	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}

	if reply == nil {
		return ethcommon.Hash{}, errors.NotFoundError(
			"No CodeHash found for %q", contractName,
		).SetComponent(component)
	}

	codeHash, err := remote.Bytes(reply, nil)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}

	return ethcommon.BytesToHash(codeHash), nil
}

// GetMethodsBySelector retrieve methods using 4 bytes unique selector
func (r *ContractRegistry) GetMethodsBySelector(ctx context.Context, req *svc.GetMethodsBySelectorRequest) (*svc.GetMethodsBySelectorResponse, error) {

	// conn := r.pool.Get()
	// defer conn.Close()

	// selector := utils.SigHashToSelector(req.GetSelector())

	// codehash, err := r.getCodehash(*req.GetAccountInstance())
	// if err == nil {
	// 	reply, err := conn.Do("LRANGE", eventKey(codeHash[:], selector[:], indexedCount), 0 , -1)
	// }

	return nil, errors.NotFoundError("method not found").SetComponent(component)

}

// GetEventsBySigHash retrieve events using hash of signature
func (r *ContractRegistry) GetEventsBySigHash(context.Context, *svc.GetEventsBySigHashRequest) (*svc.GetEventsBySigHashResponse, error) {
	return nil, errors.NotFoundError("method not found").SetComponent(component)
}

// RequestAddressUpdate request an update of the codehash of the contract address
func (r *ContractRegistry) RequestAddressUpdate(context.Context, *svc.AddressUpdateRequest) (*svc.AddressUpdateResponse, error)

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
