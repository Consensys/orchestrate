package generator

import (
	"context"
	"sync"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	jwtGenerator *JWTGenerator
	initOnce     = &sync.Once{}
)

// Init initializes key Envelope with EnabledViperKey
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if jwtGenerator != nil {
			return
		}

		vipr := viper.GetViper()
		if !vipr.GetBool(multitenancy.EnabledViperKey) {
			return
		}

		cfg, err := NewConfig(vipr)
		if err != nil {
			log.WithContext(ctx).WithError(err).Fatalf("jwt-generator: could not load config")
		}

		jwtGenerator, err = New(cfg)
		if err != nil {
			log.WithContext(ctx).WithError(err).Fatalf("jwt-generator: could not initialized")
		}

		log.WithContext(ctx).Info("jwt-generator: initialized successfully")
	})
}

func GlobalJWTGenerator() *JWTGenerator {
	return jwtGenerator
}

func SetJWTGenerator(generator *JWTGenerator) {
	jwtGenerator = generator
}
