package builder

import (
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/business/use-cases/accounts"
	"github.com/consensys/orchestrate/services/api/store"
	qkmclient "github.com/consensys/quorum-key-manager/pkg/client"
)

type accountUseCases struct {
	createAccountUC  usecases.CreateAccountUseCase
	getAccountUC     usecases.GetAccountUseCase
	deleteAccountUC  usecases.DeleteAccountUseCase
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
		deleteAccountUC:  accounts.NewDeleteAccountUseCase(db, keyManagerClient),
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

func (u *accountUseCases) DeleteAccount() usecases.DeleteAccountUseCase {
	return u.deleteAccountUC
}
