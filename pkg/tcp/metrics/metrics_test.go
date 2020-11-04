// +build unit

package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/testutils"
)

func TestTCPMetrics(t *testing.T) {
	ep := NewTCPMetrics(nil)

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

	testutils.AssertCounterFamily(t, families[0], Namespace, AcceptedConnsTotal, []float64{1}, "AcceptedConns", nil)
	testutils.AssertCounterFamily(t, families[1], Namespace, ClosedConnsTotal, []float64{1}, "ClosedConns", nil)
	testutils.AssertGaugeFamily(t, families[2], Namespace, OpenConns, []float64{1}, "OpenConns", nil)
	testutils.AssertHistogramFamily(t, families[3], Namespace, OpenConnsDurationSeconds, []uint64{1}, "ConnsLatency", nil)
}
