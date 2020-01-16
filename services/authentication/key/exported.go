package key

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	auth     *Auth
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if auth != nil {
			return
		}

		auth = NewAuth(viper.GetString(APIKeyViperKey))
	})
}

// GlobalAuth returns global Authentication Manager
func GlobalAuth() *Auth {
	return auth
}

// SetGlobalAuth sets global Authentication Manager
func SetGlobalAuth(a *Auth) {
	auth = a
	log.Debug("auth: set")
}
