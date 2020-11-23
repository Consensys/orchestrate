package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/store"
)

const getTagsComponent = component + ".get-tags"

//go:generate mockgen -source=get_tags.go -destination=mocks/mock_get_tags.go -package=mocks

type GetTagsUseCase interface {
	Execute(ctx context.Context, name string) ([]string, error)
}

// GetTags is a use case to get all tags
type GetTags struct {
	tagDataAgent store.TagDataAgent
}

// NewGetTags creates a new GetTags
func NewGetTags(tagDataAgent store.TagDataAgent) *GetTags {
	return &GetTags{
		tagDataAgent: tagDataAgent,
	}
}

// Execute gets all contract tags from DB
func (usecase *GetTags) Execute(ctx context.Context, name string) ([]string, error) {
	names, err := usecase.tagDataAgent.FindAllByName(ctx, name)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getTagsComponent)
	}

	return names, nil
}
