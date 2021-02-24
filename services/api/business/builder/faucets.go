package builder

import (
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/business/use-cases/faucets"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

type faucetUseCases struct {
	registerFaucetUC usecases.RegisterFaucetUseCase
	updateFaucetUC   usecases.UpdateFaucetUseCase
	getFaucetUC      usecases.GetFaucetUseCase
	searchFaucetUC   usecases.SearchFaucetsUseCase
	deleteFaucetUC   usecases.DeleteFaucetUseCase
}

func newFaucetUseCases(db store.DB) *faucetUseCases {
	searchFaucetsUC := faucets.NewSearchFaucets(db)

	return &faucetUseCases{
		registerFaucetUC: faucets.NewRegisterFaucetUseCase(db, searchFaucetsUC),
		updateFaucetUC:   faucets.NewUpdateFaucetUseCase(db),
		getFaucetUC:      faucets.NewGetFaucetUseCase(db),
		searchFaucetUC:   searchFaucetsUC,
		deleteFaucetUC:   faucets.NewDeleteFaucetUseCase(db),
	}
}

func (u *faucetUseCases) RegisterFaucet() usecases.RegisterFaucetUseCase {
	return u.registerFaucetUC
}

func (u *faucetUseCases) UpdateFaucet() usecases.UpdateFaucetUseCase {
	return u.updateFaucetUC
}

func (u *faucetUseCases) GetFaucet() usecases.GetFaucetUseCase {
	return u.getFaucetUC
}

func (u *faucetUseCases) SearchFaucets() usecases.SearchFaucetsUseCase {
	return u.searchFaucetUC
}

func (u *faucetUseCases) DeleteFaucet() usecases.DeleteFaucetUseCase {
	return u.deleteFaucetUC
}
