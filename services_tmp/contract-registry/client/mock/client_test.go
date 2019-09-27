package clientmock

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/common"
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

var queryContractID = &abi.ContractId{
	Name: "ERC20",
	Tag:  "v1.0.0",
}

func TestContractRegistryClient(t *testing.T) {
	var client svc.RegistryClient = New() // Break if client does not implement interface

	_, err := client.RegisterContract(context.Background(),
		&svc.RegisterContractRequest{Contract: erc20Contract},
	)
	assert.NoError(t, err)

	contractResp, err := client.GetContract(
		context.Background(),
		&svc.GetContractRequest{ContractId: queryContractID},
	)
	assert.NotNil(t, contractResp, "Method should be available")
	assert.NoError(t, err, "Should not error")

	abiResp, err := client.GetContractABI(
		context.Background(),
		&svc.GetContractRequest{ContractId: queryContractID},
	)
	assert.NotNil(t, abiResp, "Method should be available")
	assert.NoError(t, err, "Should not error")

	bytecodeResp, err := client.GetContractBytecode(
		context.Background(),
		&svc.GetContractRequest{ContractId: queryContractID},
	)
	assert.NotNil(t, bytecodeResp, "Bytecode should be available")
	assert.NoError(t, err, "Should not error")

	deployedBytecodeResp, err := client.GetContractDeployedBytecode(
		context.Background(),
		&svc.GetContractRequest{ContractId: queryContractID},
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

	catalogResp, err := client.GetCatalog(context.Background(), &svc.GetCatalogRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, catalogResp)

	tagsResp, err := client.GetTags(context.Background(), &svc.GetTagsRequest{Name: "ERC20"})
	assert.NoError(t, err)
	assert.NotNil(t, tagsResp)
}
