package jwt

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	checker  *JWT
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if checker != nil {
			return
		}

		conf := NewConfig(viper.GetViper())
		if len(conf.Certificate) == 0 {
			log.Infof("jwt: no certificate provided")
		}

		var err error
		checker, err = New(conf)
		if err != nil {
			log.WithError(err).Fatalf("jwt: could not create checker")
		}
	})
}

// GlobalChecker returns global Authentication Manager
func GlobalChecker() *JWT {
	return checker
}

// SetGlobalAuth sets global Authentication Manager
func SetGlobalChecker(c *JWT) {
	checker = c
	log.Debug("authentication manager: set")
}
