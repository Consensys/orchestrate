// +build unit

package metrics

import (
	"fmt"
	"testing"

	metrics1 "github.com/ConsenSys/orchestrate/pkg/toolkit/app/metrics"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/metrics/testutils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionSchedulerMetrics(t *testing.T) {
	ep := NewTransactionSchedulerMetrics()

	registry := prometheus.NewRegistry()
	err := registry.Register(ep)
	assert.NoError(t, err, "Registering TransactionSchedulerMetrics should not fail")

	ep.JobsLatencyHistogram().
		With("chain_uuid", "chain_uuid").
		With("status", "started").
		With("prev_status", "created").
		Observe(1)

	ep.MinedLatencyHistogram().
		With("chain_uuid", "chain_uuid").
		With("status", "started").
		With("prev_status", "created").
		Observe(1)

	families, err := registry.Gather()
	require.NoError(t, err, "Gathering metrics should not error")
	require.Len(t, families, 2, "Count of metrics families should be correct")

	testutils.AssertHistogramFamily(t, families[0], fmt.Sprintf("%s_%s", metrics1.Namespace, Subsystem), JobLatencySeconds, []uint64{1}, "Histogram of job latency between status (second). Except PENDING and MINED, see mined_latency_seconds", nil)
	testutils.AssertHistogramFamily(t, families[1], fmt.Sprintf("%s_%s", metrics1.Namespace, Subsystem), MinedLatencySeconds, []uint64{1}, "Histogram of latency between PENDING and MINED (second)", nil)
}
