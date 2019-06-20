package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	types "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/static"
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

var methodSig = []byte("isMinter(address)")
var eventSig = []byte("MinterAdded(address,address)")

var ERC20Contract = &abi.Contract{Id: &abi.ContractId{Name: "ERC20", Tag: "v1.0.0"}, Abi: ERC20, Bytecode: []byte{1, 2}, DeployedBytecode: []byte{1, 2, 3}}

var contractInstance = common.AccountInstance{
	Chain:   &chain.Chain{Id: []byte("3")},
	Account: &ethereum.Account{Raw: []byte("0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0")},
}

type RegistryServiceTestSuite struct {
	suite.Suite
	registry *RegistryService
}

func (r *RegistryServiceTestSuite) SetupTest() {
	blocks := make(map[string][]*ethtypes.Block)
	mec := mock.NewClient(blocks)

	r.registry = NewRegistryService(static.NewRegistry(mec))
}

func (r *RegistryServiceTestSuite) TestRegistry() {
	_, err := r.registry.RegisterContract(context.Background(), &types.RegisterContractRequest{Contract: ERC20Contract})
	assert.Nil(r.T(), err, "should not error")

	abiResponse, err := r.registry.GetContractABI(context.Background(), &types.GetContractRequest{Contract: &abi.Contract{Id: &abi.ContractId{Name: "ERC20", Tag: "v1.0.0"}}})
	assert.Nil(r.T(), err, "should not error")
	assert.NotNil(r.T(), abiResponse, "abi should not be nil")

	bytecodeResponse, err := r.registry.GetContractBytecode(context.Background(), &types.GetContractRequest{Contract: &abi.Contract{Id: &abi.ContractId{Name: "ERC20", Tag: "v1.0.0"}}})
	assert.Nil(r.T(), err, "should not error")
	assert.NotNil(r.T(), bytecodeResponse, "bytecode should not be nil")

	deployedBytecodeResponse, err := r.registry.GetContractDeployedBytecode(context.Background(), &types.GetContractRequest{Contract: &abi.Contract{Id: &abi.ContractId{Name: "ERC20", Tag: "v1.0.0"}}})
	assert.Nil(r.T(), err, "should not error")
	assert.NotNil(r.T(), deployedBytecodeResponse, "deloyedBytecode should not be nil")

	methodsBySelectorResponse, err := r.registry.GetMethodsBySelector(context.Background(), &types.GetMethodsBySelectorRequest{
		Selector:        crypto.Keccak256(methodSig)[:4],
		AccountInstance: &contractInstance,
	})
	assert.Nil(r.T(), err, "should not error")
	assert.NotNil(r.T(), methodsBySelectorResponse, "methodBySelectorResponse should not be nil")

	eventsBySigHashResponse, err := r.registry.GetEventsBySigHash(context.Background(), &types.GetEventsBySigHashRequest{
		SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
		AccountInstance:   &contractInstance,
		IndexedInputCount: 1,
	})
	assert.Nil(r.T(), err, "should not error")
	assert.NotNil(r.T(), eventsBySigHashResponse, "eventsBySigHashResponse should not be nil")

	addressUpdateResponse, err := r.registry.RequestAddressUpdate(context.Background(), &types.AddressUpdateRequest{
		AccountInstance: &contractInstance,
	})
	assert.Nil(r.T(), err, "should not error")
	assert.NotNil(r.T(), addressUpdateResponse, "abi should not be nil")
}

func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(RegistryServiceTestSuite))
}
