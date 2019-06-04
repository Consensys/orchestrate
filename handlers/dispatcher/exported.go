package dispatcher

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/services/chanregistry"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Dispatcher Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Channel registry
		chanregistry.Init(ctx)

		handler = Dispacher(chanregistry.GlobalChanRegistry())

		log.Infof("dispatcher: handler ready")
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
