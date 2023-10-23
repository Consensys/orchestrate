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

		switch m := nxt.(type) {
		//if the next round tripper is already a customHeadersTransport, just update the headers
		case *customHeadersTransport:
			m.headers = headers
			return m
		default:
			return &customHeadersTransport{
				T:       nxt,
				headers: headers,
			}
		}

	}
}

func (t *customHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}

	return t.T.RoundTrip(req)
}
