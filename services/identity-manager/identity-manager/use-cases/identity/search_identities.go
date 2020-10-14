package identity

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
)

const searchJobsComponent = "use-cases.search-identities"

// searchJobsUseCase is a use case to search jobs
type searchIdentitiesUseCase struct {
	db store.DB
}

func NewSearchIdentitiesUseCase(db store.DB) usecases.SearchIdentitiesUseCase {
	return &searchIdentitiesUseCase{
		db: db,
	}
}

// Execute search jobs
func (uc *searchIdentitiesUseCase) Execute(ctx context.Context, filters *entities.IdentityFilters, tenants []string) ([]*entities.Identity, error) {
	log.WithContext(ctx).WithField("filters", filters).
		WithField("tenants", tenants).
		Debug("searching identities")

	models, err := uc.db.Identity().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchJobsComponent)
	}

	var resp []*entities.Identity
	for _, model := range models {
		iden := parsers.NewIdentityEntityFromModels(model)
		resp = append(resp, iden)
	}

	log.WithContext(ctx).Debug("identities found successfully")
	return resp, nil
}
