package builder

import (
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/business/use-cases/chains"
	"github.com/consensys/orchestrate/services/api/store"
)

type chainUseCases struct {
	registerChainUC usecases.RegisterChainUseCase
	updateChainUC   usecases.UpdateChainUseCase
	getChainUC      usecases.GetChainUseCase
	searchChainsUC  usecases.SearchChainsUseCase
	deleteChainUC   usecases.DeleteChainUseCase
}

func newChainUseCases(db store.DB, ec ethclient.Client) *chainUseCases {
	searchChainsUC := chains.NewSearchChainsUseCase(db)
	getChainUC := chains.NewGetChainUseCase(db)

	return &chainUseCases{
		registerChainUC: chains.NewRegisterChainUseCase(db, searchChainsUC, ec),
		updateChainUC:   chains.NewUpdateChainUseCase(db, getChainUC),
		getChainUC:      getChainUC,
		searchChainsUC:  searchChainsUC,
		deleteChainUC:   chains.NewDeleteChainUseCase(db, getChainUC),
	}
}

func (u *chainUseCases) RegisterChain() usecases.RegisterChainUseCase {
	return u.registerChainUC
}

func (u *chainUseCases) UpdateChain() usecases.UpdateChainUseCase {
	return u.updateChainUC
}

func (u *chainUseCases) GetChain() usecases.GetChainUseCase {
	return u.getChainUC
}

func (u *chainUseCases) SearchChains() usecases.SearchChainsUseCase {
	return u.searchChainsUC
}

func (u *chainUseCases) DeleteChain() usecases.DeleteChainUseCase {
	return u.deleteChainUC
}
