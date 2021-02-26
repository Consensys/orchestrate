package testutils

import (
	"fmt"

	dto "github.com/prometheus/client_model/go"
)

func generateFamilyName(namespace, name string) string {
	return fmt.Sprintf("%s_%s", namespace, name)
}

func filterMetricFamily(family *dto.MetricFamily, labels map[string]string) []*dto.Metric {
	result := []*dto.Metric{}
	for _, metric := range family.GetMetric() {
		if containLabelValue(metric.Label, labels) {
			result = append(result, metric)
		}
	}

	return result
}

// Returns true if all labels are matched
func containLabelValue(metricLabels []*dto.LabelPair, filterLabels map[string]string) bool {
	if filterLabels == nil {
		return true
	}

	matchedLabels := 0
	for _, l := range metricLabels {
		if l == nil || l.Name == nil || l.Value == nil {
			continue
		}

		if value, ok := filterLabels[*l.Name]; ok && *l.Value == value {
			matchedLabels++
			break
		}
	}

	return len(filterLabels) == matchedLabels
}
