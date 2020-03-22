package opentracing

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize opentracing Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Controlled Faucet
		opentracing.Init(ctx)

		operationName := "GenericService"
		if ctx.Value("service-name") != nil {
			operationName = ctx.Value("service-name").(string)
		}

		// Create Handler
		handler = TxSpanFromBroker(opentracing.GetGlobalTracer(), operationName)

		log.Infof("logger: open-tracing span-creator handler ready")
	})
}

// SetGlobalHandler sets global OpenTracing Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global OpenTracing handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
