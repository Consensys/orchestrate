package token

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"

	log "github.com/sirupsen/logrus"
)

var (
	auth     authentication.Manager
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if auth != nil {
			return
		}

		auth = New()
	})
}

// GlobalAuth returns global Authentication Manager
func GlobalAuth() authentication.Manager {
	return auth
}

// SetGlobalAuth sets global Authentication Manager
func SetGlobalAuth(authManager authentication.Manager) {
	auth = authManager
	log.Debug("authentication manager: set")
}
