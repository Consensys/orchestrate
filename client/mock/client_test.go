package clientmock

import (
	"context"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/mock"
)

var ERC20 = []byte(`[
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"name":"account","type":"address"},
			{"indexed":false,"name":"account2","type":"address"}
		],
		"name":"MinterAdded",
		"type":"event"
	},
	{
		"constant":true,
		"inputs":[
			{"name":"account","type":"address"}
		],
		"name":"isMinter",
		"outputs":[
			{"name":"","type":"bool"}
		],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	}]`)

var methodSig = []byte("isMinter(address)")
var eventSig = []byte("MinterAdded(address,address)")

var erc20Contract = &abi.Contract{
	Id: &abi.ContractId{
		Name: "ERC20",
		Tag:  "v1.0.0",
	},
	Abi:              ERC20,
	Bytecode:         []byte{1, 2},
	DeployedBytecode: []byte{1, 2, 3},
}

var queryContract = &abi.Contract{
	Id: &abi.ContractId{
		Name: "ERC20",
		Tag:  "v1.0.0",
	},
}

func TestContractRegistryClient(t *testing.T) {
	blocks := make(map[string][]*ethtypes.Block)
	mec := mock.NewClient(blocks)

	client := New(mec)
	var c interface{} = client
	_, ok := c.(svc.RegistryClient)
	assert.True(t, ok, "Should match ContractRegistryClient interface")

	_, err := client.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: erc20Contract},
	)
	assert.NoError(t, err)

	abiResp, err := client.GetContractABI(
		context.Background(),
		&svc.GetContractRequest{Contract: queryContract},
	)
	assert.NotNil(t, abiResp, "Method should be available")
	assert.NoError(t, err, "Should not error")

	bytecodeResp, err := client.GetContractBytecode(
		context.Background(),
		&svc.GetContractRequest{Contract: queryContract},
	)
	assert.NotNil(t, bytecodeResp, "Bytecode should be available")
	assert.NoError(t, err, "Should not error")

	deployedBytecodeResp, err := client.GetContractDeployedBytecode(
		context.Background(),
		&svc.GetContractRequest{Contract: queryContract},
	)
	assert.NotNil(t, deployedBytecodeResp, "DeployedBytecode should be available")
	assert.NoError(t, err, "Should not error")

	methodResp, err := client.GetMethodsBySelector(
		context.Background(),
		&svc.GetMethodsBySelectorRequest{
			Selector:        crypto.Keccak256(methodSig)[:4],
			AccountInstance: &common.AccountInstance{},
		})
	assert.NoError(t, err)
	assert.NotNil(t, methodResp)

	eventResp, err := client.GetEventsBySigHash(context.Background(),
		&svc.GetEventsBySigHashRequest{
			SigHash:           crypto.Keccak256Hash(eventSig).Bytes(),
			AccountInstance:   &common.AccountInstance{},
			IndexedInputCount: 1})
	assert.NoError(t, err)
	assert.NotNil(t, eventResp)
}
