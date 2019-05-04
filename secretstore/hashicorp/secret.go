package hashicorp

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Secret contains a key/value secret
type Secret struct {
	key    string
	value  string
	client *api.Client
}

// NewSecret creates a Secret from key and value
func NewSecret(key, value string) *Secret {
	return &Secret{
		key:    key,
		value:  value,
		client: nil,
	}
}

// SetClient setter of attribute client for Secret struct object
func (s *Secret) SetClient(client *api.Client) *Secret {
	s.client = client
	return s
}

// SaveNew stores a new Secret in the vault
func (s *Secret) SaveNew() (err error) {
	fetched, _, err := s.GetValue()
	if fetched != "" {
		return fmt.Errorf("secret %q already exists", s.key)
	}
	if err != nil {
		return err
	}
	return s.Update()
}

// GetValue fetch the value from AWS SecretManager by key
func (s *Secret) GetValue() (value string, ok bool, err error) {
	// Read secret from Vault
	logical := s.client.Logical()
	res, err := logical.Read(
		fmt.Sprintf("%v/%v", GetSecretPath(), s.key),
	)
	if err != nil {
		return "", false, err
	}

	// When the secret is missing the client returns nil, nil.
	// We catch it here
	if res == nil {
		return "", false, nil
	}
	s.value = res.Data["value"].(string)

	return s.value, true, nil
}

// Update the Secret value stored in the aws Secret manager
func (s *Secret) Update() error {
	// Load secret to Vault
	logical := s.client.Logical()
	_, err := logical.Write(
		fmt.Sprintf("%v/%v", GetSecretPath(), s.key),
		map[string]interface{}{"value": s.value},
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete remove the key from the Secret manager
func (s *Secret) Delete() error {
	// Delete secret in Vault
	logical := s.client.Logical()
	_, err := logical.Delete(
		fmt.Sprintf("%v/%v", GetSecretPath(), s.key),
	)
	if err != nil {
		return err
	}

	return nil
}

// List retrieve all the keys availables in the Secret manager
func (s *Secret) List(subPath string) ([]string, error) {

	logical := s.client.Logical()
	fullPath := GetSecretPath()

	if subPath != "" && subPath[0] == '/' {
		subPath = subPath[1:]
	}

	if subPath != "" {
		fullPath = fmt.Sprintf("%v/%v", fullPath, subPath)
	}

	res, err := logical.List(fullPath)
	if err != nil {
		return nil, err
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
