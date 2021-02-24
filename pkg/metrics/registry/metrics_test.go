// +build unit

package registry

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
	"github.com/ConsenSys/orchestrate/pkg/metrics/mock"
)

func TestRegistryMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	metricOne := mock.NewMockPrometheus(ctrl)
	metricTwo := mock.NewMockDynamicPrometheus(ctrl)
	metricThree := mock.NewMockPrometheus(ctrl)
	metricFour := mock.NewMockDynamicPrometheus(ctrl)

	reg := NewMetricRegistry(metricOne, metricTwo)
	reg.Add(metricThree)
	reg.Add(metricFour)
	
	metricOne.EXPECT().Collect(gomock.Any())
	metricTwo.EXPECT().Collect(gomock.Any())
	metricThree.EXPECT().Collect(gomock.Any())
	metricFour.EXPECT().Collect(gomock.Any())
	reg.Collect(make(chan prometheus.Metric, 1))
	
	metricOne.EXPECT().Describe(gomock.Any())
	metricTwo.EXPECT().Describe(gomock.Any())
	metricThree.EXPECT().Describe(gomock.Any())
	metricFour.EXPECT().Describe(gomock.Any())
	reg.Describe(make(chan *prometheus.Desc, 1))
	
	cfg := &dynamic.Configuration{}
	metricTwo.EXPECT().Switch(gomock.Eq(cfg))
	metricFour.EXPECT().Switch(gomock.Eq(cfg))
	err := reg.SwitchDynConfig(cfg)
	assert.NoError(t, err)
}
