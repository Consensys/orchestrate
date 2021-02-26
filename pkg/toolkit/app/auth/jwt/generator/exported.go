package generator

import (
	"context"
	"sync"

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
		var err error
		jwtGenerator, err = New(NewConfig(viper.GetViper()))
		if err != nil {
			log.WithContext(ctx).WithError(err).Fatalf("jwt-generator: could not create jwtGenerator")
		}
	})
}

func GlobalJWTGenerator() *JWTGenerator {
	return jwtGenerator
}

func SetJWTGenerator(generator *JWTGenerator) {
	jwtGenerator = generator
}
