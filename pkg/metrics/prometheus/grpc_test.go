package prometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/prometheus/testutils"
)

func TestGRPCServer(t *testing.T) {
	grpcCollector := NewGRPCServer(testCfg)

	registry := prometheus.NewRegistry()
	err := registry.Register(grpcCollector)
	assert.NoError(t, err, "Registering GRPCServer should not fail")

	grpcCollector.StartedCounter().
		With("tenant_id", "foo", "type", "test-type", "service", "test-service", "method", "test-method").
		Add(1)

	grpcCollector.HandledCounter().
		With("tenant_id", "foo", "type", "test-type", "service", "test-service", "method", "test-method", "code", "test-code").
		Add(1)

	grpcCollector.StreamMsgReceivedCounter().
		With("tenant_id", "foo", "type", "test-type", "service", "test-service", "method", "test-method").
		Add(1)

	grpcCollector.StreamMsgSentCounter().
		With("tenant_id", "foo", "type", "test-type", "service", "test-service", "method", "test-method").
		Add(1)

	grpcCollector.HandledDurationHistogram().
		With("tenant_id", "foo", "type", "test-type", "service", "test-service", "method", "test-method", "code", "OK").
		Observe(1)

	families, err := registry.Gather()
	require.NoError(t, err, "Gathering metrics should not error")
	require.Len(t, families, 5, "Count of metrics families should be correct")

	testutils.AssertHistogramFamily(t, families[0], "grpc_server_handled_seconds", []uint64{1}, "HandledLatency")
	testutils.AssertCounterFamily(t, families[1], "grpc_server_handled_total", []float64{1}, "Handled")
	testutils.AssertCounterFamily(t, families[2], "grpc_server_msg_received_total", []float64{1}, "StreamMsgReceived")
	testutils.AssertCounterFamily(t, families[3], "grpc_server_msg_sent_total", []float64{1}, "StreamMsgSent")
	testutils.AssertCounterFamily(t, families[4], "grpc_server_started_total", []float64{1}, "Started")
}
