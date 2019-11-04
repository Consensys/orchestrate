package memory

import (
	"sync"
)

// SecretStore holds a pool of private keys in memory
type SecretStore struct {
	secrets *sync.Map
}

// NewSecretStore creates a new static signer
func NewSecretStore() *SecretStore {
	return &SecretStore{
		secrets: &sync.Map{},
	}
}

// Store secret
func (s *SecretStore) Store(key, value string) error {
	s.secrets.Store(key, value)
	return nil
}

// Load secret
func (s *SecretStore) Load(key string) (value string, ok bool, err error) {
	v, ok := s.secrets.Load(key)
	if !ok {
		return "", false, nil
	}
	return v.(string), true, nil
}

// Delete secret
func (s *SecretStore) Delete(key string) error {
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
