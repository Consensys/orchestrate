package contractregistry

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

func TestInit(t *testing.T) {
	viper.Set("abis", "ERC20[v0.1.3]:[{\"constant\":true,\"inputs\":[],\"name\":\"myFunction\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]")
	Init(context.Background())

	abiResp, err := GlobalRegistry().GetContractABI(
		context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v0.1.3",
			},
		})
	assert.NotNil(t, abiResp, "Method should be available")
	assert.NoError(t, err, "Should not error")

	bytecodeResp, err := GlobalRegistry().GetContractBytecode(
		context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v0.1.3",
			},
		})
	assert.NotNil(t, bytecodeResp, "Bytecode should be available")
	assert.NoError(t, err, "Should not error")

	deployedBytecodeResp, err := GlobalRegistry().GetContractDeployedBytecode(
		context.Background(),
		&svc.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: "ERC20",
				Tag:  "v0.1.3",
			},
		})
	assert.NotNil(t, deployedBytecodeResp, "DeployedBytecode should be available")
	assert.NoError(t, err, "Should not error")
}
