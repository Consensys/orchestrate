package contracts

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
)

const getTagsComponent = "use-cases.get-tags"

type getTagsUseCase struct {
	agent  store.TagAgent
	logger *log.Logger
}

func NewGetTagsUseCase(agent store.TagAgent) usecases.GetContractTagsUseCase {
	return &getTagsUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(getTagsComponent),
	}
}

func (uc *getTagsUseCase) Execute(ctx context.Context, name string) ([]string, error) {
	ctx = log.WithFields(ctx, log.Field("contract_name", name))
	names, err := uc.agent.FindAllByName(ctx, name)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getTagsComponent)
	}

	uc.logger.WithContext(ctx).Debug("get tags executed successfully")
	return names, nil
}
