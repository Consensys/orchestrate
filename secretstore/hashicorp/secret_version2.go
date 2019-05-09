package hashicorp

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// SecretKV2 contains a key/value secret
type SecretKV2 struct {
	key    string
	value  string
	client *api.Client
}

// NewSecretV2 creates a Secret from key and value
func NewSecretV2(key, value string) *SecretKV2 {
	return &SecretKV2{
		key:    key,
		value:  value,
		client: nil,
	}
}

// SetClient setter of attribute client for Secret struct object
func (s *SecretKV2) SetClient(client *api.Client) {
	s.client = client
}

// SaveNew stores a new Secret in the vault
func (s *SecretKV2) SaveNew() (err error) {
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
func (s *SecretKV2) GetValue() (value string, ok bool, err error) {
	// Read secret from Vault
	logical := s.client.Logical()
	res, err := logical.Read(
		fmt.Sprintf("secret/data/%v/%v", GetSecretPath(), s.key),
	)
	if err != nil {
		return "", false, err
	}

	// When the secret does not exist the client returns nil, nil.
	// When the secret is deleted the keystore returns an empty data.
	// We catch it here
	if res == nil || res.Data["data"] == nil {
		return "", false, nil
	}
	s.value = res.Data["data"].(map[string]interface{})["value"].(string)

	return s.value, true, nil
}

// Update the Secret value stored in the aws Secret manager
func (s *SecretKV2) Update() error {
	// Load secret to Vault
	logical := s.client.Logical()
	_, err := logical.Write(
		fmt.Sprintf("secret/data/%v/%v", GetSecretPath(), s.key),
		map[string]interface{}{
			"data": map[string]interface{}{"value": s.value},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete remove the key from the Secret manager
func (s *SecretKV2) Delete() error {
	// Delete secret in Vault
	logical := s.client.Logical()
	_, err := logical.Delete(
		fmt.Sprintf("secret/metadata/%v/%v", GetSecretPath(), s.key),
	)
	if err != nil {
		return err
	}

	return nil
}

// List retrieve all the keys availables in the Secret manager
func (s *SecretKV2) List(subPath string) ([]string, error) {

	logical := s.client.Logical()
	fullPath := fmt.Sprintf("secret/metadata/%v", GetSecretPath())

	if subPath != "" && subPath[0] == '/' {
		subPath = subPath[1:]
	}

	if subPath != "" {
		fullPath = fmt.Sprintf("%v/%v", fullPath, subPath)
	}

	println(fullPath)

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
