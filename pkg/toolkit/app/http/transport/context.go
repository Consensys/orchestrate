package transport

import (
	"net/http"

	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
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
	if authutils.GetAuthorizationHeader(req) == "" {
		authorization := authutils.AuthorizationFromContext(req.Context())
		if authorization != "" {
			authutils.AddAuthorizationHeaderValue(req, authorization)
		}
	}

	return t.T.RoundTrip(req)
}
