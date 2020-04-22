package tessera

import "context"

//go:generate mockgen -destination=mock/mock.go -package=mock . Client

type Client interface {
	// StoreRaw stores "data" field of a transaction in Tessera privacy enclave
	// It returns a hash of a stored transaction that should be used instead of transaction data
	StoreRaw(ctx context.Context, endpoint string, data []byte, privateFrom string) (enclaveKey string, err error)

	// GetStatus returns status of Tessera enclave if it is up or an error if it is down
	GetStatus(ctx context.Context, endpoint string) (status string, err error)
}
