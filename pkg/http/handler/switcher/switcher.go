package switcher

import (
	"net/http"
	"sync/atomic"

	"github.com/ConsenSys/orchestrate/pkg/http/handler/dummy"
)

// Switcher allows hot switching of http.ServeMux
type Switcher struct {
	handler *atomic.Value
}

type handlerValue struct {
	handler http.Handler
}

// NewH builds a new instance of HTTPHandlerSwitcher
func New() *Switcher {
	switcher := &Switcher{
		handler: &atomic.Value{},
	}

	switcher.handler.Store(&handlerValue{&dummy.Dummy{}})
	return switcher
}

func (h *Switcher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	v, ok := h.handler.Load().(*handlerValue)
	if ok {
		v.handler.ServeHTTP(rw, req)
	}
}

// Handler returns current http Handler
func (h *Switcher) Handler() http.Handler {
	handler := h.handler.Load().(*handlerValue).handler
	return handler
}

// Switch safely switch current http handler
func (h *Switcher) Switch(handler http.Handler) {
	h.handler.Store(&handlerValue{handler})
}
