package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth"
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/containous/traefik/v2/pkg/log"
)

type Builder struct {
	jwt, key     auth.Checker
	multitenancy bool
}

func NewBuilder(jwt, key auth.Checker, multitenancyEnabled bool) *Builder {
	return &Builder{
		jwt:          jwt,
		key:          key,
		multitenancy: multitenancyEnabled,
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.Auth)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	m := New(b.jwt, b.key, b.multitenancy)
	return m.Handler, nil, nil
}

type Auth struct {
	jwt, key     auth.Checker
	multitenancy bool
}

func New(jwt, key auth.Checker, multitenancyEnabled bool) *Auth {
	return &Auth{
		jwt:          jwt,
		key:          key,
		multitenancy: multitenancyEnabled,
	}
}

func (a *Auth) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !a.multitenancy {
			a.serveNext(rw, req, h)
			return
		}

		// Extract Authorization credentials from HTTP headers
		authCtx := authutils.WithAuthorization(
			req.Context(),
			req.Header.Get(authutils.AuthorizationHeader),
		)

		// Extract API Key credentials from HTTP headers
		authCtx = authutils.WithAPIKey(
			authCtx,
			req.Header.Get(authutils.APIKeyHeader),
		)

		// Extract TenantID from HTTP headers
		authCtx = multitenancy.WithTenantID(
			authCtx,
			req.Header.Get(multitenancy.TenantIDHeader),
		)

		// Perform API Key Authentication
		checkedCtx, err := a.key.Check(authCtx)
		if err != nil {
			log.FromContext(authCtx).WithError(err).Errorf("unauthorized request")
			a.writeUnauthorized(rw, err)
			return
		}
		if checkedCtx != nil {
			// Bypass JWT authentication
			log.FromContext(checkedCtx).
				WithField("tenant_id", multitenancy.TenantIDFromContext(checkedCtx)).
				WithField("allowed_tenants", multitenancy.AllowedTenantsFromContext(checkedCtx)).
				Debugf("authentication succeeded (API-Key)")
			a.serveNext(rw, req.WithContext(checkedCtx), h)
			return
		}

		// Perform JWT Authentication
		checkedCtx, err = a.jwt.Check(authCtx)
		if err != nil {
			log.FromContext(authCtx).WithError(err).Errorf("unauthorized request")
			a.writeUnauthorized(rw, err)
			return
		}
		if checkedCtx != nil {
			// JWT Authentication succeeded
			log.FromContext(checkedCtx).
				WithField("tenant_id", multitenancy.TenantIDFromContext(checkedCtx)).
				WithField("allowed_tenants", multitenancy.AllowedTenantsFromContext(checkedCtx)).
				Debugf("authentication succeeded (JWT)")

			a.serveNext(rw, req.WithContext(checkedCtx), h)
			return
		}

		err = errors.UnauthorizedError("missing required credentials")
		log.FromContext(authCtx).WithError(err).Errorf("unauthorized request")
		a.writeUnauthorized(rw, err)
	})
}

func (a *Auth) writeUnauthorized(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusUnauthorized)
	_, _ = rw.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusUnauthorized, err.Error())))
}

func (a *Auth) serveNext(rw http.ResponseWriter, req *http.Request, h http.Handler) {
	// Remove authorization header
	// So possibly another Authorization will be set by Proxy
	req.Header.Del(authutils.AuthorizationHeader)
	req.Header.Del(authutils.APIKeyHeader)

	// Execute next handlers
	h.ServeHTTP(rw, req)
}
