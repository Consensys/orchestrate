package contracts

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getTagsComponent = "use-cases.get-tags"

type getTagsUseCase struct {
	agent store.TagAgent
}

func NewGetTagsUseCase(agent store.TagAgent) usecases.GetContractTagsUseCase {
	return &getTagsUseCase{
		agent: agent,
	}
}

func (usecase *getTagsUseCase) Execute(ctx context.Context, name string) ([]string, error) {
	log.WithContext(ctx).Debug("get tags starting...")
	names, err := usecase.agent.FindAllByName(ctx, name)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getTagsComponent)
	}

	log.WithContext(ctx).Debug("get tags executed successfully")
	return names, nil
}
