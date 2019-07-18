package static

import (
	"math/big"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
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
	err := r.RegisterContract(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "v1.0.0",
		},
		Abi: []byte{},
	})
	assert.NoError(t, err, "Should not error on empty things")

	err = r.RegisterContract(ERC20Contract)
	assert.NoError(t, err, "Should register contract properly")

	err = r.RegisterContract(ERC20Contract)
	assert.NoError(t, err, "Should register contract properly twice")
}

func TestContractRegistryBySig(t *testing.T) {
	blocks := make(map[string][]*ethtypes.Block)
	mec := mock.NewClient(blocks)

	r := NewRegistry(mec)
	err := r.RegisterContract(ERC20Contract)
	assert.NoError(t, err)
	err = r.RegisterContract(ERC20ContractBis)
	assert.NoError(t, err)

	// Get ABI
	result, err := r.GetContractABI(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "v1.0.0",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, ERC20Contract.Abi, result)
	result, err = r.GetContractABI(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "covfefe",
		},
	})
	assert.Error(t, err, "GetContractABI should error when unknown contract")
	ierr, ok := err.(*ierror.Error)
	assert.True(t, ok, "GetContractABI error should cast to internal error")
	assert.Equal(t, "abi.registry.static", ierr.GetComponent(), "GetContractABI error component should be correct")
	assert.True(t, errors.IsStorageError(ierr), "GetContractABI error should be a storage error")
	assert.Nil(t, result)

	// Get Bytecode
	result, err = r.GetContractBytecode(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "v1.0.0",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, ERC20Contract.Bytecode, result)
	result, err = r.GetContractBytecode(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "covfefe",
		},
	})
	assert.Error(t, err, "GetContractBytecode should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(t, ok, "GetContractBytecode error should cast to internal error")
	assert.Equal(t, "abi.registry.static", ierr.GetComponent(), "GetContractBytecode error component should be correct")
	assert.True(t, errors.IsStorageError(ierr), "GetContractBytecode error should be a storage error")
	assert.Nil(t, result)

	// Get DeployedBytecode
	result, err = r.GetContractDeployedBytecode(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "v1.0.0",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, ERC20Contract.DeployedBytecode, result)
	result, err = r.GetContractDeployedBytecode(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "covfefe",
		},
	})
	assert.Error(t, err, "Should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(t, ok, "GetContractDeployedBytecode should cast to internal error")
	assert.Equal(t, "abi.registry.static", ierr.GetComponent(), "GetContractDeployedBytecode error component should be correct")
	assert.True(t, errors.IsStorageError(ierr), "GetContractDeployedBytecode error should be a storage error")
	assert.Nil(t, result)

	// Get MethodBySelector on default
	var sel [4]byte
	copy(sel[:], crypto.Keccak256(methodSig)[:4])
	method, defaultMethod, err := r.GetMethodsBySelector(sel, common.AccountInstance{})
	assert.NoError(t, err)
	assert.Nil(t, method)
	expectedMethod := ERC20ABI.Methods["isMinter"]
	assert.Equal(t, []*ethabi.Method{&expectedMethod}, defaultMethod)

	// Get EventsBySigHash wrong indexed count
	event, defaultEvent, err := r.GetEventsBySigHash(crypto.Keccak256Hash(eventSig), ContractInstance, 0)
	assert.Error(t, err)
	ierr, ok = err.(*ierror.Error)
	assert.True(t, ok, "GetEventsBySigHash error should cast to internal error")
	assert.Equal(t, "abi.registry.static", ierr.GetComponent(), "GetEventsBySigHash error component should be correct")
	assert.True(t, errors.IsStorageError(ierr), "GetEventsBySigHash error should be a storage error")
	assert.Nil(t, result)
	assert.Nil(t, event)
	assert.Nil(t, defaultEvent)

	// Get EventsBySigHash
	event, defaultEvent, err = r.GetEventsBySigHash(crypto.Keccak256Hash(eventSig), ContractInstance, 1)
	assert.NoError(t, err)
	expectedEvent := ERC20ABI.Events["MinterAdded"]
	expectedEventBis := ERC20ABIBis.Events["MinterAdded"]
	assert.Nil(t, event)
	assert.Equal(t, []*ethabi.Event{&expectedEvent, &expectedEventBis}, defaultEvent)

	// Update smart-contract address
	err = r.RequestAddressUpdate(ContractInstance)
	assert.NoError(t, err)

	// Get MethodBySelector
	method, defaultMethod, err = r.GetMethodsBySelector(sel, ContractInstance)
	assert.NoError(t, err)
	assert.Equal(t, &expectedMethod, method)
	assert.Nil(t, defaultMethod)

	// Get EventsBySigHash
	event, defaultEvent, err = r.GetEventsBySigHash(crypto.Keccak256Hash(eventSig), ContractInstance, 1)
	assert.NoError(t, err)
	assert.Equal(t, &expectedEvent, event)
	assert.Nil(t, defaultEvent)
}
