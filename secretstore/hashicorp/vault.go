package hashicorp

import (
	"sync"

	"github.com/hashicorp/vault/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/aws"
)

// Hashicorps wraps a hashicorps client an manage the unsealing
type Hashicorps struct {
	Client             *api.Client
	creds              *credentials
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
		creds:  creds,
	}, nil
}

// InitVault fetches the new token and sets the values in AWS
func (hash *Hashicorps) InitVault() (err error) {
	err = hash.creds.FetchFromVaultInit(hash.Client)
	if err != nil {
		return err
	}

	hash.SetToken(hash.creds.Token)

	err = hash.Unseal(hash.creds.Keys[0])
	if err != nil {
		return err
	}

	return nil
}

// SendToCredStore stores the vault credentials in AWS
func (hash *Hashicorps) SendToCredStore(credsStore *aws.AWS, tokenName string) (err error) {
	err = hash.creds.SendToAWS(credsStore, tokenName)
	if err != nil {
		return err
	}

	return nil
}

// InitFromAWS manages vault token auth and unsealing
func (hash *Hashicorps) InitFromAWS(credsStore *aws.AWS, tokenName string) (err error) {
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

// Unseal the vault
// Warning call Unseal is Unsafe
func (hash *Hashicorps) Unseal(unsealKey string) (err error) {
	sys := hash.Client.Sys()
	sys.Unseal(unsealKey)
	if err != nil {
		return err
	}

	return nil
}

// SetToken authorize the client to access vault
// Warning call SetToken is Unsafe
func (hash *Hashicorps) SetToken(token string) {
	hash.Client.SetToken(token)
}

// Store writes in the vault
func (hash *Hashicorps) Store(key, value string) (err error) {
	sec := NewSecret(key, value).SetClient(hash.Client)
	return sec.Update()
}

// Load reads in the vault
func (hash *Hashicorps) Load(key string) (value string, ok bool, err error) {
	sec := NewSecret(key, "").SetClient(hash.Client)
	res, err := sec.GetValue()
	if err != nil {
		return "", false, err
	}
	return res, ok, nil
}

// Delete removes a path in the vault
func (hash *Hashicorps) Delete(key string) (err error) {
	sec := NewSecret(key, "").SetClient(hash.Client)
	return sec.Delete()
}

// List returns the list of all secrets stored in the vault
func (hash *Hashicorps) List() (keys []string, err error) {
	sec := NewSecret("", "").SetClient(hash.Client)
	keys, err = sec.List()
	return keys, err
}
