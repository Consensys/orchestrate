// Code generated by MockGen. DO NOT EDIT.
// Source: exported.go

// Package mock is a generated GoMock package.
package mock

import (
	metrics "github.com/go-kit/kit/metrics"
	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
	reflect "reflect"
)

// MockTPCMetrics is a mock of TPCMetrics interface
type MockTPCMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockTPCMetricsMockRecorder
}

// MockTPCMetricsMockRecorder is the mock recorder for MockTPCMetrics
type MockTPCMetricsMockRecorder struct {
	mock *MockTPCMetrics
}

// NewMockTPCMetrics creates a new mock instance
func NewMockTPCMetrics(ctrl *gomock.Controller) *MockTPCMetrics {
	mock := &MockTPCMetrics{ctrl: ctrl}
	mock.recorder = &MockTPCMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTPCMetrics) EXPECT() *MockTPCMetricsMockRecorder {
	return m.recorder
}

// AcceptedConnsCounter mocks base method
func (m *MockTPCMetrics) AcceptedConnsCounter() metrics.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptedConnsCounter")
	ret0, _ := ret[0].(metrics.Counter)
	return ret0
}

// AcceptedConnsCounter indicates an expected call of AcceptedConnsCounter
func (mr *MockTPCMetricsMockRecorder) AcceptedConnsCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptedConnsCounter", reflect.TypeOf((*MockTPCMetrics)(nil).AcceptedConnsCounter))
}

// ClosedConnsCounter mocks base method
func (m *MockTPCMetrics) ClosedConnsCounter() metrics.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClosedConnsCounter")
	ret0, _ := ret[0].(metrics.Counter)
	return ret0
}

// ClosedConnsCounter indicates an expected call of ClosedConnsCounter
func (mr *MockTPCMetricsMockRecorder) ClosedConnsCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClosedConnsCounter", reflect.TypeOf((*MockTPCMetrics)(nil).ClosedConnsCounter))
}

// ConnsLatencyHistogram mocks base method
func (m *MockTPCMetrics) ConnsLatencyHistogram() metrics.Histogram {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnsLatencyHistogram")
	ret0, _ := ret[0].(metrics.Histogram)
	return ret0
}

// ConnsLatencyHistogram indicates an expected call of ConnsLatencyHistogram
func (mr *MockTPCMetricsMockRecorder) ConnsLatencyHistogram() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnsLatencyHistogram", reflect.TypeOf((*MockTPCMetrics)(nil).ConnsLatencyHistogram))
}

// OpenConnsGauge mocks base method
func (m *MockTPCMetrics) OpenConnsGauge() metrics.Gauge {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenConnsGauge")
	ret0, _ := ret[0].(metrics.Gauge)
	return ret0
}

// OpenConnsGauge indicates an expected call of OpenConnsGauge
func (mr *MockTPCMetricsMockRecorder) OpenConnsGauge() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenConnsGauge", reflect.TypeOf((*MockTPCMetrics)(nil).OpenConnsGauge))
}

// Describe mocks base method
func (m *MockTPCMetrics) Describe(arg0 chan<- *prometheus.Desc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Describe", arg0)
}

// Describe indicates an expected call of Describe
func (mr *MockTPCMetricsMockRecorder) Describe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockTPCMetrics)(nil).Describe), arg0)
}

// Collect mocks base method
func (m *MockTPCMetrics) Collect(arg0 chan<- prometheus.Metric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Collect", arg0)
}

// Collect indicates an expected call of Collect
func (mr *MockTPCMetricsMockRecorder) Collect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect", reflect.TypeOf((*MockTPCMetrics)(nil).Collect), arg0)
}
