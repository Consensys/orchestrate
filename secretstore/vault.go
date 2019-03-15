package secretstore

import (
	"github.com/hashicorp/vault/api"
	"sync"
)

type Hashicorps struct {
	client *api.Client
	creds *credentials
	retrieveSecretOnce *sync.Once
}

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

func (hash *Hashicorps) Store(key, value string) (err error) {
	sec := NewVaultSecret().SetKey(key).SetValue(value).SetClient(hash.client)
	return sec.Update()
}

func (hash *Hashicorps) Load(key string) (value string, err error) {
	sec := NewVaultSecret().SetKey(key).SetClient(hash.client)
	res, err := sec.GetValue()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (hash *Hashicorps) Delete(key string) (err error) {
	sec := NewVaultSecret().SetKey(key).SetClient(hash.client)
	return sec.Delete()
}

func (hash *Hashicorps) List() (keys []string, err error) {
	sec := NewVaultSecret().SetClient(hash.client)
	keys, err = sec.List()
	return keys, err
}
