package multitenancy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
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
		if authutils.HasAllPrivileges(req.Context()) {
			h.ServeHTTP(rw, req)
			return
		}

		tenantID := multitenancy.TenantIDFromContext(req.Context())
		if m.tenantID != tenantID {
			log.FromContext(req.Context()).
				WithField("expected", m.tenantID).
				WithField("received", tenantID).
				Debugf("invalid tenant id")
			m.serveNotFound(rw)
		} else {
			h.ServeHTTP(rw, req)
		}
	})
}

func (m *MultiTenant) serveNotFound(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusNotFound)
	_, _ = rw.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusNotFound, http.StatusText(http.StatusNotFound))))
}
