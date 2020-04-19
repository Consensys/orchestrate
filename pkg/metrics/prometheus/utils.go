package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

func ToLabels(labelValues ...string) prometheus.Labels {
	if len(labelValues)%2 != 0 {
		labelValues = append(labelValues, "unknown")
	}

	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}

	return labels
}
