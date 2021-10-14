package jwt

import (
	"context"
	"sync"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	checker  *JWT
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		vipr := viper.GetViper()
		isMultiTenancyEnabled := vipr.GetBool(multitenancy.EnabledViperKey)
		if !isMultiTenancyEnabled {
			return
		}

		logger := log.WithContext(ctx)
		if checker != nil {
			return
		}

		conf, err := NewConfig(vipr)
		if err != nil {
			logger.WithError(err).Fatalf("jwt: failed to init")
		}

		if len(conf.Certificates) == 0 {
			logger.Fatalf("jwt: no certificate provided")
		}

		checker, err = New(conf)
		if err != nil {
			logger.WithError(err).Fatalf("jwt: could not create checker")
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
}
