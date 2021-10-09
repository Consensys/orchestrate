// +build unit
// +build !race

package metrics

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	metrics1 "github.com/consensys/orchestrate/pkg/toolkit/app/metrics"
	"github.com/consensys/orchestrate/pkg/toolkit/app/metrics/testutils"
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPMetrics(t *testing.T) {
	httpCollector := NewHTTPMetrics(nil)

	registry := prometheus.NewRegistry()
	err := registry.Register(httpCollector)
	assert.NoError(t, err, "Registering HTTP should not fail")

	httpCollector.RequestsCounter().
		With("tenant_id", "test-tenant", "entrypoint", "app", "protocol", "http", "service", "service1", "method", http.MethodGet, "code", strconv.Itoa(http.StatusOK)).
		Add(1)

	httpCollector.TLSRequestsCounter().
		With("tenant_id", "test-tenant", "entrypoint", "app", "service", "service1", "tls_version", "1.2", "tls_cipher", "unknown").
		Add(1)

	httpCollector.RequestsLatencyHistogram().
		With("tenant_id", "test-tenant", "entrypoint", "app", "protocol", "http", "service", "service1", "method", http.MethodGet, "code", strconv.Itoa(http.StatusOK)).
		Observe(1)

	httpCollector.OpenConnsGauge().
		With("tenant_id", "test-tenant", "entrypoint", "app", "protocol", "hcv", "service", "service1", "method", http.MethodGet).
		Set(1)

	httpCollector.RetriesCounter().
		With("tenant_id", "test-tenant", "entrypoint", "app", "service", "service1").
		Add(1)

	httpCollector.ServerUpGauge().
		With("tenant_id", "test-tenant", "entrypoint", "app", "service", "service1", "url", "test-url").
		Set(1)

	time.Sleep(time.Second)

	families, err := registry.Gather()
	require.NoError(t, err, "Gathering metrics should not error")
	require.Len(t, families, 6, "Count of metrics families should be correct")

	testutils.AssertGaugeFamily(t, families[0], fmt.Sprintf("%s_%s", metrics1.Namespace, Subsystem), OpenConnections, []float64{1}, "OpenConns", nil)
	testutils.AssertHistogramFamily(t, families[1], fmt.Sprintf("%s_%s", metrics1.Namespace, Subsystem), RequestsLatencySeconds, []uint64{1}, "RequestsLatency", nil)
	testutils.AssertCounterFamily(t, families[2], fmt.Sprintf("%s_%s", metrics1.Namespace, Subsystem), RequestsTLSTotal, []float64{1}, "TLSRequests", nil)
	testutils.AssertCounterFamily(t, families[3], fmt.Sprintf("%s_%s", metrics1.Namespace, Subsystem), RequestsTotal, []float64{1}, "Requests", nil)
	testutils.AssertCounterFamily(t, families[4], fmt.Sprintf("%s_%s", metrics1.Namespace, Subsystem), RetriesTotal, []float64{1}, "Retries", nil)
	testutils.AssertGaugeFamily(t, families[5], fmt.Sprintf("%s_%s", metrics1.Namespace, Subsystem), ServerUp, []float64{1}, "ServerUp", nil)
}

func TestReloadConfiguration(t *testing.T) {
	httpCollector := NewHTTPMetrics(nil)

	registry := prometheus.NewRegistry()
	err := registry.Register(httpCollector)
	assert.NoError(t, err, "Registering HTTP should not fail")

	// First test with empty dynamic config
	httpCollector.RequestsCounter().
		With("tenant_id", "test-tenant", "entrypoint", "ep-foo", "protocol", "http", "service", "dashboard@provider1", "method", http.MethodGet, "code", strconv.Itoa(http.StatusOK)).
		Add(1)

	httpCollector.ServerUpGauge().
		With("tenant_id", "test-tenant", "entrypoint", "ep-foo", "service", "proxy@provider1", "url", "http://test.com").
		Set(1)

	// #1 Gather a 1st time: should retrieve metrics
	families, err := registry.Gather()
	require.NoError(t, err, "#1 Gathering metrics should not error")
	assert.Len(t, families, 2, "#1 Count of metrics families should be correct")

	// #2 Gather again: metrics should have been removed
	families, err = registry.Gather()
	require.NoError(t, err, "#2 Gathering metrics should not error")
	assert.Len(t, families, 0, "#2 Count of metrics families should be correct")

	// Second test with dynamic config set
	dynCfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"router-proxy@provider1": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"ep-foo"},
						Rule:        "Host(`proxy.com`)",
						Service:     "proxy@provider1",
					},
				},
			},
			Services: map[string]*dynamic.Service{
				"proxy@provider1": {
					ReverseProxy: &dynamic.ReverseProxy{
						LoadBalancer: &dynamic.LoadBalancer{
							Servers: []*dynamic.Server{
								&dynamic.Server{
									URL: "http://test.com",
								},
							},
						},
					},
				},
				"dashboard@provider1": {
					Dashboard: &dynamic.Dashboard{},
				},
			},
		},
	}

	_ = httpCollector.Switch(dynCfg)

	// Increase waiting time to complete dynamic cfg switch
	time.Sleep(time.Second)

	httpCollector.RequestsCounter().
		With("tenant_id", "test-tenant", "entrypoint", "ep-foo", "protocol", "http", "service", "dashboard@provider1", "method", http.MethodGet, "code", strconv.Itoa(http.StatusOK)).
		Add(1)

	httpCollector.RequestsCounter().
		With("tenant_id", "test-tenant", "entrypoint", "ep-foo", "protocol", "http", "service", "unknown", "method", http.MethodGet, "code", strconv.Itoa(http.StatusOK)).
		Add(1)

	httpCollector.ServerUpGauge().
		With("tenant_id", "test-tenant", "entrypoint", "ep-foo", "service", "proxy@provider1", "url", "http://test.com").
		Set(1)

	httpCollector.ServerUpGauge().
		With("tenant_id", "test-tenant", "entrypoint", "ep-foo", "service", "proxy@provider1", "url", "http://unknown.com").
		Set(1)

	time.Sleep(time.Second)

	// #3 Gather a 1st time: should retrieve metrics
	families, err = registry.Gather()
	require.NoError(t, err, "#3 Gathering metrics should not error")
	require.Len(t, families, 2, "#3 Count of metrics families should be correct")
	assert.Len(t, families[0].GetMetric(), 2, "#3 Count of requests metrics should be correct")
	assert.Len(t, families[1].GetMetric(), 2, "#3 Count of server up metrics should be correct")

	// #3 Gather a 2nd time: should retrieve metrics that are in dynamic config
	families, err = registry.Gather()
	require.NoError(t, err, "#4 Gathering metrics should not error")
	require.Len(t, families, 2, "#4 Count of metrics families should be correct")
	assert.Len(t, families[0].GetMetric(), 1, "#4 Count of requests metrics should be correct")
	assert.Len(t, families[1].GetMetric(), 1, "#4 Count of server up metrics should be correct")
}
