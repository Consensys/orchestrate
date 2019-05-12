package store

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	grpcstore "gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Sender Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Context store
		grpcstore.Init(ctx)

		// Create Handler
		handler = EnvelopeLoader(grpcstore.GlobalEnvelopeStore())

		log.Infof("envelope-store: handler ready")
	})
}

// SetGlobalHandler sets global Gas Pricer Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Gas Pricer handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
