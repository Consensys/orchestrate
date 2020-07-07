// Code generated by MockGen. DO NOT EDIT.
// Source: transactions.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	entities "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	reflect "reflect"
)

// MockTransactionValidator is a mock of TransactionValidator interface
type MockTransactionValidator struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionValidatorMockRecorder
}

// MockTransactionValidatorMockRecorder is the mock recorder for MockTransactionValidator
type MockTransactionValidatorMockRecorder struct {
	mock *MockTransactionValidator
}

// NewMockTransactionValidator creates a new mock instance
func NewMockTransactionValidator(ctrl *gomock.Controller) *MockTransactionValidator {
	mock := &MockTransactionValidator{ctrl: ctrl}
	mock.recorder = &MockTransactionValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTransactionValidator) EXPECT() *MockTransactionValidatorMockRecorder {
	return m.recorder
}

// ValidateFields mocks base method
func (m *MockTransactionValidator) ValidateFields(ctx context.Context, txRequest *entities.TxRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateFields", ctx, txRequest)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateFields indicates an expected call of ValidateFields
func (mr *MockTransactionValidatorMockRecorder) ValidateFields(ctx, txRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateFields", reflect.TypeOf((*MockTransactionValidator)(nil).ValidateFields), ctx, txRequest)
}

// ValidateChainExists mocks base method
func (m *MockTransactionValidator) ValidateChainExists(ctx context.Context, chainUUID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateChainExists", ctx, chainUUID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateChainExists indicates an expected call of ValidateChainExists
func (mr *MockTransactionValidatorMockRecorder) ValidateChainExists(ctx, chainUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateChainExists", reflect.TypeOf((*MockTransactionValidator)(nil).ValidateChainExists), ctx, chainUUID)
}

// ValidateMethodSignature mocks base method
func (m *MockTransactionValidator) ValidateMethodSignature(methodSignature string, args []interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateMethodSignature", methodSignature, args)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateMethodSignature indicates an expected call of ValidateMethodSignature
func (mr *MockTransactionValidatorMockRecorder) ValidateMethodSignature(methodSignature, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateMethodSignature", reflect.TypeOf((*MockTransactionValidator)(nil).ValidateMethodSignature), methodSignature, args)
}

// ValidateContract mocks base method
func (m *MockTransactionValidator) ValidateContract(ctx context.Context, params *types.ETHTransactionParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateContract", ctx, params)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateContract indicates an expected call of ValidateContract
func (mr *MockTransactionValidatorMockRecorder) ValidateContract(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateContract", reflect.TypeOf((*MockTransactionValidator)(nil).ValidateContract), ctx, params)
}
