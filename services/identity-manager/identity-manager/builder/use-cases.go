package builder

import (
	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/use-cases/account"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	client3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
)

func NewUseCases(db store.DB, keyManagerClient client.KeyManagerClient, registryClient client2.ChainRegistryClient,
	txSchedulerClient client3.TransactionSchedulerClient) usecases.AccountUseCases {
	searchIdentityUC := account.NewSearchAccountsUseCase(db)
	fundingIdentityUC := account.NewFundingAccountUseCase(registryClient, txSchedulerClient)
	return &useCases{
		createAccountUC:  account.NewCreateAccountUseCase(db, searchIdentityUC, fundingIdentityUC, keyManagerClient),
		getAccountUC:     account.NewGetAccountUseCase(db),
		searchAccountsUC: searchIdentityUC,
		updateAccountUC:  account.NewUpdateAccountUseCase(db),
		fundingAccountUC: fundingIdentityUC,
	}
}

type useCases struct {
	createAccountUC  usecases.CreateAccountUseCase
	getAccountUC     usecases.GetAccountUseCase
	searchAccountsUC usecases.SearchAccountsUseCase
	updateAccountUC  usecases.UpdateAccountUseCase
	fundingAccountUC usecases.FundingAccountUseCase
}

func (u useCases) GetAccount() usecases.GetAccountUseCase {
	return u.getAccountUC
}

func (u useCases) SearchAccounts() usecases.SearchAccountsUseCase {
	return u.searchAccountsUC
}

func (u useCases) CreateAccount() usecases.CreateAccountUseCase {
	return u.createAccountUC
}

func (u useCases) UpdateAccount() usecases.UpdateAccountUseCase {
	return u.updateAccountUC
}

func (u useCases) FundingAccount() usecases.FundingAccountUseCase {
	return u.fundingAccountUC
}