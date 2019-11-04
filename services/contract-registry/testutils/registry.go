package testutils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/services/contract-registry"
	rcommon "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
	"golang.org/x/net/context"
)

// ContractRegistryTestSuite  is a test suit for EnvelopeStore
type ContractRegistryTestSuite struct {
	suite.Suite
	R svc.RegistryServer
}

// ERC20 is a unittest value
var ERC20 = []byte(
	`[{
    "anonymous": false,
    "inputs": [
      {"indexed": true, "name": "account", "type": "address"},
      {"indexed": false, "name": "account2", "type": "address"}
    ],
    "name": "MinterAdded",
    "type": "event"
  },
  {
	"anonymous": false,
	"inputs": [
	  {"indexed": false, "name": "account", "type": "address"},
	  {"indexed": true, "name": "account2", "type": "address"}
	],
	"name": "MinterAddedBis",
	"type": "event"
  },
  {
    "constant": true,
    "inputs": [
      {"name": "account", "type": "address"}
    ],
    "name": "isMinter",
    "outputs": [
      {"name": "", "type": "bool"}
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
    }]`)

// ERC20bis is a unittest value
var ERC20bis = []byte(
	`[{
	"anonymous": false,
	"inputs": [
	  {"indexed": false, "name": "account", "type": "address"},
	  {"indexed": true, "name": "account2", "type": "address"}
	],
	"name": "MinterAdded",
	"type": "event"
  },
  {
	"anonymous": false,
	"inputs": [
	  {"indexed": false, "name": "account", "type": "address"},
	  {"indexed": true, "name": "account2", "type": "address"}
	],
	"name": "MinterAddedBis",
	"type": "event"
  },
  {
	"constant": true,
	"inputs": [
	  {"name": "account", "type": "address"}
	],
	"name": "isMinter",
	"outputs": [
	  {"name": "", "type": "bool"}
	],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
	}]`)

var methodSig = []byte("isMinter(address)")
var eventSig = []byte("MinterAdded(address,address)")

// ERC20Contract is a unittest value
var ERC20Contract = &abi.Contract{
	Id: &abi.ContractId{
		Name: "ERC20",
		Tag:  "v1.0.0",
	},
	Abi:              ERC20,
	Bytecode:         []byte{1, 2},
	DeployedBytecode: []byte{1, 2, 3},
}
var compactedERC20, _ = ERC20Contract.GetABICompacted()

// ERC20ContractBis is a unittest value
var ERC20ContractBis = &abi.Contract{
	Id: &abi.ContractId{
		Name: "ERC20",
	},
	Abi:              ERC20bis,
	Bytecode:         []byte{1, 3},
	DeployedBytecode: []byte{1, 2, 4},
}

var methodJSONs, eventJSONs, _ = rcommon.ParseJSONABI(ERC20Contract.Abi)
var _, eventJSONsBis, _ = rcommon.ParseJSONABI(ERC20ContractBis.Abi)

// ContractInstance is a unittest value
var ContractInstance = common.AccountInstance{
	Chain:   &chain.Chain{Id: big.NewInt(3).Bytes()},
	Account: ethereum.HexToAccount("0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0"),
}

// TestRegisterContract unit test for contract registration
func (s *ContractRegistryTestSuite) TestRegisterContract() {
	_, _ = s.R.RegisterContract(
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

	// TODO: Harmonize behavior between mock and redis/contract-registry
	// Mock allow user to provide incomplete contract data
	// While redis enforce that all data is correctly passed

	_, err := s.R.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(s.T(), err, "Should register contract properly")

	_, err = s.R.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(s.T(), err, "Should register contract properly twice")
}

// TestContractRegistryBySig checks the self-consistency of the contract-registry
func (s *ContractRegistryTestSuite) TestContractRegistryBySig() {
	_, err := s.R.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(s.T(), err)
	_, err = s.R.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20ContractBis},
	)
	assert.NoError(s.T(), err)

	// Get Contract
	contractResp, err := s.R.GetContract(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v1.0.0",
			},
		})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), compactedERC20, contractResp.GetContract().GetAbi())

	abiResp, err := s.R.GetContractABI(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "covfefe",
			},
		})
	assert.Error(s.T(), err, "GetContractABI should error when unknown contract")
	ierr, ok := err.(*ierror.Error)
	assert.True(s.T(), ok, "GetContractABI error should cast to internal error")
	assert.Equal(s.T(), "contract-registry", ierr.GetComponent()[:17], "GetContractABI error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetContractABI error should be a storage error")
	assert.Nil(s.T(), abiResp)

	// Get ABI
	abiResp, err = s.R.GetContractABI(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v1.0.0",
			},
		})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), compactedERC20, abiResp.GetAbi())

	abiResp, err = s.R.GetContractABI(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "covfefe",
			},
		})
	assert.Error(s.T(), err, "GetContractABI should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(s.T(), ok, "GetContractABI error should cast to internal error")
	assert.Equal(s.T(), "contract-registry", ierr.GetComponent()[:17], "GetContractABI error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetContractABI error should be a storage error")
	assert.Nil(s.T(), abiResp)

	// Get Bytecode
	bytecodeResp, err := s.R.GetContractBytecode(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v1.0.0",
			},
		})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ERC20Contract.Bytecode, bytecodeResp.GetBytecode())
	bytecodeResp, err = s.R.GetContractBytecode(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "covfefe",
			},
		})
	assert.Error(s.T(), err, "GetContractBytecode should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(s.T(), ok, "GetContractBytecode error should cast to internal error")
	assert.Equal(s.T(), "contract-registry", ierr.GetComponent()[:17], "GetContractBytecode error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetContractBytecode error should be a storage error")
	assert.Nil(s.T(), bytecodeResp)

	// Get DeployedBytecode
	deployedBytecodeResp, err := s.R.GetContractDeployedBytecode(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v1.0.0",
			},
		})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ERC20Contract.DeployedBytecode, deployedBytecodeResp.GetDeployedBytecode())
	deployedBytecodeResp, err = s.R.GetContractDeployedBytecode(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "covfefe",
			},
		})
	assert.Error(s.T(), err, "Should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(s.T(), ok, "GetContractDeployedBytecode should cast to internal error")
	assert.Equal(s.T(), "contract-registry", ierr.GetComponent()[:17], "GetContractDeployedBytecode error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetContractDeployedBytecode error should be a storage error")
	assert.Nil(s.T(), deployedBytecodeResp)

	// Get Catalog
	namesResp, err := s.R.GetCatalog(context.Background(), &svc.GetCatalogRequest{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{"ERC20"}, namesResp.GetNames())

	// Get Tags
	tagsResp, err := s.R.GetTags(context.Background(), &svc.GetTagsRequest{Name: "ERC20"})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{"latest", "v1.0.0"}, tagsResp.GetTags())

	// Get MethodBySelector on default
	methodResp, err := s.R.GetMethodsBySelector(context.Background(),
		&svc.GetMethodsBySelectorRequest{
			Selector:        crypto.Keccak256(methodSig)[:4],
			AccountInstance: &common.AccountInstance{},
		})
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), methodResp.GetMethod())
	assert.Equal(s.T(), [][]byte{methodJSONs["isMinter"]}, methodResp.GetDefaultMethods())

	// Get EventsBySigHash wrong indexed count
	eventResp, err := s.R.GetEventsBySigHash(context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 0})
	assert.Error(s.T(), err)
	ierr, ok = err.(*ierror.Error)
	assert.True(s.T(), ok, "GetEventsBySigHash error should cast to internal error")
	assert.Equal(s.T(), "contract-registry", ierr.GetComponent()[:17], "GetEventsBySigHash error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetEventsBySigHash error should be a storage error")
	assert.Nil(s.T(), eventResp.GetEvent())
	assert.Nil(s.T(), eventResp.GetDefaultEvents())

	// Get EventsBySigHash
	eventResp, err = s.R.GetEventsBySigHash(context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 1})
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), eventResp.GetEvent())
	assert.Equal(s.T(), [][]byte{eventJSONs["MinterAdded"], eventJSONsBis["MinterAdded"]}, eventResp.GetDefaultEvents())

	// Update smart-contract address
	_, err = s.R.SetAccountCodeHash(context.Background(),
		&svc.SetAccountCodeHashRequest{
			AccountInstance: &ContractInstance,
			CodeHash:        crypto.Keccak256([]byte{1, 2, 3}),
		})
	assert.NoError(s.T(), err)

	// Get MethodBySelector
	methodResp, err = s.R.GetMethodsBySelector(context.Background(),
		&svc.GetMethodsBySelectorRequest{
			Selector:        crypto.Keccak256(methodSig)[:4],
			AccountInstance: &ContractInstance})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), methodJSONs["isMinter"], methodResp.GetMethod())
	assert.Nil(s.T(), methodResp.GetDefaultMethods())

	// Get EventsBySigHash
	eventResp, err = s.R.GetEventsBySigHash(
		context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 1})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), eventJSONs["MinterAdded"], eventResp.GetEvent())
	assert.Nil(s.T(), eventResp.GetDefaultEvents())
}
