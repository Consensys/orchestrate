package hashicorp

import (
	"context"
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

// SecretStore wraps a HashiCorp client an manage the unsealing
type SecretStore struct {
	mut        sync.Mutex
	rtl        *RenewTokenLoop
	Client     *Hashicorp
	Config     *Config
	KeyBuilder *multitenancy.KeyBuilder
}

// NewSecretStore construct a new HashiCorp vault given a configfile or nil
func NewSecretStore(config *Config, keyBuilder *multitenancy.KeyBuilder) (*SecretStore, error) {
	hash, err := NewVaultClient(config)
	if err != nil {
		return nil, errors.InternalError("HashiCorp: Could not start vault: %v", err)
	}

	err = hash.SetTokenFromConfig(config)
	if err != nil {
		return nil, errors.InternalError("HashiCorp: Could not start vault: %v", err)
	}

	store := &SecretStore{
		Client:     hash,
		Config:     config,
		KeyBuilder: keyBuilder,
	}

	err = store.ManageToken()
	if err != nil {
		return nil, err
	}

	return store, nil
}

// ManageToken starts a loop that will renew the token automatically
func (store *SecretStore) ManageToken() error {
	secret, err := store.Client.Auth().Token().LookupSelf()
	if err != nil {
		return errors.InternalError("HashiCorp: Initial token lookup failed: %v", err)
	}

	log.Infof("HashiCorp: Token data %q", secret.Data)
	tokenTTL64, err := secret.Data["creation_ttl"].(json.Number).Int64()
	if err != nil {
		return errors.InternalError("HashiCorp: Could not read vault creation ttl: %v", err)
	}

	if int(tokenTTL64) == 0 {
		log.Info("HashiCorp: token does never expire(root token)")
		return nil
	}

	tokenExpireIn64, err := secret.Data["ttl"].(json.Number).Int64()
	if err != nil {
		return errors.InternalError("HashiCorp: Could not read vault ttl: %v", err)
	}
	log.Debugf("HashiCorp: Vault token expires in %d seconds", tokenExpireIn64)

	store.rtl = &RenewTokenLoop{
		TTL:               int(tokenExpireIn64),
		Quit:              make(chan bool, 1),
		Hash:              store,
		RtlTimeRetry:      2,
		RtlMaxNumberRetry: 3,
	}

	err = store.rtl.Refresh()
	if err != nil {
		return errors.InternalError("HashiCorp: Initial token refresh failed: %v", err)
	}

	log.Info("HashiCorp: Initial token refresh succeeded")

	// Start refresh token loop
	store.rtl.Run()

	return nil
}

// Store writes in the vault
func (store *SecretStore) Store(ctx context.Context, key, value string) error {
	storedValue, ok, err := store.Client.Logical.Read(key)
	if err != nil {
		return errors.ConnectionError(err.Error()).ExtendComponent(component)
	}

	if ok {
		if storedValue == value {
			return nil
		}
		return errors.AlreadyExistsError("HashiCorp: A different secret already exists for key: %v", key).ExtendComponent(component)
	}

	err = store.Client.Logical.Write(key, value)
	if err != nil {
		return errors.ConnectionError(err.Error()).ExtendComponent(component)
	}

	log.WithField("key", key).Info("HashiCorp: secret has been stored successfully")
	return nil
}

// Load reads in the vault
func (store *SecretStore) Load(_ context.Context, key string) (value string, ok bool, e error) {
	v, ok, err := store.Client.Logical.Read(key)
	if err != nil {
		e = errors.ConnectionError(err.Error()).ExtendComponent(component)
	} else if ok {
		return v, ok, nil
	}

	return "", false, e
}

// Delete removes a path in the vault
func (store *SecretStore) Delete(ctx context.Context, key string) error {
	err := store.Client.Logical.Delete(key)
	if err != nil {
		return errors.ConnectionError(err.Error()).ExtendComponent(component)
	}
	return nil
}

// List returns the list of all secrets stored in the vault
func (store *SecretStore) List() (keys []string, err error) {
	keys, err = store.Client.Logical.List("")
	if err != nil {
		return []string{}, errors.ConnectionError(err.Error()).ExtendComponent(component)
	}
	return keys, err
}
