package multitenancy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/traefik/traefik/v2/pkg/log"
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

	m := New(cfg.Tenant, cfg.OwnerID)

	return m.Handler, nil, nil
}

type MultiTenant struct {
	tenantID string
	ownerID  string
}

func New(tenantID, ownerID string) *MultiTenant {
	return &MultiTenant{
		tenantID: tenantID,
		ownerID:  ownerID,
	}
}

func (m *MultiTenant) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		userInfo := multitenancy.UserInfoValue(req.Context())
		if !userInfo.HasTenantAccess(m.tenantID) || !userInfo.HasUsernameAccess(m.ownerID) {
			log.FromContext(req.Context()).
				WithField("expected", m.tenantID).
				WithField("received", userInfo.AllowedTenants).
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
