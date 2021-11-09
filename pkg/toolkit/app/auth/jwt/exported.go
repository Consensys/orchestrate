package jwt

import (
	"context"
	"sync"

	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/jwt/jose"

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
		cfg := jose.NewConfig(vipr)
		if cfg == nil {
			logger.Fatalf("jwt: no identity provider issuer url provided")
		}

		validator, err := jose.NewValidator(cfg)
		if err != nil {
			logger.WithError(err).Fatalf("jwt: could not create jwt validator")
		}

		checker = New(validator)
	})
}

// GlobalChecker returns global Authentication Manager
func GlobalChecker() *JWT {
	return checker
}

// SetGlobalChecker sets global Authentication Manager
func SetGlobalChecker(c *JWT) {
	checker = c
}
