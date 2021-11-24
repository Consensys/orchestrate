package transport

import (
	"net/http"

	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
)

type ContextAuthHeadersTransport struct {
	T http.RoundTripper
}

// NewContextAuthHeadersTransport creates a new transport to attach context authentication values into request headers
func NewContextAuthHeadersTransport() Middleware {
	return func(nxt http.RoundTripper) http.RoundTripper {
		return &ContextAuthHeadersTransport{
			T: nxt,
		}
	}
}

func (t *ContextAuthHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	userInfo := multitenancy.UserInfoValue(req.Context())
	if userInfo != nil && userInfo.AuthMode == multitenancy.AuthMethodJWT {
		if userInfo.AuthValue != "" {
			authutils.AddAuthorizationHeaderValue(req, userInfo.AuthValue)
		}
	}

	return t.T.RoundTrip(req)
}
