package multitenancy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(_ context.Context, _ string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.MultiTenancy)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	m := New(cfg.Tenant)

	return m.Handler, nil, nil
}

type MultiTenant struct {
	tenantID string
}

func New(tenantID string) *MultiTenant {
	return &MultiTenant{
		tenantID: tenantID,
	}
}

func (m *MultiTenant) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		tenants := multitenancy.AllowedTenantsFromContext(req.Context())
		if !multitenancy.IsAllowed(m.tenantID, tenants) {
			log.FromContext(req.Context()).
				WithField("expected", m.tenantID).
				WithField("received", tenants).
				Debugf("invalid tenant id")
			m.serveNotFound(rw)
			return
		}

		h.ServeHTTP(rw, req)
	})
}

func (m *MultiTenant) serveNotFound(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusNotFound)
	_, _ = rw.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusNotFound, http.StatusText(http.StatusNotFound))))
}
