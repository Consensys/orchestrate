package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
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
		if a.multitenancy {
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

			// Extract API Key credentials from HTTP headers
			authCtx = multitenancy.WithTenantID(
				authCtx,
				req.Header.Get(multitenancy.TenantIDHeader),
			)

			// Perform API Key Authentication
			checkedCtx, err := a.key.Check(authCtx)
			if err == nil {
				// Bypass JWT authentication
				log.FromContext(checkedCtx).
					WithField("tenant_id", multitenancy.TenantIDFromContext(checkedCtx)).
					Debugf("authentication succeeded (API-Key)")
				a.serveNext(rw, req.WithContext(checkedCtx), h)
				return
			}

			// Perform JWT Authentication
			checkedCtx, err = a.jwt.Check(authCtx)
			if err == nil {
				// JWT Authentication succeeded
				log.FromContext(checkedCtx).
					WithField("tenant_id", multitenancy.TenantIDFromContext(checkedCtx)).
					Debugf("authentication succeeded (JWT)")

				a.serveNext(rw, req.WithContext(checkedCtx), h)
			} else {
				log.FromContext(checkedCtx).WithError(err).Errorf("authentication failed")
				a.writeUnauthorized(rw, err)
			}
		} else {
			a.serveNext(rw, req, h)
		}
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
