package usecases

import (
	"context"

	signer "github.com/ethereum/go-ethereum/signer/core"
)

//go:generate mockgen -source=ethereum.go -destination=mocks/ethereum.go -package=mocks

type ETHUseCases interface {
	SignTypedData() SignTypedDataUseCase
	VerifySignature() VerifyETHSignatureUseCase
	VerifyTypedDataSignature() VerifyTypedDataSignatureUseCase
}

type SignTypedDataUseCase interface {
	Execute(ctx context.Context, address, namespace string, typedData *signer.TypedData) (string, error)
}

type VerifyETHSignatureUseCase interface {
	Execute(ctx context.Context, address, signature, payload string) error
}

type VerifyTypedDataSignatureUseCase interface {
	Execute(ctx context.Context, address, signature string, typedData *signer.TypedData) error
}
