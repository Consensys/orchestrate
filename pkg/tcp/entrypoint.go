package tcp

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
)

// EntryPoint is a TCP server
type EntryPoint struct {
	addr string

	handler Handler

	tracker *ConnTracker

	lis *atomic.Value

	timeouts  *traefikstatic.RespondingTimeouts
	lifecycle *traefikstatic.LifeCycle
}

type listenerValue struct {
	l net.Listener
}

// NewEntryPoint creates a new EntryPoint
func NewEntryPoint(config *traefikstatic.EntryPoint, handler Handler) *EntryPoint {
	return &EntryPoint{
		addr:      config.Address,
		handler:   handler,
		tracker:   NewConnTracker(),
		timeouts:  config.Transport.RespondingTimeouts,
		lifecycle: config.Transport.LifeCycle,
		lis:       &atomic.Value{},
	}
}

func (e *EntryPoint) Addr() string {
	lis := e.listener()
	if lis != nil {
		return lis.Addr().String()
	}
	return ""
}

func (e *EntryPoint) listener() net.Listener {
	v, ok := e.lis.Load().(*listenerValue)
	if ok {
		return v.l
	}
	return nil
}

func (e *EntryPoint) Serve(ctx context.Context, l net.Listener) error {
	logger := log.FromContext(ctx)
	logger.WithField("address", l.Addr()).Infof("start serving tcp entrypoint")

	lis, ok := l.(KeepAliveListener)
	if !ok {
		lis = KeepAliveListener{l.(*net.TCPListener)}
	}

	e.lis.Store(&listenerValue{lis})

	for {
		conn, err := lis.Accept()
		if err != nil {
			logger.Error(err)
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				continue
			}

			return err
		}

		// TODO: switch this to Trace when updating logger interface
		logger.WithFields(logrus.Fields{
			"addr.remote": conn.RemoteAddr(),
			"addr.local":  conn.LocalAddr(),
		}).Debugf("accepted tcp connection %v", conn)

		writeCloser, err := writeCloser(conn)
		if err != nil {
			panic(err)
		}

		go func() {
			// Enforce read/write deadlines at the connection level,
			// because when we're peeking the first byte to determine whether we are doing TLS,
			// the deadlines at the server level are not taken into account.
			if e.timeouts.ReadTimeout > 0 {
				err = writeCloser.SetReadDeadline(time.Now().Add(time.Duration(e.timeouts.ReadTimeout)))
				if err != nil {
					logger.WithError(err).Errorf("could not set read deadline")
				}
			}

			if e.timeouts.WriteTimeout > 0 {
				err = writeCloser.SetWriteDeadline(time.Now().Add(time.Duration(e.timeouts.WriteTimeout)))
				if err != nil {
					logger.WithError(err).Errorf("could not set write deadline")
				}
			}

			e.handler.ServeTCP(NewTrackedConn(writeCloser, e.tracker))
		}()
	}
}

// Serve handler provided on entrypoint
func (e *EntryPoint) ListenAndServe(ctx context.Context) error {
	listener, err := Listen(e.addr)
	if err != nil {
		return err
	}
	return e.Serve(ctx, listener)
}

// Shutdown stops the TCP connections
func (e *EntryPoint) Shutdown(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Infof("shutting down tcp entrypoint...")

	reqAcceptGraceTimeOut := time.Duration(e.lifecycle.RequestAcceptGraceTimeout)
	if reqAcceptGraceTimeOut > 0 {
		logger.Infof("waiting %s for incoming requests to cease...", reqAcceptGraceTimeOut)
		time.Sleep(reqAcceptGraceTimeOut)
	}

	graceTimeOut := time.Duration(e.lifecycle.GraceTimeOut)
	ctx, cancel := context.WithTimeout(ctx, graceTimeOut)
	defer cancel()
	logger.Infof("waiting %s seconds before killing connections...", graceTimeOut)

	var wg sync.WaitGroup
	if handler, ok := e.handler.(Shutdownable); ok {
		wg.Add(1)
		go func() {
			_ = Shutdown(ctx, handler)
			wg.Done()
		}()
	}

	if e.tracker != nil {
		wg.Add(1)
		go func() {
			_ = Shutdown(ctx, e.tracker)
			wg.Done()
		}()
	}

	wg.Wait()

	lis := e.listener()
	if lis != nil {
		return lis.Close()
	}
	logger.Infof("tcp entrypoint closed")
	return nil
}

func (e *EntryPoint) Close() error {
	var wg sync.WaitGroup
	if handler, ok := e.handler.(io.Closer); ok {
		wg.Add(1)
		go func() {
			_ = Close(handler)
			wg.Done()
		}()
	}

	if e.tracker != nil {
		wg.Add(1)
		go func() {
			_ = Close(e.tracker)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}
