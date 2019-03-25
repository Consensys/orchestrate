package secretstore

// SecretStore is an interface implemented by helpers that store and retrive secrets
type SecretStore interface {
	// Store secret
	Store(key, value string) (err error)

	// Load secret
	Load(key string) (value string, err error)

	// Delete secret
	Delete(key string) (err error)

	// List secrets
	List() (keys []string, err error)
}
