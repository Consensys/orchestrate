package multitenancy

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

const component = "handler.multitenancy"

var (
	handler     engine.HandlerFunc
	authHandler engine.HandlerFunc
	initOnce    = &sync.Once{}
)

// Init initialize Multi Tenancy Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil && authHandler != nil {
			return
		}

		// Initialize Authentication Manager
		jwt.Init(ctx)
		log.Infof("multitenancy enable: %v", viper.GetBool(multitenancy.EnabledViperKey))

		// Create Handler
		handler = ExtractTenant(viper.GetBool(multitenancy.EnabledViperKey), nil)
		authHandler = ExtractTenant(viper.GetBool(multitenancy.EnabledViperKey), jwt.GlobalChecker())

		log.Infof("authentication multi-tenancy: handler ready")
	})
}

// SetGlobalHandler sets global Gas Estimator Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

func GlobalHandler() engine.HandlerFunc {
	return handler
}
