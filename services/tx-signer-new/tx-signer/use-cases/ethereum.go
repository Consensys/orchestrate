package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"github.com/golang/protobuf/proto"
)

//go:generate mockgen -source=ethereum.go -destination=mocks/ethereum.go -package=mocks

type EthereumUseCases interface {
	SignTransaction() SignTransactionUseCase
	SignEEATransaction() SignEEATransactionUseCase
	SendEnvelope() SendEnvelopeUseCase
}

type SignTransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (signedRaw string, txHash string, err error)
}

type SignEEATransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (signedRaw string, txHash string, err error)
}

type SendEnvelopeUseCase interface {
	Execute(ctx context.Context, protoMessage proto.Message, topic, partitionKey string) error
}
