package hashicorp

import (
	"sync"
	"time"

	"github.com/hashicorp/vault/api"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

// renewTokenLoop handle the token renewal of the application
type renewTokenLoop struct {
	ttl           int
	quit          chan bool
	client        *api.Client
	mut           *sync.Mutex
	retryInterval int
	maxRetries    int
}

func newRenewTokenLoop(tokenExpireIn64 int64, client *api.Client) *renewTokenLoop {
	return &renewTokenLoop{
		ttl:           int(tokenExpireIn64),
		quit:          make(chan bool, 1),
		client:        client,
		retryInterval: 2,
		maxRetries:    3,
		mut:           &sync.Mutex{},
	}
}

// Refresh the token
func (loop *renewTokenLoop) Refresh() error {
	retry := 0
	for {
		// Regularly try renewing the token
		newTokenSecret, err := loop.client.Auth().Token().RenewSelf(0)

		if err == nil {
			loop.mut.Lock()
			loop.client.SetToken(newTokenSecret.Auth.ClientToken)
			loop.mut.Unlock()
			log.Info("Hashicorp Vault token was refreshed successfully")
			return nil
		}

		retry++
		if retry < loop.maxRetries {
			errMessage := "reached max number of retries to renew vault token"
			log.WithField("retries", retry).Error(errMessage)
			return errors.InternalError(errMessage)
		}

		time.Sleep(time.Duration(loop.retryInterval) * time.Second)
	}
}

// Run contains the token regeneration routine
func (loop *renewTokenLoop) Run() {
	go func() {
		timeToWait := time.Duration(
			int(float64(loop.ttl)*0.75), // We wait 75% of the TTL to refresh
		) * time.Second

		// Max token refresh loop of 1h
		if timeToWait > time.Hour {
			log.Info("HashiCorp: forcing token refresh to maximum one hour")
			timeToWait = time.Hour
		}

		ticker := time.NewTicker(timeToWait)
		defer ticker.Stop()

		log.Infof("HashiCorp: token refresh loop started (every %d seconds)", timeToWait/time.Second)
		for {
			select {
			case <-ticker.C:
				err := loop.Refresh()
				if err != nil {
					loop.quit <- true
				}

			// TODO: Be able to graceful shutdown every other services in the infra
			case <-loop.quit:
				// The token parameter is ignored
				_ = loop.client.Auth().Token().RevokeSelf("this parameter is ignored")
				// Erase the local value of the token
				loop.mut.Lock()
				loop.client.SetToken("")
				loop.mut.Unlock()
				// Wait 5 seconds for the ongoing requests to return
				time.Sleep(time.Duration(5) * time.Second)
				// Crash the tx-signer to force restart
				log.Fatal("gracefully shutting down the vault client, the token has been revoked")
			}
		}
	}()
}
