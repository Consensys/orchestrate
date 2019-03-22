package secretstore

// SecretStore is an interface implemented by helpers that store and retrive secrets
type SecretStore interface {
	Store(key, value string) (err error)
	Load(key string) (value string, err error)
	Delete(key string) (err error)
	List() (keys []string, err error)
}

