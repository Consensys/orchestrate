package transport

import (
	"net/http"

	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	authutils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/utils"
)

type APIKeyHeadersTransport struct {
	apiKey string
	T      http.RoundTripper
}

// NewAPIKeyHeadersTransport creates a new transport to attach API-KEY as part of request headers
func NewAPIKeyHeadersTransport(apiKey string) Middleware {
	return func(nxt http.RoundTripper) http.RoundTripper {
		return &APIKeyHeadersTransport{
			T:      nxt,
			apiKey: apiKey,
		}
	}
}

func (t *APIKeyHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if authutils.GetAuthorizationHeader(req) == "" {
		authutils.AddXAPIKeyHeaderValue(req, t.apiKey)
		multitenancy.AddTenantIDHeader(req)
	}

	return t.T.RoundTrip(req)
}
