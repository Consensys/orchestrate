package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

const setCodeHashComponent = component + ".set-codehash"

//go:generate mockgen -source=set_codehash.go -destination=mocks/mock_set_codehash.go -package=mocks

type SetCodeHashUseCase interface {
	Execute(ctx context.Context, account *common.AccountInstance, hash string) error
}

// SetCodeHash is a use case to set the codehash of a contract
type SetCodeHash struct {
	codehashDataAgent store.CodeHashDataAgent
}

// NewSetCodeHash creates a new SetCodeHash
func NewSetCodeHash(codehashDataAgent store.CodeHashDataAgent) *SetCodeHash {
	return &SetCodeHash{
		codehashDataAgent: codehashDataAgent,
	}
}

// Execute sets the codehash of a contract in DB
func (usecase *SetCodeHash) Execute(ctx context.Context, account *common.AccountInstance, hash string) error {
	codehash := &models.CodehashModel{
		ChainID:  account.GetChainId(),
		Address:  account.GetAccount(),
		Codehash: hash,
	}
	err := usecase.codehashDataAgent.Insert(ctx, codehash)
	if err != nil {
		return errors.FromError(err).ExtendComponent(setCodeHashComponent)
	}

	return nil
}
