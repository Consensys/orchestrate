package transport

import (
	"net/http"

	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
)

type AuthHeadersTransport struct {
	auth string
	T    http.RoundTripper
}

func NewAuthHeadersTransport(jwt string) Middleware {
	return func(nxt http.RoundTripper) http.RoundTripper {
		return &AuthHeadersTransport{
			T:    nxt,
			auth: jwt,
		}
	}
}

func (t *AuthHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if authutils.GetAuthorizationHeader(req) == "" {
		authutils.AddAuthorizationHeaderValue(req, t.auth)
	}

	return t.T.RoundTrip(req)
}
