package trainjector

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing"
)

// TraceInjector inserts a span in the txctx carrier from txctx.Context
func TraceInjector(tracer *opentracing.Tracer, _ string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Next()
		_ = tracer.InjectFromContext(txctx)
	}
}
