package testutils

import (
	"github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/types/api"
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
		CodeHash: FakeHash(),
	}
}

func FakeSearchContractRequest() *api.SearchContractRequest {
	return &api.SearchContractRequest{
		CodeHash: FakeHash(),
		Address:  FakeAddress(),
	}
}

func FakeContractResponse() *api.ContractResponse {
	c := FakeContract()
	return &api.ContractResponse{
		Name:             c.Name,
		Tag:              c.Tag,
		ABI:              c.ABI,
		Bytecode:         c.Bytecode,
		DeployedBytecode: c.DeployedBytecode,
	}
}
