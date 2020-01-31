package auth

import (
	"fmt"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/gorilla/mux"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/key"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

type Auth struct {
	next         http.Handler
	authjwt      *authjwt.Auth
	authkey      *authkey.Auth
	multitenancy bool
}

func New(
	jwt *authjwt.Auth,
	key *authkey.Auth,
	multiEnabled bool,
	next http.Handler,
) *Auth {
	return &Auth{
		next:         next,
		authjwt:      jwt,
		authkey:      key,
		multitenancy: multiEnabled,
	}
}

func (a *Auth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if a.multitenancy {

		// Extract API Key credentials from HTTP headers
		apiKey := req.Header.Get(authentication.APIKeyHeader)
		apiKeyCtx := authutils.WithAPIKey(req.Context(), apiKey)

		// Perform API Key Authentication
		checkedCtx, err := a.authkey.Check(apiKeyCtx)
		if err == nil {
			// Bypass JWT authentication
			log.FromContext(req.Context()).Debugf("API Key Authentication succeeded")
			a.serveNext(rw, req.WithContext(checkedCtx))
			return
		}

		// Extract Authorization credentials from HTTP headers
		authorization := req.Header.Get(authentication.AuthorizationHeader)
		authCtx := authutils.WithAuthorization(req.Context(), authorization)

		// Perform JWT Authentication
		checkedCtx, err = a.authjwt.Check(authCtx)
		if err == nil {
			// JWT Authentication succeeded
			// We now control that tenantID passed in request pass and JWT token are the same
			jwtTenantID := multitenancy.TenantIDFromContext(checkedCtx)
			pathTenantID, ok := mux.Vars(req)["tenantID"]
			if !ok || pathTenantID != jwtTenantID {
				log.FromContext(req.Context()).Errorf("Permissioned to access tenant denied")
				a.writeNotFound(rw)
			} else {
				log.FromContext(req.Context()).Debugf("JWT Authentication succeeded")
				a.serveNext(rw, req.WithContext(checkedCtx))
			}
		} else {
			// JWT Authentication failed
			log.FromContext(req.Context()).WithError(err).Errorf("Authentication failed")
			a.writeUnauthorized(rw, err)
		}
	} else {
		// Auth is deactivated
		a.serveNext(rw, req)
	}
}

func (a *Auth) writeNotFound(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusNotFound)
	_, _ = rw.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusNotFound, http.StatusText(http.StatusNotFound))))
}

func (a *Auth) writeUnauthorized(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusUnauthorized)
	_, _ = rw.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusUnauthorized, err.Error())))
}

func (a *Auth) serveNext(rw http.ResponseWriter, req *http.Request) {
	// Remove authorization header
	// So possibly another Authorization will be set by Proxy
	req.Header.Del(authentication.AuthorizationHeader)
	req.Header.Del(authentication.APIKeyHeader)

	// Execute next handlers
	a.next.ServeHTTP(rw, req)
}
