package transport

import (
	"net/http"
)

type customHeadersTransport struct {
	headers map[string]string
	T       http.RoundTripper
}

func NewCustomHeadersTransport(headers map[string]string) Middleware {
	return func(nxt http.RoundTripper) http.RoundTripper {
		return &customHeadersTransport{
			T:       nxt,
			headers: headers,
		}
	}
}

func (t *customHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}

	return t.T.RoundTrip(req)
}
