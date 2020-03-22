package tcp

import (
	"sync/atomic"
)

type Switcher struct {
	handler *atomic.Value
}

type handlerValue struct {
	handler Handler
}

func NewSwitcher() *Switcher {
	return &Switcher{
		handler: &atomic.Value{},
	}
}

// ServeTCP forwards the TCP connection to the current active router
func (s *Switcher) ServeTCP(conn WriteCloser) {
	handler := s.handler.Load()
	v, ok := handler.(*handlerValue)
	if ok {
		v.handler.ServeTCP(conn)
	} else {
		conn.Close()
	}
}

// Switch sets the new router for new connections
func (s *Switcher) Switch(handler Handler) {
	s.handler.Store(&handlerValue{handler})
}
