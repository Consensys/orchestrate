package memory

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

// SecretStore holds a pool of private keys in memory
type SecretStore struct {
	secrets    *sync.Map
	KeyBuilder *multitenancy.KeyBuilder
}

// NewSecretStore creates a new static signer
func NewSecretStore(keyBuilder *multitenancy.KeyBuilder) *SecretStore {
	return &SecretStore{
		secrets:    &sync.Map{},
		KeyBuilder: keyBuilder,
	}
}

// Store secret
func (s *SecretStore) Store(ctx context.Context, key, value string) error {
	v, ok := s.secrets.Load(key)

	if ok {
		if v == value {
			return nil
		}
		return errors.AlreadyExistsError("A different secret already exists for key: %v", key).ExtendComponent(component)
	}

	s.secrets.Store(key, value)
	return nil
}

// Load secret
func (s *SecretStore) Load(ctx context.Context, key string) (value string, ok bool, e error) {
	v, ok := s.secrets.Load(key)
	if ok {
		return v.(string), true, nil
	}

	return "", false, nil
}

// Delete secret
func (s *SecretStore) Delete(ctx context.Context, key string) error {
	s.secrets.Delete(key)
	return nil
}

// List secret
func (s *SecretStore) List() (keys []string, err error) {
	keys = []string{}
	s.secrets.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys, nil
}
