package httpclient

import (
	"net/http"
	"strconv"
	"time"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

// AuthHeaderTransport is an internal Transport for Orchestrate
type AuthHeadersTransport struct {
	T http.RoundTripper
}

// NewAuthHeadersTransport creates a new transport
func NewAuthHeadersTransport(t http.RoundTripper) *AuthHeadersTransport {
	return &AuthHeadersTransport{
		T: t,
	}
}

func (t *AuthHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.addAuthorizationHeaders(req)
	return t.T.RoundTrip(req)
}

func (t *AuthHeadersTransport) addAuthorizationHeaders(req *http.Request) {
	authutils.AddAuthorizationHeader(req)
	authutils.AddXAPIKeyHeader(req)
}

type Retry429Transport struct {
	T http.RoundTripper
}

func NewRetry429Transport(t http.RoundTripper) *Retry429Transport {
	return &Retry429Transport{
		T: t,
	}
}

func (t *Retry429Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for {
		resp, err := t.T.RoundTrip(req)
		if err != nil {
			return resp, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter, _ := strconv.ParseInt(
				resp.Header.Get("Retry-After"),
				10, 64,
			)

			if retryAfter > 0 {
				select {
				case <-time.After(time.Duration(1000000000 * retryAfter)):
					continue
				case <-req.Context().Done():
					return nil, req.Context().Err()
				}
			}
		}

		return resp, nil
	}
}

func NewTransport(t http.RoundTripper) http.RoundTripper {
	return NewAuthHeadersTransport(NewRetry429Transport(t))
}
