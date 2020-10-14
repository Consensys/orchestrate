package services

import "context"

//go:generate mockgen -source=secret_store.go -destination=mocks/secret_store.go -package=mocks

// SecretStore is an interface implemented by helpers that store and retrieve secrets
type SecretStore interface {
	// Store secret
	Store(ctx context.Context, key, value string) (err error)

	// Load secret
	Load(ctx context.Context, key string) (value string, ok bool, err error)

	// Delete secret
	Delete(ctx context.Context, key string) (err error)

	// List secrets
	List() (keys []string, err error)
}
