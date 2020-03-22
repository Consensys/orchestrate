package memory

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/error"
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
func (s *SecretStore) Store(ctx context.Context, rawKey, value string) error {
	key, err := s.KeyBuilder.BuildKey(ctx, rawKey)
	if err != nil {
		return err.(*ierror.Error).ExtendComponent(component)
	}
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
func (s *SecretStore) Load(ctx context.Context, rawKey string) (value string, ok bool, err error) {
	key, err := s.KeyBuilder.BuildKey(ctx, rawKey)
	if err != nil {
		return "", false, err.(*ierror.Error).ExtendComponent(component)
	}
	v, ok := s.secrets.Load(key)
	if !ok {
		return "", false, nil
	}
	return v.(string), true, nil
}

// Delete secret
func (s *SecretStore) Delete(ctx context.Context, rawKey string) error {
	key, err := s.KeyBuilder.BuildKey(ctx, rawKey)
	if err != nil {
		return err.(*ierror.Error).ExtendComponent(component)
	}
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
