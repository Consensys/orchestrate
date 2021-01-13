package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
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
