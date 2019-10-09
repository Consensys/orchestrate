package hashicorp

import (
	"encoding/json"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

// SecretStore wraps a hashicorps client an manage the unsealing
type SecretStore struct {
	mut    sync.Mutex
	rtl    *RenewTokenLoop
	Client *Hashicorp
	Config *Config
}

// NewSecretStore construct a new hashicorps vault given a configfile or nil
func NewSecretStore(config *Config) (*SecretStore, error) {

	client, err := NewVaultClient(config)
	if err != nil {
		log.Fatalf("Could not start vault: %v", err)
	}

	err = client.SetTokenFromConfig(config)
	if err != nil {
		log.Fatalf("Could not start vault: %v", err)
	}

	hash := &SecretStore{
		Client: client,
		Config: config,
	}

	store.ManageToken()
	return hash, nil
}

// ManageToken starts a loop that will renew the token automatically
func (store *SecretStore) ManageToken() {
	secret, err := store.Client.Auth().Token().LookupSelf()
	if err != nil {
		log.Fatalf("Initial token lookup failed: %v", err)
	}

	vaultTTL64, err := secret.Data["ttl"].(json.Number).Int64()
	if err != nil {
		log.Fatalf("Could not read vault ttl: %v", err)
	}

	vaultTokenTTL := int(vaultTTL64)
	if vaultTokenTTL < 1 {
		// case where the tokenTTL is infinite
		return
	}

	log.Debugf("Vault TTL: %v", vaultTokenTTL)
	log.Debugf("64: %v", vaultTTL64)

	timeToWait := time.Duration(
		int(float64(
			vaultTokenTTL,
		)*0.75), // We wait 75% of the TTL to refresh
	) * time.Second

	ticker := time.NewTicker(timeToWait)
	log.Debugf("time to wait: %v", timeToWait)

	store.rtl = &RenewTokenLoop{
		TTL:    vaultTokenTTL,
		ticker: ticker,
		Quit:   make(chan bool, 1),
		Hash:   store,

		RtlTimeRetry:      2,
		RtlMaxNumberRetry: 3,
	}

	err = store.rtl.Refresh()
	if err != nil {
		log.Fatalf("Initial token refresh failed: %v", err)
	}
}

// Store writes in the vault
func (store *SecretStore) Store(key, value string) error {
	_, ok, err := store.Client.Logical.Read(key)
	if err != nil {
		return errors.ConnectionError(err.Error()).ExtendComponent(component)
	}

	if ok {
		return errors.AlreadyExistsError("A secret already exists for key: %v", key).ExtendComponent(component)
	}

	err = store.Client.Logical.Write(key, value)
	if err != nil {
		return errors.ConnectionError(err.Error()).ExtendComponent(component)
	}
	return nil
}

// Load reads in the vault
func (store *SecretStore) Load(key string) (value string, ok bool, err error) {
	value, ok, err = store.Client.Logical.Read(key)
	if err != nil {
		return "", false, errors.ConnectionError(err.Error()).ExtendComponent(component)
	}
	return value, ok, err
}

// Delete removes a path in the vault
func (store *SecretStore) Delete(key string) error {
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
