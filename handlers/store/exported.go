package store

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	storeclient "gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/client"
)

var (
	component = "handler.listener.store"
	handler   engine.HandlerFunc
	initOnce  = &sync.Once{}
)

// Init initialize Sender Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Envelope store client
		storeclient.Init(ctx)

		// Create Handler
		handler = EnvelopeLoader(storeclient.GlobalEnvelopeStoreClient())

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
