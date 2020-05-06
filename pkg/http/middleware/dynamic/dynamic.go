package dynamic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/middlewares/addprefix"
	"github.com/containous/traefik/v2/pkg/middlewares/buffering"
	"github.com/containous/traefik/v2/pkg/middlewares/circuitbreaker"
	"github.com/containous/traefik/v2/pkg/middlewares/compress"
	"github.com/containous/traefik/v2/pkg/middlewares/inflightreq"
	"github.com/containous/traefik/v2/pkg/middlewares/ipwhitelist"
	"github.com/containous/traefik/v2/pkg/middlewares/passtlsclientcert"
	"github.com/containous/traefik/v2/pkg/middlewares/ratelimiter"
	"github.com/containous/traefik/v2/pkg/middlewares/redirect"
	"github.com/containous/traefik/v2/pkg/middlewares/replacepath"
	"github.com/containous/traefik/v2/pkg/middlewares/replacepathregex"
	"github.com/containous/traefik/v2/pkg/middlewares/retry"
	"github.com/containous/traefik/v2/pkg/middlewares/stripprefix"
	"github.com/containous/traefik/v2/pkg/middlewares/stripprefixregex"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/cors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/headers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/httptrace"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/loadbalancer"
	reflectmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/reflect"
)

var errBadConf = errors.New("cannot create middleware: multi-types middleware not supported, consider declaring two different pieces of middleware instead")

type Builder struct {
	reflect *reflectmid.Builder
	traefik *TraefikBuilder
}

func NewBuilder() *Builder {
	b := &Builder{
		reflect: reflectmid.NewBuilder(),
		traefik: &TraefikBuilder{},
	}

	b.AddBuilder(reflect.TypeOf(&dynamic.Cors{}), &cors.Builder{})
	b.AddBuilder(reflect.TypeOf(&dynamic.Headers{}), &headers.Builder{})
	b.AddBuilder(reflect.TypeOf(&dynamic.LoadBalancer{}), &loadbalancer.Builder{})
	b.AddBuilder(reflect.TypeOf(&dynamic.HTTPTrace{}), &httptrace.Builder{})

	return b
}

func (b *Builder) AddBuilder(typ reflect.Type, builder middleware.Builder) {
	b.reflect.AddBuilder(typ, builder)
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.Middleware)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	field, err := cfg.Field()
	if err != nil {
		return nil, nil, err
	}

	if trefikCfg, ok := field.(*traefikdynamic.Middleware); ok {
		return b.traefik.Build(ctx, name, trefikCfg)
	}

	return b.reflect.Build(ctx, name, field)
}

type TraefikBuilder struct{}

func (b *TraefikBuilder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*traefikdynamic.Middleware)
	if !ok {
		return nil, nil, fmt.Errorf("expected configuration type %T but got %T", cfg, configuration)
	}

	// AddPrefix
	if cfg.AddPrefix != nil {
		if cfg.AddPrefix.Prefix == "" {
			return nil, nil, fmt.Errorf("prefix cannot be empty")
		}

		mid = func(next http.Handler) http.Handler {
			h, _ := addprefix.New(ctx, next, *cfg.AddPrefix, name)
			return h
		}
	}

	// Buffering
	if cfg.Buffering != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := buffering.New(ctx, next, *cfg.Buffering, name)
			return h
		}
	}

	// CircuitBreaker
	if cfg.CircuitBreaker != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := circuitbreaker.New(ctx, next, *cfg.CircuitBreaker, name)
			return h
		}
	}

	// Compress
	if cfg.Compress != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := compress.New(ctx, next, *cfg.Compress, name)
			return h
		}
	}

	// ContentType
	if cfg.ContentType != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !cfg.ContentType.AutoDetect {
					rw.Header()["Content-Type"] = nil
				}
				next.ServeHTTP(rw, req)
			})
		}
	}

	// IPWhiteList
	if cfg.IPWhiteList != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := ipwhitelist.New(ctx, next, *cfg.IPWhiteList, name)
			return h
		}
	}

	// InFlightReq
	if cfg.InFlightReq != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := inflightreq.New(ctx, next, *cfg.InFlightReq, name)
			return h
		}
	}

	// PassTLSClientCert
	if cfg.PassTLSClientCert != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := passtlsclientcert.New(ctx, next, *cfg.PassTLSClientCert, name)
			return h
		}
	}

	// RateLimit
	if cfg.RateLimit != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := ratelimiter.New(ctx, next, *cfg.RateLimit, name)
			return h
		}
	}

	// RedirectRegex
	if cfg.RedirectRegex != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := redirect.NewRedirectRegex(ctx, next, *cfg.RedirectRegex, name)
			return h
		}
	}

	// RedirectScheme
	if cfg.RedirectScheme != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := redirect.NewRedirectScheme(ctx, next, *cfg.RedirectScheme, name)
			return h
		}
	}

	// ReplacePath
	if cfg.ReplacePath != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := replacepath.New(ctx, next, *cfg.ReplacePath, name)
			return h
		}
	}

	// ReplacePathRegex
	if cfg.ReplacePathRegex != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := replacepathregex.New(ctx, next, *cfg.ReplacePathRegex, name)
			return h
		}
	}

	// Retry
	if cfg.Retry != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			// FIXME missing metrics / accessLog
			h, _ := retry.New(ctx, next, *cfg.Retry, retry.Listeners{}, name)
			return h
		}
	}

	// StripPrefix
	if cfg.StripPrefix != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := stripprefix.New(ctx, next, *cfg.StripPrefix, name)
			return h
		}
	}

	// StripPrefixRegex
	if cfg.StripPrefixRegex != nil {
		if mid != nil {
			return nil, nil, errBadConf
		}
		mid = func(next http.Handler) http.Handler {
			h, _ := stripprefixregex.New(ctx, next, *cfg.StripPrefixRegex, name)
			return h
		}
	}

	return mid, nil, nil
}
