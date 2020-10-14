package builder

import (
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases/identity"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
)

func NewUseCases(db store.DB, client client.KeyManagerClient) usecases.IdentityUseCases {
	searchUC := identity.NewSearchIdentitiesUseCase(db)
	return &useCases{
		createIdentityUC: identity.NewCreateIdentityUseCase(db, searchUC, client),
		searchIdentityUC: searchUC,
	}
}

type useCases struct{
	createIdentityUC usecases.CreateIdentityUseCase
	searchIdentityUC usecases.SearchIdentitiesUseCase
}

func (u useCases) SearchIdentity() usecases.SearchIdentitiesUseCase {
	return u.searchIdentityUC
}

func (u useCases) CreateIdentity() usecases.CreateIdentityUseCase {
	return u.createIdentityUC
}

