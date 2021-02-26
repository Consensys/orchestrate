package tcp

import (
	"fmt"
	"net"
	"sync"
)

type Forwarder struct {
	net.Listener
	conns chan net.Conn

	closeOnce *sync.Once
	closed    chan struct{}
}

func NewForwarder(ln net.Listener) *Forwarder {
	return &Forwarder{
		Listener:  ln,
		conns:     make(chan net.Conn),
		closed:    make(chan struct{}),
		closeOnce: &sync.Once{},
	}
}

// ServeTCP uses the connection to serve it later in "Accept"
func (f *Forwarder) ServeTCP(conn WriteCloser) {
	select {
	case <-f.closed:
		close(f.conns)
	default:
		f.conns <- conn
	}
}

// Accept retrieves a served connection in ServeTCP
func (f *Forwarder) Accept() (net.Conn, error) {
	select {
	case <-f.closed:
		return nil, fmt.Errorf("listener closed")
	default:
		conn, ok := <-f.conns
		if !ok {
			return nil, fmt.Errorf("listener closed")
		}
		return conn, nil
	}
}

func (f *Forwarder) Close() error {
	f.closeOnce.Do(func() {
		close(f.closed)
		f.ServeTCP(nil)
	})
	return nil
}
