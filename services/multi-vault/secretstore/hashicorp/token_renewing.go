package hashicorp

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// RenewTokenLoop handle the token renewal of the application
type RenewTokenLoop struct {
	TTL    int
	ticker *time.Ticker
	Quit   chan bool
	Hash   *SecretStore

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
			log.Debugf("Successfully refreshed token, TokenTTL is %v", loop.TTL)
			return nil
		}

		retry++
		if retry < loop.RtlMaxNumberRetry {
			// Max number number of retry reached: graceful shutdown
			log.Error("Graceful shutdown of the vault, the token could not be renewed")
			return errors.InternalError("token refresh failed (%v)", err).SetComponent(component)
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
			loop.Hash.Client.Client.SetToken("")
			// Wait 5 seconds for the ongoing requests to return
			time.Sleep(time.Duration(5) * time.Second)
			// Crash the app
			log.Fatal("Graceful shutdown of the vault, the token has been revoked")
		}
	}

}
