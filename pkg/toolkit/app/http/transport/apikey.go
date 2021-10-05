package transport

import (
	"net/http"

	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
)

type XAPIKeyHeadersTransport struct {
	apiKey string
	T      http.RoundTripper
}

// NewXAPIKeyHeadersTransport creates a new transport to attach API-KEY as part of request headers
func NewXAPIKeyHeadersTransport(apiKey string) Middleware {
	return func(nxt http.RoundTripper) http.RoundTripper {
		return &XAPIKeyHeadersTransport{
			T:      nxt,
			apiKey: apiKey,
		}
	}
}

func (t *XAPIKeyHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if authutils.GetAuthorizationHeader(req) == "" {
		authutils.AddXAPIKeyHeaderValue(req, t.apiKey)
		multitenancy.AddTenantIDHeader(req)
	}

	return t.T.RoundTrip(req)
}
