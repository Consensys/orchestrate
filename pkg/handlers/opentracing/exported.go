package opentracing

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/tracing/opentracing/jaeger"
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

		operationName := "GenericService"
		if ctx.Value("service-name") != nil {
			operationName = ctx.Value("service-name").(string)
		}

		// Create Handler
		handler = TxSpanFromBroker(opentracing.GlobalTracer(), operationName)

		log.Infof("logger: open-tracing handler ready")
	})
}

// SetGlobalHandler sets global Opentracing Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Opentracing handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
