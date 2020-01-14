package generator

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

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
		jwtGenerator = New(viper.GetBool(multitenancy.EnabledViperKey),
			viper.GetString(authentication.TenantNamespaceViperKey),
			viper.GetString(AuthServicePrivateKeyViperKey))
	})
}

func GlobalJWTGenerator() *JWTGenerator {
	return jwtGenerator
}

func SetJWTGenerator(jwt *JWTGenerator) {
	jwtGenerator = jwt
}
