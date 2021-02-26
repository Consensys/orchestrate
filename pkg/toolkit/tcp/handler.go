package tcp

import (
	"crypto/tls"
)

//go:generate mockgen -source=handler.go  -destination=mock/handler.go -package=mock

// Handler is the TCP Handlers interface
type Handler interface {
	ServeTCP(conn WriteCloser)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as handlers.
type HandlerFunc func(conn WriteCloser)

// ServeTCP serves tcp
func (f HandlerFunc) ServeTCP(conn WriteCloser) {
	f(conn)
}

// TLSHandler handles TLS connections
type TLSHandler struct {
	Next   Handler
	Config *tls.Config
}

// ServeTCP terminates the TLS connection
func (t *TLSHandler) ServeTCP(conn WriteCloser) {
	t.Next.ServeTCP(tls.Server(conn, t.Config))
}
