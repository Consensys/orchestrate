package usecases

import (
	"context"

	signer "github.com/ethereum/go-ethereum/signer/core"
)

//go:generate mockgen -source=use-cases.go -destination=mocks/use-cases.go -package=mocks

type UseCases interface {
	SignTypedData() SignTypedDataUseCase
	VerifySignature() VerifySignatureUseCase
	VerifyTypedDataSignature() VerifyTypedDataSignatureUseCase
}

type SignTypedDataUseCase interface {
	Execute(ctx context.Context, address, namespace string, typedData *signer.TypedData) (string, error)
}

type VerifySignatureUseCase interface {
	Execute(ctx context.Context, address, signature, payload string) error
}

type VerifyTypedDataSignatureUseCase interface {
	Execute(ctx context.Context, address, signature string, typedData *signer.TypedData) error
}
