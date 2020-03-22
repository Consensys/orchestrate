package multitenancy

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

const component = "handler.multitenancy"

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Multi Tenancy Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Key Envelope
		multitenancy.Init(ctx)

		// Initialize Authentication Manager
		jwt.Init(ctx)

		log.Infof("multitenancy enable: %v", viper.GetBool(multitenancy.EnabledViperKey))

		// Create Handler
		handler = ExtractTenant(jwt.GlobalChecker())

		log.Infof("authentication multi-tenancy: handler ready")
	})
}

// SetGlobalHandler sets global Gas Estimator Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Gas Estimator handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
