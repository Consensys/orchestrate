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

// NewSecretFromKey creates a Secret from a key, it does not fetch the associated value.
func NewSecretFromKey(key string) *Secret {
	return &Secret{
		key:    key,
		value:  "",
		client: nil,
	}
}

// SetKey setter of attribute key for Secret struct object
func (s *Secret) SetKey(key string) *Secret {
	s.key = key
	return s
}

// SetValue setter of attribute value for Secret struct object
func (s *Secret) SetValue(value string) *Secret {
	s.value = value
	return s
}

// SetClient setter of attribute client for Secret struct object
func (s *Secret) SetClient(client *api.Client) *Secret {
	s.client = client
	return s
}

// SaveNew stores a new Secret in the vault
func (s *Secret) SaveNew() (err error) {
	fetched, err := s.GetValue()
	if fetched != "" {
		return fmt.Errorf("Secret %q already exists", s.key)
	}
	return s.Update()
}

// GetValue fetch the value from AWS SecretManager by key
func (s *Secret) GetValue() (string, error) {
	// Read secret from Vault
	logical := s.client.Logical()
	res, err := logical.Read(fmt.Sprintf("secret/secret/%v", s.key))
	if err != nil {
		return "", err
	}

	// When the secret is missing the client returns nil, nil.
	// We catch it here
	if res == nil {
		return "", fmt.Errorf("No secret for key %q", s.key)
	}
	s.value = res.Data["value"].(string)

	return s.value, nil
}

// Update the Secret value stored in the aws Secret manager
func (s *Secret) Update() error {
	// Load secret to Vault
	logical := s.client.Logical()
	_, err := logical.Write(
		fmt.Sprintf("secret/secret/%v", s.key),
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
	_, err := logical.Delete(fmt.Sprintf("secret/secret/%v", s.key))
	if err != nil {
		return err
	}

	return nil
}

// List retrieve all the keys availables in the Secret manager
func (s *Secret) List() ([]string, error) {
	logical := s.client.Logical()
	res, err := logical.List("secret/secret")
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
