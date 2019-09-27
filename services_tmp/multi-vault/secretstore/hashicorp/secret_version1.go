package hashicorp

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// SecretKV1 contains a key/value secret
type SecretKV1 struct {
	key    string
	value  string
	client *api.Client
}

// NewSecretV1 creates a Secret from key and value
func NewSecretV1(key, value string) *SecretKV1 {
	return &SecretKV1{
		key:    key,
		value:  value,
		client: nil,
	}
}

// SetClient setter of attribute client for Secret struct object
func (s *SecretKV1) SetClient(client *api.Client) {
	s.client = client
}

// SaveNew stores a new Secret in the vault
func (s *SecretKV1) SaveNew() (err error) {
	fetched, _, err := s.GetValue()
	if fetched != "" {
		return errors.AlreadyExistsError("secret %q already exists", s.key).SetComponent(component)
	}
	if err != nil {
		return err
	}
	return s.Update()
}

// GetValue fetch the value from vault by key
func (s *SecretKV1) GetValue() (value string, ok bool, err error) {
	// Read secret from Vault
	logical := s.client.Logical()
	res, err := logical.Read(
		fmt.Sprintf("%v/%v/%v", GetMountPoint(), GetSecretPath(), s.key),
	)
	if err != nil {
		return "", false, errors.ConnectionError(err.Error()).SetComponent(component)
	}

	// When the secret is missing the client returns nil, nil.
	// We catch it here
	if res == nil {
		return "", false, nil
	}
	s.value = res.Data["value"].(string)

	return s.value, true, nil
}

// Update the Secret value stored in the vault
func (s *SecretKV1) Update() error {
	// Load secret to Vault
	logical := s.client.Logical()
	_, err := logical.Write(
		fmt.Sprintf("%v/%v/%v", GetMountPoint(), GetSecretPath(), s.key),
		map[string]interface{}{"value": s.value},
	)
	if err != nil {
		return errors.ConnectionError(err.Error()).SetComponent(component)
	}

	return nil
}

// Delete remove the key from the vault
func (s *SecretKV1) Delete() error {
	// Delete secret in Vault
	logical := s.client.Logical()
	_, err := logical.Delete(
		fmt.Sprintf("%v/%v/%v", GetMountPoint(), GetSecretPath(), s.key),
	)
	if err != nil {
		return errors.ConnectionError(err.Error()).SetComponent(component)
	}

	return nil
}

// List retrieve all the keys availables in the vault
func (s *SecretKV1) List(subPath string) ([]string, error) {

	logical := s.client.Logical()
	fullPath := fmt.Sprintf("%v/%v", GetMountPoint(), GetSecretPath())

	if subPath != "" && subPath[0] == '/' {
		subPath = subPath[1:]
	}

	if subPath != "" {
		fullPath = fmt.Sprintf("%v/%v", fullPath, subPath)
	}

	res, err := logical.List(fullPath)
	if err != nil {
		return nil, errors.ConnectionError(err.Error()).SetComponent(component)
	}

	if res == nil {
		return []string{}, nil
	}

	secrets := res.Data["keys"].([]interface{})
	rv := make([]string, len(secrets))
	for i, elem := range secrets {
		rv[i] = fmt.Sprintf("%v", elem)
	}

	return rv, nil
}
