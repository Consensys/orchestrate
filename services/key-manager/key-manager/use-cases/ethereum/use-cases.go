package ethereum

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
)

//go:generate mockgen -source=use-cases.go -destination=mocks/use-cases.go -package=mocks

/**
Ethereum Use Cases
*/
type UseCases interface {
	CreateAccount() CreateAccountUseCase
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, namespace, importedPrivKey string) (*entities.ETHAccount, error)
}
