// Code generated by MockGen. DO NOT EDIT.
// Source: update_job.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	reflect "reflect"
)

// MockUpdateJobUseCase is a mock of UpdateJobUseCase interface
type MockUpdateJobUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateJobUseCaseMockRecorder
}

// MockUpdateJobUseCaseMockRecorder is the mock recorder for MockUpdateJobUseCase
type MockUpdateJobUseCaseMockRecorder struct {
	mock *MockUpdateJobUseCase
}

// NewMockUpdateJobUseCase creates a new mock instance
func NewMockUpdateJobUseCase(ctrl *gomock.Controller) *MockUpdateJobUseCase {
	mock := &MockUpdateJobUseCase{ctrl: ctrl}
	mock.recorder = &MockUpdateJobUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUpdateJobUseCase) EXPECT() *MockUpdateJobUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockUpdateJobUseCase) Execute(ctx context.Context, jobEntity *types.Job, newStatus, logMessage string, tenants []string) (*types.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, jobEntity, newStatus, logMessage, tenants)
	ret0, _ := ret[0].(*types.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockUpdateJobUseCaseMockRecorder) Execute(ctx, jobEntity, newStatus, logMessage, tenants interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockUpdateJobUseCase)(nil).Execute), ctx, jobEntity, newStatus, logMessage, tenants)
}
