package builder

import (
	qkmclient "github.com/consensys/quorum-key-manager/pkg/client"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/business/use-cases/accounts"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

type accountUseCases struct {
	createAccountUC  usecases.CreateAccountUseCase
	getAccountUC     usecases.GetAccountUseCase
	searchAccountsUC usecases.SearchAccountsUseCase
	updateAccountUC  usecases.UpdateAccountUseCase
}

func newAccountUseCases(
	db store.DB,
	keyManagerClient qkmclient.EthClient,
	searchChainsUC usecases.SearchChainsUseCase,
	sendTxUC usecases.SendTxUseCase,
	getFaucetCandidateUC usecases.GetFaucetCandidateUseCase,
) *accountUseCases {
	searchAccountsUC := accounts.NewSearchAccountsUseCase(db)
	fundAccountUC := accounts.NewFundAccountUseCase(searchChainsUC, sendTxUC, getFaucetCandidateUC)

	return &accountUseCases{
		createAccountUC:  accounts.NewCreateAccountUseCase(db, searchAccountsUC, fundAccountUC, keyManagerClient),
		getAccountUC:     accounts.NewGetAccountUseCase(db),
		searchAccountsUC: searchAccountsUC,
		updateAccountUC:  accounts.NewUpdateAccountUseCase(db),
	}
}

func (u *accountUseCases) GetAccount() usecases.GetAccountUseCase {
	return u.getAccountUC
}

func (u *accountUseCases) SearchAccounts() usecases.SearchAccountsUseCase {
	return u.searchAccountsUC
}

func (u *accountUseCases) CreateAccount() usecases.CreateAccountUseCase {
	return u.createAccountUC
}

func (u *accountUseCases) UpdateAccount() usecases.UpdateAccountUseCase {
	return u.updateAccountUC
}
