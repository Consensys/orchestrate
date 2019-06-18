package registry

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	viper.Set("abis", "ERC20[v0.1.3]:[{\"constant\":true,\"inputs\":[],\"name\":\"myFunction\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]")
	Init(context.Background())

	res, err := GlobalRegistry().GetContractABI(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "v0.1.3",
		},
	})
	assert.NotNil(t, res, "Method should be available")
	assert.NoError(t, err, "Should not error")

	res, err = GlobalRegistry().GetContractBytecode(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "v0.1.3",
		},
	})
	assert.NotNil(t, res, "Bytecode should be available")
	assert.NoError(t, err, "Should not error")

	res, err = GlobalRegistry().GetContractDeployedBytecode(&abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "v0.1.3",
		},
	})
	assert.NotNil(t, res, "DeployedBytecode should be available")
	assert.NoError(t, err, "Should not error")
}
