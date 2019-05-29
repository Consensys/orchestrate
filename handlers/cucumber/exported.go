package cucumber

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/cucumber/chanregistry"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Gas Pricer Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Channel registry
		chanregistry.Init(ctx)

		handler = Cucumber(chanregistry.GlobalChanRegistry())

		log.Infof("cucumber: handler ready")
	})
}

// SetGlobalHandler sets global Cucumber Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Cucumber handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
