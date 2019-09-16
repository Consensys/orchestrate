package suites

import (
	"testing"
	"math/big"

	"golang.org/x/net/context"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"github.com/ethereum/go-ethereum/crypto"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/utils"
)

// RegistrySuite gathers tests that are common to the registry no matter 
type RegistrySuite struct {
	suite.Suite
	Server svc.RegistryServer
}

// NewRegistrySuite returns a testing suite
func NewRegistrySuite(server svc.RegistryServer) *RegistrySuite {
	suite := new(RegistrySuite)
	suite.Server = server
	return suite
}

// Run the full registry test suite
func RunRegistry(t *testing.T, server svc.RegistryServer) {
	suite.Run(t, NewRegistrySuite(server))
}

// ERC20 is used for testing
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

// ERC20bis is used for testing
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

// ERC20Contract is used for testing
var ERC20Contract = &abi.Contract{
	Id: &abi.ContractId{
		Name: "ERC20",
		Tag:  "v1.0.0",
	},
	Abi:              ERC20,
	Bytecode:         []byte{1, 2},
	DeployedBytecode: []byte{1, 2, 3},
}

// ERC20ContractBis is used for testing
var ERC20ContractBis = &abi.Contract{
	Id: &abi.ContractId{
		Name: "ERC20",
		Tag:  "v1.0.1",
	},
	Abi:              ERC20bis,
	Bytecode:         []byte{1, 3},
	DeployedBytecode: []byte{1, 2, 4},
}

var methodJSONs, eventJSONs, _ = utils.ParseJSONABI(ERC20Contract.Abi)
var _, eventJSONsBis, _ = utils.ParseJSONABI(ERC20ContractBis.Abi)

// ContractInstance used for testing
var ContractInstance = common.AccountInstance{
	Chain:   &chain.Chain{Id: big.NewInt(3).Bytes()},
	Account: ethereum.HexToAccount("0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0"),
}

// TestRegisterContract unit test conract registration
func (s *RegistrySuite) TestRegisterContract() {
	_, err := s.Server.RegisterContract(
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
	assert.NoError(s.T(), err, "Should not error on empty things")

	_, err = s.Server.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(s.T(), err, "Should register contract properly")

	_, err = s.Server.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(s.T(), err, "Should register contract properly twice")
}

// TestContractRegistryBySig unit tests contract registration
func (s *RegistrySuite) TestContractRegistryBySig() {
	_, err := s.Server.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20Contract},
	)
	assert.NoError(s.T(), err)
	_, err = s.Server.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: ERC20ContractBis},
	)
	assert.NoError(s.T(), err)

	// Get Contract
	contractResp, err := s.Server.GetContract(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v1.0.0",
			},
		})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ERC20Contract.Abi, contractResp.GetContract().GetAbi())

	abiResp, err := s.Server.GetContractABI(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "covfefe",
			},
		})
	assert.Error(s.T(), err, "GetContractABI should error when unknown contract")
	ierr, ok := err.(*ierror.Error)
	assert.True(s.T(), ok, "GetContractABI error should cast to internal error")
	assert.Equal(s.T(), "contract-registry.mock", ierr.GetComponent(), "GetContractABI error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetContractABI error should be a storage error")
	assert.Nil(s.T(), abiResp)

	// Get ABI
	abiResp, err = s.Server.GetContractABI(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v1.0.0",
			},
		})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ERC20Contract.Abi, abiResp.GetAbi())

	abiResp, err = s.Server.GetContractABI(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "covfefe",
			},
		})
	assert.Error(s.T(), err, "GetContractABI should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(s.T(), ok, "GetContractABI error should cast to internal error")
	assert.Equal(s.T(), "contract-registry.mock", ierr.GetComponent(), "GetContractABI error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetContractABI error should be a storage error")
	assert.Nil(s.T(), abiResp)

	// Get Bytecode
	bytecodeResp, err := s.Server.GetContractBytecode(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v1.0.0",
			},
		})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ERC20Contract.Bytecode, bytecodeResp.GetBytecode())
	bytecodeResp, err = s.Server.GetContractBytecode(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "covfefe",
			},
		})
	assert.Error(s.T(), err, "GetContractBytecode should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(s.T(), ok, "GetContractBytecode error should cast to internal error")
	assert.Equal(s.T(), "contract-registry.mock", ierr.GetComponent(), "GetContractBytecode error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetContractBytecode error should be a storage error")
	assert.Nil(s.T(), bytecodeResp)

	// Get DeployedBytecode
	deployedBytecodeResp, err := s.Server.GetContractDeployedBytecode(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v1.0.0",
			},
		})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ERC20Contract.DeployedBytecode, deployedBytecodeResp.GetDeployedBytecode())
	deployedBytecodeResp, err = s.Server.GetContractDeployedBytecode(context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "covfefe",
			},
		})
	assert.Error(s.T(), err, "Should error when unknown contract")
	ierr, ok = err.(*ierror.Error)
	assert.True(s.T(), ok, "GetContractDeployedBytecode should cast to internal error")
	assert.Equal(s.T(), "contract-registry.mock", ierr.GetComponent(), "GetContractDeployedBytecode error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetContractDeployedBytecode error should be a storage error")
	assert.Nil(s.T(), deployedBytecodeResp)

	// Get Catalog
	namesResp, err := s.Server.GetCatalog(context.Background(), &svc.GetCatalogRequest{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{"ERC20"}, namesResp.GetNames())

	// Get Tags
	tagsResp, err := s.Server.GetTags(context.Background(), &svc.GetTagsRequest{Name: "ERC20"})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{"v1.0.0", "v1.0.1"}, tagsResp.GetTags())

	// Get MethodBySelector on default
	methodResp, err := s.Server.GetMethodsBySelector(context.Background(),
		&svc.GetMethodsBySelectorRequest{
			Selector:        crypto.Keccak256(methodSig)[:4],
			AccountInstance: &common.AccountInstance{},
		})
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), methodResp.GetMethod())
	assert.Equal(s.T(), [][]byte{methodJSONs["isMinter"]}, methodResp.GetDefaultMethods())

	// Get EventsBySigHash wrong indexed count
	eventResp, err := s.Server.GetEventsBySigHash(context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 0})
	assert.Error(s.T(), err)
	ierr, ok = err.(*ierror.Error)
	assert.True(s.T(), ok, "GetEventsBySigHash error should cast to internal error")
	assert.Equal(s.T(), "contract-registry.mock", ierr.GetComponent(), "GetEventsBySigHash error component should be correct")
	assert.True(s.T(), errors.IsStorageError(ierr), "GetEventsBySigHash error should be a storage error")
	assert.Nil(s.T(), eventResp.GetEvent())
	assert.Nil(s.T(), eventResp.GetDefaultEvents())

	// Get EventsBySigHash
	eventResp, err = s.Server.GetEventsBySigHash(context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 1})
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), eventResp.GetEvent())
	assert.Equal(s.T(), [][]byte{eventJSONs["MinterAdded"], eventJSONsBis["MinterAdded"]}, eventResp.GetDefaultEvents())

	// Update smart-contract address
	_, err = s.Server.SetAccountCodeHash(context.Background(),
		&svc.SetAccountCodeHashRequest{
			AccountInstance: &ContractInstance,
			CodeHash:        crypto.Keccak256([]byte{1, 2, 3}),
		})
	assert.NoError(s.T(), err)

	// Get MethodBySelector
	methodResp, err = s.Server.GetMethodsBySelector(context.Background(),
		&svc.GetMethodsBySelectorRequest{
			Selector:        crypto.Keccak256(methodSig)[:4],
			AccountInstance: &ContractInstance})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), methodJSONs["isMinter"], methodResp.GetMethod())
	assert.Nil(s.T(), methodResp.GetDefaultMethods())

	// Get EventsBySigHash
	eventResp, err = s.Server.GetEventsBySigHash(
		context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &ContractInstance,
			IndexedInputCount: 1})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), eventJSONs["MinterAdded"], eventResp.GetEvent())
	assert.Nil(s.T(), eventResp.GetDefaultEvents())
}



