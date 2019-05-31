package hashicorp

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// HashiCorp wraps a hashicorps client an manage the unsealing
type HashiCorp struct {
	mut    sync.Mutex
	rtl    *RenewTokenLoop
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

	_ = WithVaultToken(client) // If the token was not found. The error is ignored

	hash := &HashiCorp{
		Client: client,
	}

	hash.manageToken()
	return hash, nil
}

func (hash *HashiCorp) manageToken() {

	secret, err := hash.Client.Auth().Token().LookupSelf()
	if err != nil {
		log.Fatalf("Initial token lookup failed : %v", err)
	}

	vaultTTL64, err := secret.Data["ttl"].(json.Number).Int64()
	if err != nil {
		log.Fatalf("Could not read vault ttl : %v", err)
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

	hash.rtl = &RenewTokenLoop{
		TTL:    vaultTokenTTL,
		ticker: ticker,
		Quit:   make(chan bool, 1),
		Hash:   hash,

		RtlTimeRetry:      2,
		RtlMaxNumberRetry: 3,
	}

	err = hash.rtl.Refresh()
	if err != nil {
		log.Fatalf("Initial token refresh failed : %v", err)
	}
}

// Store writes in the vault
func (hash *HashiCorp) Store(key, value string) (err error) {
	sec := NewSecret(key, value)
	sec.SetClient(hash.Client)

	hash.mut.Lock()
	defer hash.mut.Unlock()
	return sec.Update()
}

// Load reads in the vault
func (hash *HashiCorp) Load(key string) (value string, ok bool, err error) {
	sec := NewSecret(key, "")
	sec.SetClient(hash.Client)

	hash.mut.Lock()
	defer hash.mut.Unlock()
	return sec.GetValue()
}

// Delete removes a path in the vault
func (hash *HashiCorp) Delete(key string) (err error) {
	sec := NewSecret(key, "")
	sec.SetClient(hash.Client)

	hash.mut.Lock()
	defer hash.mut.Unlock()
	return sec.Delete()
}

// List returns the list of all secrets stored in the vault
func (hash *HashiCorp) List() (keys []string, err error) {
	sec := NewSecret("", "")
	sec.SetClient(hash.Client)

	hash.mut.Lock()
	defer hash.mut.Unlock()
	return sec.List("")
}

// RenewTokenLoop handle the token renewal of the application
type RenewTokenLoop struct {
	TTL    int
	ticker *time.Ticker
	Quit   chan bool
	Hash   *HashiCorp

	RtlTimeRetry      int // RtlTimeRetry : Time between each retry of token renewal
	RtlMaxNumberRetry int // RtlMaxNumberRetry : Max number of retry for token renewal
}

// Refresh the token
func (loop *RenewTokenLoop) Refresh() error {
	retry := 0
	for {
		// Regularly try renewing the token
		newTokenSecret, err := loop.
			Hash.Client.Auth().Token().RenewSelf(0)

		if err == nil {
			loop.Hash.mut.Lock()
			loop.Hash.Client.SetToken(
				newTokenSecret.Auth.ClientToken,
			)
			loop.Hash.mut.Unlock()
			log.Debugf("Successfully refreshed token, TokenTTL is %v", loop.TTL)
			return nil
		}

		retry++
		if retry < loop.RtlMaxNumberRetry {
			// Max number number of retry reached : graceful shutdown
			log.Error("Graceful shutdown of the vault, the token could not be renewed")
			return fmt.Errorf("token refresh failed, we got over the max_retry : %v ", err.Error())
		}

		time.Sleep(time.Duration(loop.RtlTimeRetry) * time.Second)
	}
}

// Run contains the token regeneration routine
func (loop *RenewTokenLoop) Run() {

	for {
		select {
		case <-loop.ticker.C:
			err := loop.Refresh()
			if err != nil {
				loop.Quit <- true
			}

		// TODO: Be able to graceful shutdown every other services in the infra
		case <-loop.Quit:
			// The token parameter is ignored
			_ = loop.
				Hash.Client.Auth().Token().RevokeSelf("this parameter is ignored")
			// Erase the local value of the token
			loop.Hash.Client.SetToken("")
			// Wait 5 seconds for the ongoing requests to return
			time.Sleep(time.Duration(5) * time.Second)
			// Crash the app
			log.Fatal("Graceful shutdown of the vault, the token has been revoked")
		}
	}

}
