package infra

import (
	"context"
	"sync"
	"sync/atomic"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry"
)

var i *Infra

func init() {
	i = New()
}

// Infra elements of the application
type Infra struct {
	ctx context.Context

	tracer opentracing.Tracer

	initOnce *sync.Once
	cancel   func()

	ready *atomic.Value
}

// New creates a new infrastructure
func New() *Infra {
	infra := &Infra{
		ready:    &atomic.Value{},
		initOnce: &sync.Once{},
	}
	i.ready.Store(false)
	return infra
}

// Init intilialize infrastructure
func Init() {
	i.initOnce.Do(func() {
		i.ctx, i.cancel = context.WithCancel(context.Background())
		InitTracing(i)
		registry.Init(i.ctx)
		i.ready.Store(true)
	})
}

// Tracer returns tracer
func Tracer() opentracing.Tracer {
	if !Ready() {
		log.Fatal("Infra is not ready. Please call Init() first")
	}
	return i.tracer
}

// Registry returns contract registry
func Registry() registry.Registry {
	if !Ready() {
		log.Fatal("Infra is not ready. Please call Init() first")
	}
	return registry.GlobalRegistry()
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
	i.cancel()
}
