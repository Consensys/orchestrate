// Code generated by MockGen. DO NOT EDIT.
// Source: checker.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entities "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	reflect "reflect"
)

// MockChecker is a mock of Checker interface
type MockChecker struct {
	ctrl     *gomock.Controller
	recorder *MockCheckerMockRecorder
}

// MockCheckerMockRecorder is the mock recorder for MockChecker
type MockCheckerMockRecorder struct {
	mock *MockChecker
}

// NewMockChecker creates a new mock instance
func NewMockChecker(ctrl *gomock.Controller) *MockChecker {
	mock := &MockChecker{ctrl: ctrl}
	mock.recorder = &MockCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockChecker) EXPECT() *MockCheckerMockRecorder {
	return m.recorder
}

// Check mocks base method
func (m *MockChecker) Check(ctx context.Context, job *entities.Job) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, job)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check
func (mr *MockCheckerMockRecorder) Check(ctx, job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockChecker)(nil).Check), ctx, job)
}

// OnFailure mocks base method
func (m *MockChecker) OnFailure(ctx context.Context, job *entities.Job, jobErr error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnFailure", ctx, job, jobErr)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnFailure indicates an expected call of OnFailure
func (mr *MockCheckerMockRecorder) OnFailure(ctx, job, jobErr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnFailure", reflect.TypeOf((*MockChecker)(nil).OnFailure), ctx, job, jobErr)
}

// OnSuccess mocks base method
func (m *MockChecker) OnSuccess(ctx context.Context, job *entities.Job) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnSuccess", ctx, job)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnSuccess indicates an expected call of OnSuccess
func (mr *MockCheckerMockRecorder) OnSuccess(ctx, job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnSuccess", reflect.TypeOf((*MockChecker)(nil).OnSuccess), ctx, job)
}
