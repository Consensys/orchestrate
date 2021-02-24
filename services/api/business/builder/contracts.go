package builder

import (
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/business/use-cases/contracts"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

type contractUseCases struct {
	GetContractsCatalogUC       usecases.GetContractsCatalogUseCase
	getContractEvents           usecases.GetContractEventsUseCase
	getContractMethodSignatures usecases.GetContractMethodSignaturesUseCase
	getContractMethods          usecases.GetContractMethodsUseCase
	getContractTags             usecases.GetContractTagsUseCase
	setContractCodeHash         usecases.SetContractCodeHashUseCase
	registerContractUC          usecases.RegisterContractUseCase
	getContractUC               usecases.GetContractUseCase
}

func newContractUseCases(db store.DB) *contractUseCases {
	getContractUC := contracts.NewGetContractUseCase(db.Artifact())

	return &contractUseCases{
		registerContractUC:          contracts.NewRegisterContractUseCase(db),
		getContractUC:               getContractUC,
		GetContractsCatalogUC:       contracts.NewGetCatalogUseCase(db.Repository()),
		getContractEvents:           contracts.NewGetEventsUseCase(db.Event()),
		getContractMethodSignatures: contracts.NewGetMethodSignaturesUseCase(getContractUC),
		getContractMethods:          contracts.NewGetMethodsUseCase(db.Method()),
		getContractTags:             contracts.NewGetTagsUseCase(db.Tag()),
		setContractCodeHash:         contracts.NewSetCodeHashUseCase(db.CodeHash()),
	}
}

func (u *contractUseCases) GetContract() usecases.GetContractUseCase {
	return u.getContractUC
}

func (u *contractUseCases) RegisterContract() usecases.RegisterContractUseCase {
	return u.registerContractUC
}

func (u *contractUseCases) GetContractsCatalog() usecases.GetContractsCatalogUseCase {
	return u.GetContractsCatalogUC
}

func (u *contractUseCases) GetContractEvents() usecases.GetContractEventsUseCase {
	return u.getContractEvents
}

func (u *contractUseCases) GetContractMethodSignatures() usecases.GetContractMethodSignaturesUseCase {
	return u.getContractMethodSignatures
}

func (u *contractUseCases) GetContractMethods() usecases.GetContractMethodsUseCase {
	return u.getContractMethods
}

func (u *contractUseCases) GetContractTags() usecases.GetContractTagsUseCase {
	return u.getContractTags
}

func (u *contractUseCases) SetContractCodeHash() usecases.SetContractCodeHashUseCase {
	return u.setContractCodeHash
}
