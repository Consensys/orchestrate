package infra

import (
	"context"
	"sync"
	"sync/atomic"

	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra"
)

var i *Infra

func init() {
	i = New()
}

// Infra elements of the application
type Infra struct {
	ctx context.Context

	tracer opentracing.Tracer
	store  infra.TraceStore

	initOnce, closeOnce *sync.Once
	cancel              func()

	ready *atomic.Value
}

// New creates a new infrastructure
func New() *Infra {
	return &Infra{
		initOnce:  &sync.Once{},
		closeOnce: &sync.Once{},
		ready:     &atomic.Value{},
	}
}

// Init intilialize infrastructure
func Init() {
	i.initOnce.Do(func() {
		i.ctx, i.cancel = context.WithCancel(context.Background())
		InitTracing(i)
		InitStore(i)
		i.ready.Store(true)
	})
}

// Tracer returns tracer
func Tracer() opentracing.Tracer {
	if !Ready() {
		panic("Infra is not ready. Please call Init() first")
	}
	return i.tracer
}

// Store returns trace store
func Store() infra.TraceStore {
	if !Ready() {
		panic("Infra is not ready. Please call Init() first")
	}
	return i.store
}

// Ready indicate if infra is ready
func Ready() bool {
	select {
	case <-i.ctx.Done():
		return false
	default:
		return i.ready.Load().(bool)
	}
}

// Close infra
func Close() {
	i.closeOnce.Do(i.cancel)
}
