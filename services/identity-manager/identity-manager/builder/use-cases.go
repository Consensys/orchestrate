package builder

import (
	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases/identity"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	client3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

func NewUseCases(db store.DB, keyManagerClient client.KeyManagerClient, registryClient client2.ChainRegistryClient,
	txSchedulerClient client3.TransactionSchedulerClient) usecases.IdentityUseCases {
	searchIdentityUC := identity.NewSearchIdentitiesUseCase(db)
	fundingIdentityUC := identity.NewFundingIdentityUseCase(registryClient, txSchedulerClient)
	return &useCases{
		createIdentityUC:  identity.NewCreateIdentityUseCase(db, searchIdentityUC, fundingIdentityUC, keyManagerClient),
		searchIdentityUC:  searchIdentityUC,
		fundingIdentityUC: fundingIdentityUC,
	}
}

type useCases struct {
	createIdentityUC  usecases.CreateIdentityUseCase
	searchIdentityUC  usecases.SearchIdentitiesUseCase
	fundingIdentityUC usecases.FundingIdentityUseCase
}

func (u useCases) SearchIdentity() usecases.SearchIdentitiesUseCase {
	return u.searchIdentityUC
}

func (u useCases) CreateIdentity() usecases.CreateIdentityUseCase {
	return u.createIdentityUC
}

func (u useCases) FundingIdentity() usecases.FundingIdentityUseCase {
	return u.fundingIdentityUC
}
