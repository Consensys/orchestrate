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
func (usecase *getContractUseCase) Execute(ctx context.Context, name, tag string) (*entities.Contract, error) {
	logger := log.WithContext(ctx).WithField("name", name).WithField("tag", tag)
	logger.Debug("getting contract")

	artifact, err := usecase.agent.FindOneByNameAndTag(ctx, name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	}

	contract := &entities.Contract{
		Name:             name,
		Tag:              tag,
		ABI:              artifact.ABI,
		Bytecode:         artifact.Bytecode,
		DeployedBytecode: artifact.DeployedBytecode,
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

	logger.Debug("contract found successfully")
	return contract, nil
}
