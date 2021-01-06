package contracts

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
)

const getContractComponent = "use-cases.get-contract"

type getContractUseCase struct {
	agent store.ArtifactAgent
}

func NewGetContractUseCase(agent store.ArtifactAgent) usecases.GetContractUseCase {
	return &getContractUseCase{
		agent: agent,
	}
}

// Execute gets a contract from DB
func (usecase *getContractUseCase) Execute(ctx context.Context, id *entities.ContractID) (*entities.Contract, error) {
	logger := log.WithContext(ctx).WithField("contract", id.Short())
	logger.Debug("get contract is starting...")

	artifact, err := usecase.agent.FindOneByNameAndTag(ctx, id.Name, id.Tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	}

	contract := &entities.Contract{
		ID:               *id,
		ABI:              artifact.ABI,
		Bytecode:         artifact.Bytecode,
		DeployedBytecode: artifact.DeployedBytecode,
		Methods:          []entities.Method{},
	}

	contractABI, err := contract.ToABI()
	if err != nil {
		errMessage := "failed to get contract ABI"
		logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(getMethodSignaturesComponent)
	}

	for _, method := range contractABI.Methods {
		contract.Methods = append(contract.Methods, entities.Method{
			Signature: method.Sig(),
		})
	}
	for _, event := range contractABI.Events {
		contract.Events = append(contract.Events, entities.Event{
			Signature: event.Sig(),
		})
	}
	contract.Constructor = entities.Method{Signature: contractABI.Constructor.Sig()}

	logger.Debug("get contract executed successfully")
	return contract, nil
}
