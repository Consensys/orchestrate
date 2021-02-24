package testutils

import (
	"github.com/ConsenSys/orchestrate/pkg/encoding/json"
	"github.com/ConsenSys/orchestrate/pkg/types/api"
)

func FakeRegisterContractRequest() *api.RegisterContractRequest {
	c := FakeContract()
	var abi interface{}
	_ = json.Unmarshal([]byte(c.ABI), &abi)

	return &api.RegisterContractRequest{
		Name:             c.Name,
		Tag:              c.Tag,
		ABI:              abi,
		Bytecode:         c.Bytecode,
		DeployedBytecode: c.DeployedBytecode,
	}
}

func FakeSetContractCodeHashRequest() *api.SetContractCodeHashRequest {
	return &api.SetContractCodeHashRequest{
		CodeHash: FakeHash().String(),
	}
}
