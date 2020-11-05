package hashicorp

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// RenewTokenLoop handle the token renewal of the application
type RenewTokenLoop struct {
	TTL  int
	Quit chan bool
	Hash *SecretStore

	RtlTimeRetry      int // RtlTimeRetry: Time between each retry of token renewal
	RtlMaxNumberRetry int // RtlMaxNumberRetry: Max number of retry for token renewal
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
			loop.Hash.Client.Client.SetToken(
				newTokenSecret.Auth.ClientToken,
			)
			loop.Hash.mut.Unlock()
			log.Info("SecretStore: Vault token was refreshed successfully")
			return nil
		}

		retry++
		if retry < loop.RtlMaxNumberRetry {
			// Max number number of retry reached: graceful shutdown
			log.Error("SecretStore: Graceful shutdown of the vault, the token could not be renewed")
			return errors.InternalError("SecretStore: Token refresh failed (%v)", err).SetComponent(component)
		}

		time.Sleep(time.Duration(loop.RtlTimeRetry) * time.Second)
	}
}

// Run contains the token regeneration routine
func (loop *RenewTokenLoop) Run() {
	go func() {
		timeToWait := time.Duration(
			int(float64(loop.TTL)*0.75), // We wait 75% of the TTL to refresh
		) * time.Second

		// Max token refresh loop of 1h
		if timeToWait > time.Hour {
			log.Infof("HashiCorp: forcing token refresh to maximum one hour")
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
					loop.Quit <- true
				}

			// TODO: Be able to graceful shutdown every other services in the infra
			case <-loop.Quit:
				// The token parameter is ignored
				_ = loop.Hash.Client.Auth().Token().RevokeSelf("this parameter is ignored")
				// Erase the local value of the token
				loop.Hash.mut.Lock()
				loop.Hash.Client.Client.SetToken("")
				loop.Hash.mut.Unlock()
				// Wait 5 seconds for the ongoing requests to return
				time.Sleep(time.Duration(5) * time.Second)
				// Crash the tx-signer to force restart
				log.Fatal("SecretStore: Graceful shutdown of the vault, the token has been revoked")
			}
		}
	}()
}
