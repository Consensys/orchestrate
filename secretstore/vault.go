package secretstore

import import (
	"github.com/hashicorp/vault/api"
	"strings"
	"fmt"
)

type Hashicorps struct {
	client *api.Client
	creds *credentials
	retrieveSecretOnce *sync.Once
}

func NewHashicorps(config api.Config) (*Hashicorps) {

	if config == nil {
		config = api.DefaultConfig()
	}

	client := api.NewClient(config)
	creds := &credentials{}
	return &Hashicorps{
		client: client
		creds: creds
	}
}

func (hash *Hashicorps) Init(credsStore *AWS, tokenName string) (*Hashicorps, error) {

	err = hash.creds.FetchFromAWS(credsStore, tokenName)
	if err != nil {
		return err
	}

	hash.creds.AttachTo(s.client)

	err = hash.creds.Unseal(s.client)
	if err != nil {
		return err
	}
}

func (hash *Hashicorps) Store(key, value string) (err error) {
	sec := NewVaultSecret().SetKey(key).setValue(value).SetClient(hash.client)
	err := sec.Update()
	return err
}

func (hash *Hashicorps) Load(key string) (value string, err error) {
	sec := NewVaultSecret().SetKey(key).SetClient(hash.client)
	res, err := sec.GetValue()
	if err != nil {
		return nil, err
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
