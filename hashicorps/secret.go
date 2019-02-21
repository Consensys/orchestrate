package hashicorps

import (
	"github.com/hashicorp/vault/api"
	"strings"
	"fmt"
)

// Secret contains a key value secret
type Secret struct {
	key string
	value string
	client *api.Client
}

// NewSecret is the default constructor for Secret
func NewSecret() *Secret {
	return &Secret{}
}

// CreateSecret creates a Secret from key and value
func CreateSecret(key, value string) (*Secret) {
	return &Secret{
		key: key,
		value: value,
		client: nil,
	}
}

// SecretFromKey creates a secret from a key, it does not fetch the associated value.
func SecretFromKey(key string) (*Secret) {
	return &Secret{
		key: key,
		value: "",
		client: nil,
	}
}

// SetKey setter of attribute key for Secret struct object
func (sec *Secret) SetKey(key string) *Secret {
	sec.key = key
	return sec
}

// SetValue setter of attribute value for Secret struct object
func (sec *Secret) SetValue(value string) *Secret {
	sec.value = value
	return sec	
}

// SetClient setter of attribute client for Secret struct object
func (sec *Secret) SetClient(client *api.Client) *Secret {
	sec.client = client
	return sec	
}

// SaveNew stores a new secret in the vault
func (sec *Secret) SaveNew() (res *api.Secret, err error) {

	fetched, err := sec.GetValue()
	if fetched != "" {
		return nil, fmt.Errorf("This secret already exists : " + sec.key)
	}

	return sec.Update()
}

// GetValue fetch the value from AWS SecretManager by key
func (sec *Secret) GetValue() (string, error) {

	log := sec.client.Logical()
	res, err := log.Read(
		strings.Join([]string{"secret", sec.key}, "/"),
	)

	if err != nil {
		return "", err
	}

	sec.value = res.Data["value"].(string)

	return sec.value, nil
}

// Update the secret value stored in the aws secret manager
func (sec *Secret) Update() (*api.Secret, error) {

	log := sec.client.Logical()
	res, err := log.Write(
		strings.Join([]string{"secret", sec.key}, "/"),
		map[string]interface{}{ "value": sec.value },
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}


// Delete remove the key from the secret manager
func (sec *Secret) Delete() (*api.Secret, error) {

	log := sec.client.Logical()
	res, err := log.Delete(
		strings.Join([]string{"secret", sec.key}, "/"),
	)

	if err != nil {
		return nil, err
	}

	return res, nil	
} 

// List retrieve all the keys availables in the secret manager
func (sec *Secret) List() ([]string, error) {

	log := sec.client.Logical()
	res, err := log.List("secret")

	if err != nil {
		return nil, err
	}

	list := make([]string, len(res.Data))
	i := 0
	for address := range res.Data {
		list[i] = address
		i++
	}

	return list, nil	
} 