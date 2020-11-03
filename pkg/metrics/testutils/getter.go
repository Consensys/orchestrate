package testutils

import (
	"fmt"

	dto "github.com/prometheus/client_model/go"
)

func CounterValue(metric *dto.Metric) float64 {
	return metric.GetCounter().GetValue()
}

func GaugeValue(metric *dto.Metric) float64 {
	return metric.GetGauge().GetValue()
}

func HistogramValue(metric *dto.Metric) uint64 {
	return metric.GetHistogram().GetSampleCount()
}

func CounterFamilyValue(family *dto.MetricFamily) []float64 {
	r := []float64{}
	for _, m := range family.GetMetric() {
		r = append(r, CounterValue(m))
	}
	return r
}

func GaugeFamilyValue(family *dto.MetricFamily) []float64 {
	r := []float64{}
	for _, m := range family.GetMetric() {
		r = append(r, GaugeValue(m))
	}
	return r
}

func HistogramFamilyValue(family *dto.MetricFamily) []uint64 {
	r := []uint64{}
	for _, m := range family.GetMetric() {
		r = append(r, HistogramValue(m))
	}
	return r
}

func FamilyValue(families map[string]*dto.MetricFamily, namespace, name string) (interface{}, error) {
	mf, ok := families[generateFamilyName(namespace, name)]
	if !ok {
		return nil, fmt.Errorf("metric family does not exists")
	}

	switch mf.GetType() {
	case dto.MetricType_COUNTER:
		return CounterFamilyValue(mf), nil
	case dto.MetricType_GAUGE:
		return GaugeFamilyValue(mf), nil
	case dto.MetricType_HISTOGRAM:
		return HistogramFamilyValue(mf), nil
	default:
		return nil, fmt.Errorf("invalid metric type: %s", mf.GetType().String())
	}
}
