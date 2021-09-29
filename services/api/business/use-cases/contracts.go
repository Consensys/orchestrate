package usecases

import (
	"context"

	"github.com/consensys/orchestrate/pkg/types/entities"
)

//go:generate mockgen -source=contracts.go -destination=mocks/contracts.go -package=mocks

type ContractUseCases interface {
	GetContractsCatalog() GetContractsCatalogUseCase
	GetContract() GetContractUseCase
	GetContractEvents() GetContractEventsUseCase
	GetContractMethodSignatures() GetContractMethodSignaturesUseCase
	GetContractMethods() GetContractMethodsUseCase
	GetContractTags() GetContractTagsUseCase
	SetContractCodeHash() SetContractCodeHashUseCase
	RegisterContract() RegisterContractUseCase
}

type GetContractsCatalogUseCase interface {
	Execute(ctx context.Context) ([]string, error)
}

type GetContractUseCase interface {
	Execute(ctx context.Context, name, tag string) (*entities.Contract, error)
}

type GetContractEventsUseCase interface {
	Execute(ctx context.Context, chainID, address, sighash string, indexedInputCount uint32) (abi string, eventsABI []string, err error)
}

type GetContractMethodSignaturesUseCase interface {
	Execute(ctx context.Context, name, tag, methodName string) ([]string, error)
}

type GetContractMethodsUseCase interface {
	Execute(ctx context.Context, chainID, address string, selector []byte) (abi string, methodsABI []string, err error)
}

type GetContractTagsUseCase interface {
	Execute(ctx context.Context, name string) ([]string, error)
}

type RegisterContractUseCase interface {
	Execute(ctx context.Context, contract *entities.Contract) error
}

type SetContractCodeHashUseCase interface {
	Execute(ctx context.Context, chainID, address string, hash string) error
}
