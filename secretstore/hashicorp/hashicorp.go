package hashicorp

import (
	"github.com/hashicorp/vault/api"
)

// HashiCorp wraps a hashicorps client an manage the unsealing
type HashiCorp struct {
	Client *api.Client
}

// NewHashiCorp construct a new hashicorps vault given a configfile or nil
func NewHashiCorp(config *api.Config) (*HashiCorp, error) {
	if config == nil {
		// This will read the environments variable
		config = api.DefaultConfig()
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &HashiCorp{
		Client: client,
	}, nil
}

// Store writes in the vault
func (hash *HashiCorp) Store(key, value string) (err error) {
	sec := NewSecret(key, value)
	sec.SetClient(hash.Client)
	return sec.Update()
}

// Load reads in the vault
func (hash *HashiCorp) Load(key string) (value string, ok bool, err error) {
	sec := NewSecret(key, "")
	sec.SetClient(hash.Client)
	return sec.GetValue()
}

// Delete removes a path in the vault
func (hash *HashiCorp) Delete(key string) (err error) {
	sec := NewSecret(key, "")
	sec.SetClient(hash.Client)
	return sec.Delete()
}

// List returns the list of all secrets stored in the vault
func (hash *HashiCorp) List() (keys []string, err error) {
	sec := NewSecret("", "")
	sec.SetClient(hash.Client)
	keys, err = sec.List("")
	return keys, err
}
