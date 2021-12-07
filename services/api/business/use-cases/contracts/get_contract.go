package contracts

import (
	"context"

	"github.com/consensys/orchestrate/services/api/store"
	"github.com/consensys/quorum/common/hexutil"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
)

const getContractComponent = "use-cases.get-contract"

type getContractUseCase struct {
	agent  store.ArtifactAgent
	logger *log.Logger
}

func NewGetContractUseCase(agent store.ArtifactAgent) usecases.GetContractUseCase {
	return &getContractUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(getContractComponent),
	}
}

// Execute gets a contract from DB
func (uc *getContractUseCase) Execute(ctx context.Context, name, tag string) (*entities.Contract, error) {
	ctx = log.WithFields(ctx, log.Field("contract_name", name), log.Field("contract_tag", name))
	logger := uc.logger.WithContext(ctx)

	artifact, err := uc.agent.FindOneByNameAndTag(ctx, name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	}

	contract := &entities.Contract{
		Name:             name,
		Tag:              tag,
		ABI:              artifact.ABI,
		Bytecode:         hexutil.MustDecode(artifact.Bytecode),
		DeployedBytecode: hexutil.MustDecode(artifact.DeployedBytecode),
	}

	contractABI, err := contract.ToABI()
	if err != nil {
		errMessage := "failed to get contract ABI"
		logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(getMethodSignaturesComponent)
	}

	// nolint
	for _, method := range contractABI.Methods {
		contract.Methods = append(contract.Methods, entities.Method{
			Signature: method.Sig(),
		})
	}

	// nolint
	for _, event := range contractABI.Events {
		contract.Events = append(contract.Events, entities.Event{
			Signature: event.Sig(),
		})
	}

	if contractABI.Constructor.Sig() == "" {
		contract.Constructor = entities.Method{Signature: "()"}
	} else {
		contract.Constructor = entities.Method{Signature: contractABI.Constructor.Sig()}
	}

	logger.Debug("contract found successfully")
	return contract, nil
}
