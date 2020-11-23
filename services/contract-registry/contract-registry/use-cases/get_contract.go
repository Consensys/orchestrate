package usecases

import (
	"context"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/contract-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/store"
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
	logger := log.WithContext(ctx).WithField("contract", id)
	logger.Debug("getting contract")

	name, tag, err := utils.CheckExtractNameTag(id)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	}

	artifact, err := usecase.artifactDataAgent.FindOneByNameAndTag(ctx, name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	}

	contract := &abi.Contract{
		Id: &abi.ContractId{
			Name: name,
			Tag:  tag,
		},
		Abi:              artifact.Abi,
		Bytecode:         artifact.Bytecode,
		DeployedBytecode: artifact.DeployedBytecode,
		Methods:          []*abi.Method{},
	}

	contractABI, err := contract.ToABI()
	if err != nil {
		errMessage := "failed to get contract ABI"
		logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(getMethodSignaturesComponent)
	}

	for _, method := range contractABI.Methods {
		contract.Methods = append(contract.Methods, &abi.Method{
			Signature: method.Sig(),
		})
	}
	contract.Constructor = &abi.Method{Signature: contractABI.Constructor.Sig()}

	return contract, nil
}
