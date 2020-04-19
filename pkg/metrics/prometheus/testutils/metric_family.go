package testutils

import (
	"fmt"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertCounter(t *testing.T, metric *dto.Metric, value float64, desc string) {
	assert.Equal(t, value, metric.GetCounter().GetValue(), desc+": invalid Counter value")
}

func AssertGauge(t *testing.T, metric *dto.Metric, value float64, desc string) {
	assert.Equal(t, value, metric.GetGauge().GetValue(), desc+": invalid Gauge value")
}

func AssertHistogram(t *testing.T, metric *dto.Metric, sampleCount uint64, desc string) {
	assert.Equal(t, sampleCount, metric.GetHistogram().GetSampleCount(), desc+": invalid Histogram sample count")
}

func AssertMetricFamily(t *testing.T, family *dto.MetricFamily, name string, typ dto.MetricType, desc string) {
	assert.Equal(t, typ, family.GetType(), desc+": invalid metric type")
	assert.Equal(t, name, family.GetName(), desc+": invalid name")
}

func AssertCounterFamily(t *testing.T, family *dto.MetricFamily, name string, values []float64, desc string) {
	AssertMetricFamily(t, family, name, dto.MetricType_COUNTER, desc)
	require.Len(t, family.GetMetric(), len(values), desc+": invalid count of metrics")
	for i, metric := range family.GetMetric() {
		AssertCounter(t, metric, values[i], desc+fmt.Sprintf(" (#%v)", i))
	}
}

func AssertGaugeFamily(t *testing.T, family *dto.MetricFamily, name string, values []float64, desc string) {
	AssertMetricFamily(t, family, name, dto.MetricType_GAUGE, desc)
	require.Len(t, family.GetMetric(), len(values), desc+": invalid count of metrics")
	for i, metric := range family.GetMetric() {
		AssertGauge(t, metric, values[i], desc+fmt.Sprintf(" (#%v)", i))
	}
}

func AssertHistogramFamily(t *testing.T, family *dto.MetricFamily, name string, sampleCounts []uint64, desc string) {
	AssertMetricFamily(t, family, name, dto.MetricType_HISTOGRAM, desc)
	require.Len(t, family.GetMetric(), len(sampleCounts), desc+": invalid count of metrics")
	for i, metric := range family.GetMetric() {
		AssertHistogram(t, metric, sampleCounts[i], desc+fmt.Sprintf(" (#%v)", i))
	}
}
