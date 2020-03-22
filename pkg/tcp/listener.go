package tcp

import (
	"fmt"
	"net"
	"time"
)

func Listen(addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error opening listener: %v", err)
	}

	return KeepAliveListener{listener.(*net.TCPListener)}, nil
}

// KeepAliveListener sets TCP keep-alive timeouts on accepted
// connections.
type KeepAliveListener struct {
	*net.TCPListener
}

func (ln KeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if err := tc.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err := tc.SetKeepAlivePeriod(3 * time.Minute); err != nil {
		return nil, err
	}

	return tc, nil
}

func (ln KeepAliveListener) Close() error {
	return ln.TCPListener.Close()
}
