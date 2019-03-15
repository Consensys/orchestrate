package secretstore

import (
	"github.com/hashicorp/vault/api"
	"sync"
)

// Hashicorps wraps a hashicorps client an manage the unsealing
type Hashicorps struct {
	client *api.Client
	creds *credentials
	retrieveSecretOnce *sync.Once
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

	creds := &credentials{}
	return &Hashicorps{
		client: client,
		creds: creds,
	}, nil
}

// Init manages vault token auth and unsealing
func (hash *Hashicorps) Init(credsStore *AWS, tokenName string) (err error) {

	err = hash.creds.FetchFromAWS(credsStore, tokenName)
	if err != nil {
		return err
	}

	hash.creds.AttachTo(hash.client)

	err = hash.creds.Unseal(hash.client)
	if err != nil {
		return err
	}

	return nil
}

// Store writes in the vault
func (hash *Hashicorps) Store(key, value string) (err error) {
	sec := NewVaultSecret().SetKey(key).SetValue(value).SetClient(hash.client)
	return sec.Update()
}

// Load reads in the vault
func (hash *Hashicorps) Load(key string) (value string, err error) {
	sec := NewVaultSecret().SetKey(key).SetClient(hash.client)
	res, err := sec.GetValue()
	if err != nil {
		return "", err
	}
	return res, nil
}

// Delete removes a path in the vault
func (hash *Hashicorps) Delete(key string) (err error) {
	sec := NewVaultSecret().SetKey(key).SetClient(hash.client)
	return sec.Delete()
}

// List returns the list of all secrets stored in the vault
func (hash *Hashicorps) List() (keys []string, err error) {
	sec := NewVaultSecret().SetClient(hash.client)
	keys, err = sec.List()
	return keys, err
}
