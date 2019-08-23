package mock

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/mock"
)

var ERC20 = []byte(
	`[{
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "name": "account",
        "type": "address"
      },
      {
        "indexed": false,
        "name": "account2",
        "type": "address"
      }
    ],
    "name": "MinterAdded",
    "type": "event"
  },
  {
    "constant": true,
    "inputs": [
      {
        "name": "account",
        "type": "address"
      }
    ],
    "name": "isMinter",
    "outputs": [
      {
        "name": "",
        "type": "bool"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
    }]`)

var ERC20bis = []byte(
	`[{
        "anonymous": false,
        "inputs": [
          {
            "indexed": false,
            "name": "account",
            "type": "address"
          },
          {
            "indexed": true,
            "name": "account2",
            "type": "address"
          }
        ],
        "name": "MinterAdded",
        "type": "event"
      },
	  {
        "anonymous": false,
        "inputs": [
          {
            "indexed": false,
            "name": "account",
            "type": "address"
          },
          {
            "indexed": true,
            "name": "account2",
            "type": "address"
          }
        ],
        "name": "MinterAddedBis",
        "type": "event"
      },
      {
        "constant": true,
        "inputs": [
          {
            "name": "account",
            "type": "address"
          }
        ],
        "name": "isMinter",
        "outputs": [
          {
            "name": "",
            "type": "bool"
          }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
        }]`)

var methodSig = []byte("isMinter(address)")
var eventSig = []byte("MinterAdded(address,address)")

var ERC20Contract = &abi.Contract{
	Id: &abi.ContractId{
		Name: "ERC20",
		Tag:  "v1.0.0",
	},
	Abi:              ERC20,
	Bytecode:         []byte{1, 2},
	DeployedBytecode: []byte{1, 2, 3},
}
var ERC20ContractBis = &abi.Contract{
	Id: &abi.ContractId{
		Name: "ERC20",
		Tag:  "v1.0.1",
	},
	Abi:              ERC20bis,
	Bytecode:         []byte{1, 3},
	DeployedBytecode: []byte{1, 2, 4},
}

var ERC20ABI, _ = ERC20Contract.ToABI()
var ERC20ABIBis, _ = ERC20ContractBis.ToABI()

var ContractInstance = common.AccountInstance{
	Chain:   &chain.Chain{Id: big.NewInt(3).Bytes()},
	Account: ethereum.HexToAccount("0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0"),
}

func TestRegisterContract(t *testing.T) {
	blocks := make(map[string][]*ethtypes.Block)
	mec := mock.NewClient(blocks)

	r := NewRegistry(mec)
	_, err := r.RegisterContract(
		context.Background(),
		&svc.RegisterContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "ERC20",
					Tag:  "v1.0.0",
				},
				Abi: []byte{},
			},
		},
	)
	assert.NoError(t, err, "Should not error on empty things")

	_, err = r.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(t, err, "Should register contract properly")

	_, err = r.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(t, err, "Should register contract properly twice")
}

func TestContractRegistryBySig(t *testing.T) {
	blocks := make(map[string][]*ethtypes.Block)
	mec := mock.NewClient(blocks)

	r := NewRegistry(mec)
	_, err := r.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(t, err)
	_, err = r.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20ContractBis},
	)
	assert.NoError(t, err)

	// Get ABI
	abiResp, err := r.GetContractABI(context.Background(),
		&svc.GetContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "ERC20",
					Tag:  "v1.0.0",
				},
			},
		})
	assert.NoError(t, err)
	assert.Equal(t, ERC20Contract.Abi, abiResp.GetAbi())

	abiResp, err = r.GetContractABI(context.Background(),
		&svc.GetContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "ERC20",
					Tag:  "covfefe",
				},
			},
		})
	assert.Error(t, err, "GetContractABI should error when unknown contract")
	ierr, ok := err.(*ierror.Error)
	assert.True(t, ok, "GetContractABI error should cast to internal error")
	assert.Equal(t, "contract-registry.mock", ierr.GetComponent(), "GetContractABI error component should be correct")
	assert.True(t, errors.IsStorageError(ierr), "GetContractABI error should be a storage error")
	assert.Nil(t, abiResp)

	// Get Bytecode
	bytecodeResp, err := r.GetContractBytecode(context.Background(),
		&svc.GetContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "ERC20",
					Tag:  "v1.0.0",
				},
			},
		})
	assert.NoError(t, err)
	assert.Equal(t, ERC20Contract.Bytecode, bytecodeResp.GetBytecode())
	bytecodeResp, err = r.GetContractBytecode(context.Background(),
		&svc.GetContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "ERC20",
					Tag:  "covfefe",
				},
			},
		})
	assert.Error(t, err, "GetContractBytecode should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(t, ok, "GetContractBytecode error should cast to internal error")
	assert.Equal(t, "contract-registry.mock", ierr.GetComponent(), "GetContractBytecode error component should be correct")
	assert.True(t, errors.IsStorageError(ierr), "GetContractBytecode error should be a storage error")
	assert.Nil(t, bytecodeResp)

	// Get DeployedBytecode
	deployedBytecodeResp, err := r.GetContractDeployedBytecode(context.Background(),
		&svc.GetContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "ERC20",
					Tag:  "v1.0.0",
				},
			},
		})
	assert.NoError(t, err)
	assert.Equal(t, ERC20Contract.DeployedBytecode, deployedBytecodeResp.GetDeployedBytecode())
	deployedBytecodeResp, err = r.GetContractDeployedBytecode(context.Background(),
		&svc.GetContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "ERC20",
					Tag:  "covfefe",
				},
			},
		})
	assert.Error(t, err, "Should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(t, ok, "GetContractDeployedBytecode should cast to internal error")
	assert.Equal(t, "contract-registry.mock", ierr.GetComponent(), "GetContractDeployedBytecode error component should be correct")
	assert.True(t, errors.IsStorageError(ierr), "GetContractDeployedBytecode error should be a storage error")
	assert.Nil(t, deployedBytecodeResp)

	// Get MethodBySelector on default
	methodResp, err := r.GetMethodsBySelector(context.Background(),
		&svc.GetMethodsBySelectorRequest{
			Selector:        crypto.Keccak256(methodSig)[:4],
			AccountInstance: &common.AccountInstance{},
		})
	assert.NoError(t, err)
	assert.Nil(t, methodResp.GetMethod())
	expectedMethod := ERC20ABI.Methods["isMinter"]
	assert.Equal(t, []*ethabi.Method{&expectedMethod}, methodResp.GetDefaultMethods())

	// Get EventsBySigHash wrong indexed count
	eventResp, err := r.GetEventsBySigHash(context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 0})
	assert.Error(t, err)
	ierr, ok = err.(*ierror.Error)
	assert.True(t, ok, "GetEventsBySigHash error should cast to internal error")
	assert.Equal(t, "contract-registry.mock", ierr.GetComponent(), "GetEventsBySigHash error component should be correct")
	assert.True(t, errors.IsStorageError(ierr), "GetEventsBySigHash error should be a storage error")
	assert.Nil(t, eventResp.GetEvent())
	assert.Nil(t, eventResp.GetDefaultEvents())

	// Get EventsBySigHash
	eventResp, err = r.GetEventsBySigHash(context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 1})
	assert.NoError(t, err)
	expectedEvent := ERC20ABI.Events["MinterAdded"]
	expectedEventBis := ERC20ABIBis.Events["MinterAdded"]
	assert.Nil(t, eventResp.GetEvent())
	assert.Equal(t, []*ethabi.Event{&expectedEvent, &expectedEventBis}, eventResp.GetDefaultEvents())

	// Update smart-contract address
	_, err = r.RequestAddressUpdate(context.Background(),
		&svc.AddressUpdateRequest{AccountInstance: &ContractInstance})
	assert.NoError(t, err)

	// Get MethodBySelector
	methodResp, err = r.GetMethodsBySelector(context.Background(),
		&svc.GetMethodsBySelectorRequest{
			Selector:        crypto.Keccak256(methodSig)[:4],
			AccountInstance: &ContractInstance})
	assert.NoError(t, err)

	var method ethabi.Method
	err = json.Unmarshal(methodResp.GetMethod(), &method)
	assert.NoError(t, err)
	assert.Equal(t, &expectedMethod, &method)
	assert.Nil(t, methodResp.GetDefaultMethods())

	// Get EventsBySigHash
	eventResp, err = r.GetEventsBySigHash(
		context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 1})
	assert.NoError(t, err)
	assert.Equal(t, &expectedEvent, eventResp.GetEvent())
	assert.Nil(t, eventResp.GetDefaultEvents())
}
