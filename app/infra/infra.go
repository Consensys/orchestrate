package infra

import (
	"context"
	"sync"
	"sync/atomic"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/store"
	pgStore "gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/store/pg"
)

var i *Infra

func init() {
	i = New()
}

// Infra elements of the application
type Infra struct {
	ctx context.Context

	tracer opentracing.Tracer

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
		pgStore.Init()
		i.ready.Store(true)
	})
}

// Tracerr returns tracer
func Tracer() opentracing.Tracer {
	if !Ready() {
		log.Fatal("Infra is not ready. Please call Init() first")
	}
	return i.tracer
}

// Store returns envelope store
func Store() store.EnvelopeStore {
	if !Ready() {
		log.Fatal("Infra is not ready. Please call Init() first")
	}
	return pgStore.GlobalEnvelopeStore()
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
