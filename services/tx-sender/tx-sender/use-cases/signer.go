package usecases

import (
	"context"

	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//go:generate mockgen -source=signer.go -destination=mocks/signer.go -package=mocks

type SignETHTransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (signedRaw hexutil.Bytes, txHash *ethcommon.Hash, err error)
}

type SignEEATransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (raw hexutil.Bytes, txHash *ethcommon.Hash, err error)
}

type SignQuorumPrivateTransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (signedRaw hexutil.Bytes, txHash *ethcommon.Hash, err error)
}
