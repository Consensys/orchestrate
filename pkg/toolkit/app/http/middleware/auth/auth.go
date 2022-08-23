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
	"github.com/traefik/traefik/v2/pkg/log"
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

func (b *Builder) Build(_ context.Context, _ string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.Auth)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	m := New(b.jwt, b.key, b.multitenancy)
	return m.Handler, nil, nil
}

type Auth struct {
	checker      auth.Checker
	multitenancy bool
}

func New(jwt, key auth.Checker, multitenancyEnabled bool) *Auth {
	return &Auth{
		checker:      auth.NewCombineCheckers(key, jwt),
		multitenancy: multitenancyEnabled,
	}
}

func (a *Auth) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !a.multitenancy {
			userInfo := multitenancy.DefaultUser()
			a.serveNext(rw, req.WithContext(multitenancy.WithUserInfo(req.Context(), userInfo)), h)
			return
		}

		// Extract Authorization credentials from HTTP headers
		authCtx := authutils.WithAuthorization(
			req.Context(),
			authutils.GetAuthorizationHeader(req),
		)

		// Extract API Key credentials from HTTP headers
		authCtx = authutils.WithAPIKey(
			authCtx,
			authutils.GetAPIKeyHeaderValue(req),
		)

		// Extract TenantID from HTTP headers
		authCtx = authutils.WithTenantID(
			authCtx,
			authutils.GetTenantIDHeaderValue(req),
		)

		// Extract Username from HTTP headers
		authCtx = authutils.WithUsername(
			authCtx,
			authutils.GetUsernameHeaderValue(req),
		)

		userInfo, err := a.checker.Check(authCtx)
		if err != nil {
			log.FromContext(authCtx).WithError(err).Errorf("unauthorized request")
			a.writeUnauthorized(rw, err)
			return
		}

		if userInfo != nil {
			// Bypass JWT authentication
			log.FromContext(authCtx).
				WithField("tenant_id", userInfo.TenantID).
				WithField("username", userInfo.Username).
				WithField("allowed_tenants", userInfo.AllowedTenants).
				Debugf("authentication succeeded (%s)", userInfo.AuthMode)

			a.serveNext(rw, req.WithContext(multitenancy.WithUserInfo(authCtx, userInfo)), h)
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
	authutils.DeleteAuthorizationHeaderValue(req)
	authutils.DeleteAPIKeyHeaderValue(req)

	// Execute next handlers
	h.ServeHTTP(rw, req)
}
