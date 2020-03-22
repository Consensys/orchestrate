package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/contract-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
)

const getContractComponent = component + ".get-contract"

//go:generate mockgen -source=get_contract.go -destination=mocks/mock_get_contract.go -package=mocks

type GetContractUseCase interface {
	Execute(ctx context.Context, id *abi.ContractId) (*abi.Contract, error)
}

// GetContract is a use case to get a contract
type GetContract struct {
	artifactDataAgent store.ArtifactDataAgent
}

// NewGetContract creates a new GetContract
func NewGetContract(artifactDataAgent store.ArtifactDataAgent) *GetContract {
	return &GetContract{
		artifactDataAgent: artifactDataAgent,
	}
}

// Execute gets a contract from DB
func (usecase *GetContract) Execute(ctx context.Context, id *abi.ContractId) (*abi.Contract, error) {
	name, tag, err := utils.CheckExtractNameTag(id)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	}

	artifact, err := usecase.artifactDataAgent.FindOneByNameAndTag(ctx, name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	}

	return &abi.Contract{
		Id: &abi.ContractId{
			Name: name,
			Tag:  tag,
		},
		Abi:              artifact.Abi,
		Bytecode:         artifact.Bytecode,
		DeployedBytecode: artifact.DeployedBytecode,
	}, nil
}
