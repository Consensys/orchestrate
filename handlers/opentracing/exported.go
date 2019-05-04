package opentracing

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Crafter Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Controlled Faucet
		jaeger.Init(ctx)

		// Create Handler
		handler = TxSpanFromBroker(opentracing.GlobalTracer(), ctx.Value("service-name").(string))

		log.Infof("logger: handler ready")
	})
}

// SetGlobalHandler sets global Faucet Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Faucet handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
