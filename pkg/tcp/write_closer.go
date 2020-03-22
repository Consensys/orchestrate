package tcp

import (
	"fmt"
	"net"

	proxyprotocol "github.com/c0va23/go-proxyprotocol"
)

//go:generate mockgen -source=write_closer.go -destination=mock/write_closer.go -package=mock

// WriteCloser describes a net.Conn with a CloseWrite method.
type WriteCloser interface {
	net.Conn
	// CloseWrite on a network connection, indicates that the issuer of the call
	// has terminated sending on that connection.
	// It corresponds to sending a FIN packet.
	CloseWrite() error
}

// writeCloserWrapper wraps together a connection, and the concrete underlying
// connection type that was found to satisfy WriteCloser.
type writeCloserWrapper struct {
	net.Conn
	writeCloser WriteCloser
}

func (c *writeCloserWrapper) CloseWrite() error {
	return c.writeCloser.CloseWrite()
}

// writeCloser returns the given connection, augmented with the WriteCloser
// implementation, if any was found within the underlying conn.
func writeCloser(conn net.Conn) (WriteCloser, error) {
	switch typedConn := conn.(type) {
	case *proxyprotocol.Conn:
		underlying, err := writeCloser(typedConn.Conn)
		if err != nil {
			return nil, err
		}
		return &writeCloserWrapper{writeCloser: underlying, Conn: typedConn}, nil
	case *net.TCPConn:
		return typedConn, nil
	default:
		return nil, fmt.Errorf("unknown connection type %T", typedConn)
	}
}
