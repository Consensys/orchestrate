package accesslog

import (
	"context"
	"fmt"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares/accesslog"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

type Builder struct {
	handlers map[string]*accesslog.Handler
}

func NewBuilder() *Builder {
	b := &Builder{
		handlers: make(map[string]*accesslog.Handler),
	}

	return b
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	h, ok := b.handlers[name]
	if !ok {
		cfg, ok := configuration.(*dynamic.AccessLog)
		if !ok {
			return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
		}

		log.FromContext(ctx).
			WithField("middleware", name).
			WithField("type", fmt.Sprintf("%T", configuration)).
			Debugf("building middleware")

		h, err = accesslog.NewHandler(cfg.ToTraefikType())
		if err != nil {
			return nil, nil, err
		}

		b.handlers[name] = h
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			h.ServeHTTP(rw, req, next)
		})
	}, nil, nil
}
