// Code generated by MockGen. DO NOT EDIT.
// Source: create_job.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entities "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	reflect "reflect"
)

// MockCreateJobUseCase is a mock of CreateJobUseCase interface.
type MockCreateJobUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCreateJobUseCaseMockRecorder
}

// MockCreateJobUseCaseMockRecorder is the mock recorder for MockCreateJobUseCase.
type MockCreateJobUseCaseMockRecorder struct {
	mock *MockCreateJobUseCase
}

// NewMockCreateJobUseCase creates a new mock instance.
func NewMockCreateJobUseCase(ctrl *gomock.Controller) *MockCreateJobUseCase {
	mock := &MockCreateJobUseCase{ctrl: ctrl}
	mock.recorder = &MockCreateJobUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateJobUseCase) EXPECT() *MockCreateJobUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockCreateJobUseCase) Execute(ctx context.Context, job *entities.Job, tenantID string) (*entities.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, job, tenantID)
	ret0, _ := ret[0].(*entities.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockCreateJobUseCaseMockRecorder) Execute(ctx, job, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCreateJobUseCase)(nil).Execute), ctx, job, tenantID)
}
