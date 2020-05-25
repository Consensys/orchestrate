// Code generated by MockGen. DO NOT EDIT.
// Source: send_tx.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entities "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	reflect "reflect"
)

// MockSendTxUseCase is a mock of SendTxUseCase interface.
type MockSendTxUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSendTxUseCaseMockRecorder
}

// MockSendTxUseCaseMockRecorder is the mock recorder for MockSendTxUseCase.
type MockSendTxUseCaseMockRecorder struct {
	mock *MockSendTxUseCase
}

// NewMockSendTxUseCase creates a new mock instance.
func NewMockSendTxUseCase(ctrl *gomock.Controller) *MockSendTxUseCase {
	mock := &MockSendTxUseCase{ctrl: ctrl}
	mock.recorder = &MockSendTxUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendTxUseCase) EXPECT() *MockSendTxUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockSendTxUseCase) Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, txRequest, tenantID)
	ret0, _ := ret[0].(*entities.TxRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockSendTxUseCaseMockRecorder) Execute(ctx, txRequest, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockSendTxUseCase)(nil).Execute), ctx, txRequest, tenantID)
}
