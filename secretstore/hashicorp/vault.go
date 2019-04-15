package hashicorp

import (
	"github.com/hashicorp/vault/api"
)

// Hashicorps wraps a hashicorps client an manage the unsealing
type Hashicorps struct {
	Client             *api.Client
}

// NewHashicorps construct a new hashicorps vault given a configfile or nil
func NewHashicorps(config *api.Config) (*Hashicorps, error) {
	if config == nil {
		config = api.DefaultConfig()
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Hashicorps{
		Client: client,
	}, nil
}

// Store writes in the vault
func (hash *Hashicorps) Store(key, value string) (err error) {
	sec := NewSecret(key, value).SetClient(hash.Client)
	return sec.Update()
}

// Load reads in the vault
func (hash *Hashicorps) Load(key string) (value string, ok bool, err error) {
	sec := NewSecret(key, "").SetClient(hash.Client)
	return sec.GetValue()
}

// Delete removes a path in the vault
func (hash *Hashicorps) Delete(key string) (err error) {
	sec := NewSecret(key, "").SetClient(hash.Client)
	return sec.Delete()
}

// List returns the list of all secrets stored in the vault
func (hash *Hashicorps) List() (keys []string, err error) {
	sec := NewSecret("", "").SetClient(hash.Client)
	keys, err = sec.List("")
	return keys, err
}
