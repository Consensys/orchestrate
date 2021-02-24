package http

import (
	"net/http"

	"github.com/ConsenSys/orchestrate/pkg/http/transport"
)

func NewTransport(t http.RoundTripper, middleware ...transport.Middleware) http.RoundTripper {
	if len(middleware) == 0 {
		return t
	} else if len(middleware) == 1 {
		return middleware[0](t)
	}

	return middleware[0](NewTransport(t, middleware[1:]...))
}
