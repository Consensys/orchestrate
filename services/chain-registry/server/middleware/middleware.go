package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/containous/alice"
	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/containous/traefik/v2/pkg/server/middleware"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/server/utils"
)

// Envelope the middleware builder
type Builder struct {
	builder *middleware.Builder

	orchestrateMiddlewares map[string]func(routerName string) alice.Constructor
}

type serviceBuilder interface {
	BuildHTTP(ctx context.Context, serviceName string, responseModifier func(*http.Response) error) (http.Handler, error)
}

// NewBuilder creates a new Envelope
func NewBuilder(
	configs map[string]*runtime.MiddlewareInfo,
	serviceBuilder serviceBuilder,
	orchestrateMiddlewares map[string]func(routerName string) alice.Constructor,
) *Builder {
	return &Builder{
		builder:                middleware.NewBuilder(configs, serviceBuilder),
		orchestrateMiddlewares: orchestrateMiddlewares,
	}
}

// BuildChain creates a middleware chain
func (b *Builder) BuildChain(ctx context.Context, middlewares []string, routerName string) *alice.Chain {
	// Compute Orchestrate custom Middleware chain
	chain := alice.New()
	var traefikMiddlewares []string
	for _, middleware := range middlewares {
		parts := strings.Split(utils.GetQualifiedName(ctx, middleware), "@")
		if constructorBuilder, ok := b.orchestrateMiddlewares[parts[0]]; ok {
			chain = chain.Append(constructorBuilder(routerName))
		} else {
			traefikMiddlewares = append(traefikMiddlewares, middleware)
		}
	}

	// Build Traeffik Middleware chain
	traefikChain := b.builder.BuildChain(ctx, traefikMiddlewares)

	// Extend Traeffik chain with Orchestrate chain
	chain = traefikChain.Extend(chain)

	return &chain
}
