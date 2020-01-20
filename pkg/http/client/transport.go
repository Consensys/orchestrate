package httpclient

import (
	"net/http"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

// Transport is an internal Transport for Orchestrate
type Transport struct {
	T http.RoundTripper
}

// NewTransport creates a new transport
func NewTransport(t http.RoundTripper) *Transport {
	return &Transport{
		T: t,
	}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.addAuthorizationHeader(req)
	return t.T.RoundTrip(req)
}

func (t *Transport) addAuthorizationHeader(req *http.Request) {
	auth := authutils.AuthorizationFromContext(req.Context())
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}
}
