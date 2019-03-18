package secretstore

import (
	"github.com/hashicorp/vault/api"
	"sync"
)

// Hashicorps wraps a hashicorps client an manage the unsealing
type Hashicorps struct {
	Client *api.Client
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
		Client: client,
		creds: creds,
	}, nil
}

// InitVault fetches the new token and sets the values in AWS
func (hash *Hashicorps) InitVault(credsStore *AWS, tokenName string) (err error) {

	err = hash.creds.FetchFromVaultInit(hash.Client)
	if err != nil {
		return err
	}

	err = hash.creds.SendToAWS(credsStore, tokenName)

	return nil

}

// InitFromAWS manages vault token auth and unsealing
func (hash *Hashicorps) InitFromAWS(credsStore *AWS, tokenName string) (err error) {

	err = hash.creds.FetchFromAWS(credsStore, tokenName)
	if err != nil {
		return err
	}

	hash.creds.AttachTo(hash.Client)

	err = hash.creds.Unseal(hash.Client)
	if err != nil {
		return err
	}

	return nil
}

// Unseal [UNSAFE] the vault 
func (hash *Hashicorps) Unseal(unsealKey string) {
	sys := hash.Client.Sys()
	sys.Unseal(unsealKey)
}

// SetToken [UNSAFE] authorize the client to access vault
func (hash *Hashicorps) SetToken(token string) {
	hash.Client.SetToken(token)
}

// Store writes in the vault
func (hash *Hashicorps) Store(key, value string) (err error) {
	sec := NewVaultSecret().SetKey(key).SetValue(value).SetClient(hash.Client)
	return sec.Update()
}

// Load reads in the vault
func (hash *Hashicorps) Load(key string) (value string, err error) {
	sec := NewVaultSecret().SetKey(key).SetClient(hash.Client)
	res, err := sec.GetValue()
	if err != nil {
		return "", err
	}
	return res, nil
}

// Delete removes a path in the vault
func (hash *Hashicorps) Delete(key string) (err error) {
	sec := NewVaultSecret().SetKey(key).SetClient(hash.Client)
	return sec.Delete()
}

// List returns the list of all secrets stored in the vault
func (hash *Hashicorps) List() (keys []string, err error) {
	sec := NewVaultSecret().SetClient(hash.Client)
	keys, err = sec.List()
	return keys, err
}
