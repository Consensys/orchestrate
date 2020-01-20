package generator

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	jwtGenerator *JWTGenerator
	initOnce     = &sync.Once{}
)

// Init initializes key Builder with EnabledViperKey
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if jwtGenerator != nil {
			return
		}
		var err error
		jwtGenerator, err = New(
			viper.GetString(jwt.ClaimsNamespaceViperKey),
			viper.GetString(PrivateKeyViperKey),
		)
		if err != nil {
			log.WithError(err).Fatalf("jwt-generator: could not create jwtGenerator")
		}
	})
}

func GlobalJWTGenerator() *JWTGenerator {
	return jwtGenerator
}

func SetJWTGenerator(generator *JWTGenerator) {
	jwtGenerator = generator
}
