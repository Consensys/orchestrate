package prometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/prometheus/testutils"
)

func TestTCP(t *testing.T) {
	ep := NewTCP(testCfg)

	registry := prometheus.NewRegistry()
	err := registry.Register(ep)
	assert.NoError(t, err, "Registering TCP should not fail")

	ep.AcceptedConnsCounter().
		With("entrypoint", "http").
		Add(1)

	ep.ClosedConnsCounter().
		With("entrypoint", "http").
		Add(1)

	ep.ConnsLatencyHistogram().
		With("entrypoint", "http").
		Observe(1)

	ep.OpenConnsGauge().
		With("entrypoint", "http").
		Set(1)

	families, err := registry.Gather()
	require.NoError(t, err, "Gathering metrics should not error")
	require.Len(t, families, 4, "Count of metrics families should be correct")

	testutils.AssertCounterFamily(t, families[0], "tcp_accepted_conns_total", []float64{1}, "AcceptedConns")
	testutils.AssertCounterFamily(t, families[1], "tcp_closed_conns_total", []float64{1}, "ClosedConns")
	testutils.AssertGaugeFamily(t, families[2], "tcp_open_conns", []float64{1}, "OpenConns")
	testutils.AssertHistogramFamily(t, families[3], "tcp_open_conns_duration_seconds", []uint64{1}, "ConnsLatency")
}
