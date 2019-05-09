package hashicorp

import (
	"github.com/hashicorp/vault/api"
)

// Secret is a generic interface that is implemented by all secret engines
type Secret interface {
	SaveNew() (err error)
	Update() error
	Delete() error
	SetClient(client *api.Client)
	List(subPath string) ([]string, error)
	GetValue() (value string, ok bool, err error)
}

// NewSecret construct a new secret with an implementation depending on which
// secret engine s being used.
func NewSecret(key, value string) Secret {
	kvVersion := GetKVVersion()
	if kvVersion == "v1" {
		return NewSecretV1(key, value)
	}
	return NewSecretV2(key, value)
}