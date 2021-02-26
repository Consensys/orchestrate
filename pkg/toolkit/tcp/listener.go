package tcp

import (
	"net"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/tcp/metrics"
	kitmetrics "github.com/go-kit/kit/metrics"
)

func Listen(network, addr string, opts ...ListenerOpt) (net.Listener, error) {
	l, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		l, err = opt(l)
		if err != nil {
			return nil, err
		}
	}

	return l, nil
}

type ListenerOpt func(l net.Listener) (net.Listener, error)

func KeepAliveOpt(period time.Duration) ListenerOpt {
	return func(l net.Listener) (net.Listener, error) {
		return &KeepAliveListener{l.(*net.TCPListener), period}, nil
	}
}

// KeepAliveListener sets TCP keep-alive timeouts on accepted
// connections.
type KeepAliveListener struct {
	*net.TCPListener
	period time.Duration
}

func (l *KeepAliveListener) Accept() (net.Conn, error) {
	tc, err := l.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if err := tc.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err := tc.SetKeepAlivePeriod(l.period); err != nil {
		return nil, err
	}

	return tc, nil
}

func (l *KeepAliveListener) Close() error {
	return l.TCPListener.Close()
}

func MetricsOpt(name string, registry metrics.TPCMetrics) ListenerOpt {
	return func(l net.Listener) (net.Listener, error) {
		return &MetricsListener{
			Listener: l,
			labels:   []string{"entrypoint", name},
			registry: registry,
		}, nil
	}
}

type MetricsListener struct {
	net.Listener
	labels   []string
	registry metrics.TPCMetrics
}

func (l *MetricsListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}

	l.registry.AcceptedConnsCounter().With(l.labels...).Add(1)
	openConnsGauge := l.registry.OpenConnsGauge().With(l.labels...)
	openConnsGauge.Add(1)

	return &trackedConn{
		Conn:                   conn,
		start:                  time.Now(),
		closedConnCounter:      l.registry.ClosedConnsCounter().With(l.labels...),
		openConnsGauge:         openConnsGauge,
		connsDurationHistogram: l.registry.ConnsLatencyHistogram().With(l.labels...),
	}, nil
}

type trackedConn struct {
	net.Conn
	start time.Time

	closedConnCounter      kitmetrics.Counter
	openConnsGauge         kitmetrics.Gauge
	connsDurationHistogram kitmetrics.Histogram
}

func (conn *trackedConn) Close() error {
	conn.closedConnCounter.Add(1)
	conn.openConnsGauge.Add(-1)
	conn.connsDurationHistogram.Observe(time.Since(conn.start).Seconds())
	return conn.Conn.Close()
}
