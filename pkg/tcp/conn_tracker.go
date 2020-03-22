package tcp

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
)

type ConnTracker struct {
	mux   *sync.RWMutex
	conns map[net.Conn]struct{}
}

func NewConnTracker() *ConnTracker {
	return &ConnTracker{
		mux:   &sync.RWMutex{},
		conns: make(map[net.Conn]struct{}),
	}
}

// AddConnection add a connection in the tracked connections list
func (t *ConnTracker) AddConnection(conn net.Conn) {
	t.mux.Lock()
	t.conns[conn] = struct{}{}
	t.mux.Unlock()
}

// RemoveConnection remove a connection from the tracked connections list
func (t *ConnTracker) RemoveConnection(conn net.Conn) {
	t.mux.Lock()
	delete(t.conns, conn)
	t.mux.Unlock()
}

func (t *ConnTracker) isEmpty() bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return len(t.conns) == 0
}

// Shutdown wait for the connection closing
func (t *ConnTracker) Shutdown(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		if t.isEmpty() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// Close close all the connections in the tracked connections list
func (t *ConnTracker) Close() error {
	t.mux.Lock()
	for conn := range t.conns {
		if err := conn.Close(); err != nil {
			log.WithoutContext().Errorf("Error while closing connection: %v", err)
		}
		delete(t.conns, conn)
	}
	t.mux.Unlock()
	return nil
}

type TrackedConn struct {
	WriteCloser
	tracker *ConnTracker
}

func NewTrackedConn(conn WriteCloser, tracker *ConnTracker) *TrackedConn {
	return &TrackedConn{
		WriteCloser: conn,
		tracker:     tracker,
	}
}

func (conn *TrackedConn) Close() error {
	conn.tracker.RemoveConnection(conn.WriteCloser)
	return conn.WriteCloser.Close()
}
