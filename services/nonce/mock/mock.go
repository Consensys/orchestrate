// Code generated by MockGen. DO NOT EDIT.
// Source: nonce.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAttributor is a mock of Attributor interface
type MockAttributor struct {
	ctrl     *gomock.Controller
	recorder *MockAttributorMockRecorder
}

// MockAttributorMockRecorder is the mock recorder for MockAttributor
type MockAttributorMockRecorder struct {
	mock *MockAttributor
}

// NewMockAttributor creates a new mock instance
func NewMockAttributor(ctrl *gomock.Controller) *MockAttributor {
	mock := &MockAttributor{ctrl: ctrl}
	mock.recorder = &MockAttributorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAttributor) EXPECT() *MockAttributorMockRecorder {
	return m.recorder
}

// GetLastAttributed mocks base method
func (m *MockAttributor) GetLastAttributed(key string) (uint64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAttributed", key)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLastAttributed indicates an expected call of GetLastAttributed
func (mr *MockAttributorMockRecorder) GetLastAttributed(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAttributed", reflect.TypeOf((*MockAttributor)(nil).GetLastAttributed), key)
}

// IncrLastAttributed mocks base method
func (m *MockAttributor) IncrLastAttributed(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrLastAttributed", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrLastAttributed indicates an expected call of IncrLastAttributed
func (mr *MockAttributorMockRecorder) IncrLastAttributed(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrLastAttributed", reflect.TypeOf((*MockAttributor)(nil).IncrLastAttributed), key)
}

// SetLastAttributed mocks base method
func (m *MockAttributor) SetLastAttributed(key string, value uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLastAttributed", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLastAttributed indicates an expected call of SetLastAttributed
func (mr *MockAttributorMockRecorder) SetLastAttributed(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastAttributed", reflect.TypeOf((*MockAttributor)(nil).SetLastAttributed), key, value)
}

// DeleteLastAttributed mocks base method
func (m *MockAttributor) DeleteLastAttributed(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastAttributed", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastAttributed indicates an expected call of DeleteLastAttributed
func (mr *MockAttributorMockRecorder) DeleteLastAttributed(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastAttributed", reflect.TypeOf((*MockAttributor)(nil).DeleteLastAttributed), key)
}

// MockSender is a mock of Sender interface
type MockSender struct {
	ctrl     *gomock.Controller
	recorder *MockSenderMockRecorder
}

// MockSenderMockRecorder is the mock recorder for MockSender
type MockSenderMockRecorder struct {
	mock *MockSender
}

// NewMockSender creates a new mock instance
func NewMockSender(ctrl *gomock.Controller) *MockSender {
	mock := &MockSender{ctrl: ctrl}
	mock.recorder = &MockSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSender) EXPECT() *MockSenderMockRecorder {
	return m.recorder
}

// GetLastSent mocks base method
func (m *MockSender) GetLastSent(key string) (uint64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastSent", key)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLastSent indicates an expected call of GetLastSent
func (mr *MockSenderMockRecorder) GetLastSent(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastSent", reflect.TypeOf((*MockSender)(nil).GetLastSent), key)
}

// IncrLastSent mocks base method
func (m *MockSender) IncrLastSent(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrLastSent", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrLastSent indicates an expected call of IncrLastSent
func (mr *MockSenderMockRecorder) IncrLastSent(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrLastSent", reflect.TypeOf((*MockSender)(nil).IncrLastSent), key)
}

// DeleteLastSent mocks base method
func (m *MockSender) DeleteLastSent(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastSent", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastSent indicates an expected call of DeleteLastSent
func (mr *MockSenderMockRecorder) DeleteLastSent(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastSent", reflect.TypeOf((*MockSender)(nil).DeleteLastSent), key)
}

// SetLastSent mocks base method
func (m *MockSender) SetLastSent(key string, value uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLastSent", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLastSent indicates an expected call of SetLastSent
func (mr *MockSenderMockRecorder) SetLastSent(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastSent", reflect.TypeOf((*MockSender)(nil).SetLastSent), key, value)
}

// IsRecovering mocks base method
func (m *MockSender) IsRecovering(key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRecovering", key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRecovering indicates an expected call of IsRecovering
func (mr *MockSenderMockRecorder) IsRecovering(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRecovering", reflect.TypeOf((*MockSender)(nil).IsRecovering), key)
}

// SetRecovering mocks base method
func (m *MockSender) SetRecovering(key string, status bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRecovering", key, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRecovering indicates an expected call of SetRecovering
func (mr *MockSenderMockRecorder) SetRecovering(key, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRecovering", reflect.TypeOf((*MockSender)(nil).SetRecovering), key, status)
}

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// GetLastAttributed mocks base method
func (m *MockManager) GetLastAttributed(key string) (uint64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAttributed", key)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLastAttributed indicates an expected call of GetLastAttributed
func (mr *MockManagerMockRecorder) GetLastAttributed(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAttributed", reflect.TypeOf((*MockManager)(nil).GetLastAttributed), key)
}

// IncrLastAttributed mocks base method
func (m *MockManager) IncrLastAttributed(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrLastAttributed", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrLastAttributed indicates an expected call of IncrLastAttributed
func (mr *MockManagerMockRecorder) IncrLastAttributed(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrLastAttributed", reflect.TypeOf((*MockManager)(nil).IncrLastAttributed), key)
}

// SetLastAttributed mocks base method
func (m *MockManager) SetLastAttributed(key string, value uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLastAttributed", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLastAttributed indicates an expected call of SetLastAttributed
func (mr *MockManagerMockRecorder) SetLastAttributed(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastAttributed", reflect.TypeOf((*MockManager)(nil).SetLastAttributed), key, value)
}

// DeleteLastAttributed mocks base method
func (m *MockManager) DeleteLastAttributed(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastAttributed", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastAttributed indicates an expected call of DeleteLastAttributed
func (mr *MockManagerMockRecorder) DeleteLastAttributed(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastAttributed", reflect.TypeOf((*MockManager)(nil).DeleteLastAttributed), key)
}

// GetLastSent mocks base method
func (m *MockManager) GetLastSent(key string) (uint64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastSent", key)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLastSent indicates an expected call of GetLastSent
func (mr *MockManagerMockRecorder) GetLastSent(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastSent", reflect.TypeOf((*MockManager)(nil).GetLastSent), key)
}

// IncrLastSent mocks base method
func (m *MockManager) IncrLastSent(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrLastSent", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrLastSent indicates an expected call of IncrLastSent
func (mr *MockManagerMockRecorder) IncrLastSent(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrLastSent", reflect.TypeOf((*MockManager)(nil).IncrLastSent), key)
}

// DeleteLastSent mocks base method
func (m *MockManager) DeleteLastSent(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastSent", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastSent indicates an expected call of DeleteLastSent
func (mr *MockManagerMockRecorder) DeleteLastSent(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastSent", reflect.TypeOf((*MockManager)(nil).DeleteLastSent), key)
}

// SetLastSent mocks base method
func (m *MockManager) SetLastSent(key string, value uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLastSent", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLastSent indicates an expected call of SetLastSent
func (mr *MockManagerMockRecorder) SetLastSent(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastSent", reflect.TypeOf((*MockManager)(nil).SetLastSent), key, value)
}

// IsRecovering mocks base method
func (m *MockManager) IsRecovering(key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRecovering", key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRecovering indicates an expected call of IsRecovering
func (mr *MockManagerMockRecorder) IsRecovering(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRecovering", reflect.TypeOf((*MockManager)(nil).IsRecovering), key)
}

// SetRecovering mocks base method
func (m *MockManager) SetRecovering(key string, status bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRecovering", key, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRecovering indicates an expected call of SetRecovering
func (mr *MockManagerMockRecorder) SetRecovering(key, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRecovering", reflect.TypeOf((*MockManager)(nil).SetRecovering), key, status)
}
