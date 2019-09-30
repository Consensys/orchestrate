package registry

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/abi"
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
