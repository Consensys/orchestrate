// +build unit

package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/testutils"
)

func TestTransactionSchedulerMetrics(t *testing.T) {
	ep := NewTransactionSchedulerMetrics()

	registry := prometheus.NewRegistry()
	err := registry.Register(ep)
	assert.NoError(t, err, "Registering TransactionSchedulerMetrics should not fail")

	ep.CreatedJobsCounter().
		With("chain_uuid", "chain_uuid").
		With("tenant_id", "tenant_id").
		Add(1)

	families, err := registry.Gather()
	require.NoError(t, err, "Gathering metrics should not error")
	require.Len(t, families, 1, "Count of metrics families should be correct")

	testutils.AssertCounterFamily(t, families[0], Namespace, CreatedJobName, []float64{1}, "Total count of created jobs.", nil)
}
