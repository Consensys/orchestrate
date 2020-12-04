package usecases

import (
	"context"

	signer "github.com/ethereum/go-ethereum/signer/core"
)

//go:generate mockgen -source=use-cases.go -destination=mocks/use-cases.go -package=mocks

type UseCases interface {
	SignTypedData() SignTypedDataUseCase
}

type SignTypedDataUseCase interface {
	Execute(ctx context.Context, address, namespace string, typedData *signer.TypedData) (string, error)
}
