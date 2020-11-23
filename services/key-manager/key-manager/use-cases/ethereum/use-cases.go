package ethereum

import (
	"context"

	quorumtypes "github.com/consensys/quorum/core/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

//go:generate mockgen -source=use-cases.go -destination=mocks/use-cases.go -package=mocks

/**
Ethereum Use Cases
*/
type UseCases interface {
	CreateAccount() CreateAccountUseCase
	SignPayload() SignUseCase
	SignTransaction() SignTransactionUseCase
	SignQuorumPrivateTransaction() SignQuorumPrivateTransactionUseCase
	SignEEATransaction() SignEEATransactionUseCase
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, namespace, importedPrivKey string) (*entities.ETHAccount, error)
}

type SignUseCase interface {
	Execute(ctx context.Context, address, namespace, data string) (string, error)
}

type SignTransactionUseCase interface {
	Execute(ctx context.Context, address, namespace, chainID string, tx *ethtypes.Transaction) (string, error)
}

type SignQuorumPrivateTransactionUseCase interface {
	Execute(ctx context.Context, address, namespace string, tx *quorumtypes.Transaction) (string, error)
}

type SignEEATransactionUseCase interface {
	Execute(
		ctx context.Context,
		address, namespace string, chainID string,
		tx *ethtypes.Transaction,
		privateArgs *entities.PrivateETHTransactionParams,
	) (string, error)
}
