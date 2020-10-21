package ethereum

import (
	"context"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
)

//go:generate mockgen -source=use-cases.go -destination=mocks/use-cases.go -package=mocks

/**
Ethereum Use Cases
*/
type UseCases interface {
	CreateAccount() CreateAccountUseCase
	SignPayload() SignUseCase
	SignTransaction() SignTransactionUseCase
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, namespace, importedPrivKey string) (*entities.ETHAccount, error)
}

type SignUseCase interface {
	Execute(ctx context.Context, address, namespace, data string) (string, error)
}

type SignTransactionUseCase interface {
	Execute(ctx context.Context, address, namespace string, chainID *big.Int, tx *ethtypes.Transaction) (string, error)
}
