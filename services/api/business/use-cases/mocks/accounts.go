// Code generated by MockGen. DO NOT EDIT.
// Source: accounts.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	multitenancy "github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	entities "github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	common "github.com/ethereum/go-ethereum/common"
	hexutil "github.com/ethereum/go-ethereum/common/hexutil"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAccountUseCases is a mock of AccountUseCases interface
type MockAccountUseCases struct {
	ctrl     *gomock.Controller
	recorder *MockAccountUseCasesMockRecorder
}

// MockAccountUseCasesMockRecorder is the mock recorder for MockAccountUseCases
type MockAccountUseCasesMockRecorder struct {
	mock *MockAccountUseCases
}

// NewMockAccountUseCases creates a new mock instance
func NewMockAccountUseCases(ctrl *gomock.Controller) *MockAccountUseCases {
	mock := &MockAccountUseCases{ctrl: ctrl}
	mock.recorder = &MockAccountUseCasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountUseCases) EXPECT() *MockAccountUseCasesMockRecorder {
	return m.recorder
}

// GetAccount mocks base method
func (m *MockAccountUseCases) GetAccount() usecases.GetAccountUseCase {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount")
	ret0, _ := ret[0].(usecases.GetAccountUseCase)
	return ret0
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockAccountUseCasesMockRecorder) GetAccount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountUseCases)(nil).GetAccount))
}

// CreateAccount mocks base method
func (m *MockAccountUseCases) CreateAccount() usecases.CreateAccountUseCase {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount")
	ret0, _ := ret[0].(usecases.CreateAccountUseCase)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount
func (mr *MockAccountUseCasesMockRecorder) CreateAccount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountUseCases)(nil).CreateAccount))
}

// UpdateAccount mocks base method
func (m *MockAccountUseCases) UpdateAccount() usecases.UpdateAccountUseCase {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount")
	ret0, _ := ret[0].(usecases.UpdateAccountUseCase)
	return ret0
}

// UpdateAccount indicates an expected call of UpdateAccount
func (mr *MockAccountUseCasesMockRecorder) UpdateAccount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockAccountUseCases)(nil).UpdateAccount))
}

// SearchAccounts mocks base method
func (m *MockAccountUseCases) SearchAccounts() usecases.SearchAccountsUseCase {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAccounts")
	ret0, _ := ret[0].(usecases.SearchAccountsUseCase)
	return ret0
}

// SearchAccounts indicates an expected call of SearchAccounts
func (mr *MockAccountUseCasesMockRecorder) SearchAccounts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAccounts", reflect.TypeOf((*MockAccountUseCases)(nil).SearchAccounts))
}

// MockGetAccountUseCase is a mock of GetAccountUseCase interface
type MockGetAccountUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockGetAccountUseCaseMockRecorder
}

// MockGetAccountUseCaseMockRecorder is the mock recorder for MockGetAccountUseCase
type MockGetAccountUseCaseMockRecorder struct {
	mock *MockGetAccountUseCase
}

// NewMockGetAccountUseCase creates a new mock instance
func NewMockGetAccountUseCase(ctrl *gomock.Controller) *MockGetAccountUseCase {
	mock := &MockGetAccountUseCase{ctrl: ctrl}
	mock.recorder = &MockGetAccountUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGetAccountUseCase) EXPECT() *MockGetAccountUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockGetAccountUseCase) Execute(ctx context.Context, address common.Address, userInfo *multitenancy.UserInfo) (*entities.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, address, userInfo)
	ret0, _ := ret[0].(*entities.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockGetAccountUseCaseMockRecorder) Execute(ctx, address, userInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockGetAccountUseCase)(nil).Execute), ctx, address, userInfo)
}

// MockCreateAccountUseCase is a mock of CreateAccountUseCase interface
type MockCreateAccountUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCreateAccountUseCaseMockRecorder
}

// MockCreateAccountUseCaseMockRecorder is the mock recorder for MockCreateAccountUseCase
type MockCreateAccountUseCaseMockRecorder struct {
	mock *MockCreateAccountUseCase
}

// NewMockCreateAccountUseCase creates a new mock instance
func NewMockCreateAccountUseCase(ctrl *gomock.Controller) *MockCreateAccountUseCase {
	mock := &MockCreateAccountUseCase{ctrl: ctrl}
	mock.recorder = &MockCreateAccountUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCreateAccountUseCase) EXPECT() *MockCreateAccountUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockCreateAccountUseCase) Execute(ctx context.Context, identity *entities.Account, privateKey hexutil.Bytes, chainName string, userInfo *multitenancy.UserInfo) (*entities.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, identity, privateKey, chainName, userInfo)
	ret0, _ := ret[0].(*entities.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockCreateAccountUseCaseMockRecorder) Execute(ctx, identity, privateKey, chainName, userInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCreateAccountUseCase)(nil).Execute), ctx, identity, privateKey, chainName, userInfo)
}

// MockSearchAccountsUseCase is a mock of SearchAccountsUseCase interface
type MockSearchAccountsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSearchAccountsUseCaseMockRecorder
}

// MockSearchAccountsUseCaseMockRecorder is the mock recorder for MockSearchAccountsUseCase
type MockSearchAccountsUseCaseMockRecorder struct {
	mock *MockSearchAccountsUseCase
}

// NewMockSearchAccountsUseCase creates a new mock instance
func NewMockSearchAccountsUseCase(ctrl *gomock.Controller) *MockSearchAccountsUseCase {
	mock := &MockSearchAccountsUseCase{ctrl: ctrl}
	mock.recorder = &MockSearchAccountsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSearchAccountsUseCase) EXPECT() *MockSearchAccountsUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockSearchAccountsUseCase) Execute(ctx context.Context, filters *entities.AccountFilters, userInfo *multitenancy.UserInfo) ([]*entities.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, filters, userInfo)
	ret0, _ := ret[0].([]*entities.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockSearchAccountsUseCaseMockRecorder) Execute(ctx, filters, userInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockSearchAccountsUseCase)(nil).Execute), ctx, filters, userInfo)
}

// MockUpdateAccountUseCase is a mock of UpdateAccountUseCase interface
type MockUpdateAccountUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateAccountUseCaseMockRecorder
}

// MockUpdateAccountUseCaseMockRecorder is the mock recorder for MockUpdateAccountUseCase
type MockUpdateAccountUseCaseMockRecorder struct {
	mock *MockUpdateAccountUseCase
}

// NewMockUpdateAccountUseCase creates a new mock instance
func NewMockUpdateAccountUseCase(ctrl *gomock.Controller) *MockUpdateAccountUseCase {
	mock := &MockUpdateAccountUseCase{ctrl: ctrl}
	mock.recorder = &MockUpdateAccountUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUpdateAccountUseCase) EXPECT() *MockUpdateAccountUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockUpdateAccountUseCase) Execute(ctx context.Context, identity *entities.Account, userInfo *multitenancy.UserInfo) (*entities.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, identity, userInfo)
	ret0, _ := ret[0].(*entities.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockUpdateAccountUseCaseMockRecorder) Execute(ctx, identity, userInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockUpdateAccountUseCase)(nil).Execute), ctx, identity, userInfo)
}

// MockFundAccountUseCase is a mock of FundAccountUseCase interface
type MockFundAccountUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockFundAccountUseCaseMockRecorder
}

// MockFundAccountUseCaseMockRecorder is the mock recorder for MockFundAccountUseCase
type MockFundAccountUseCaseMockRecorder struct {
	mock *MockFundAccountUseCase
}

// NewMockFundAccountUseCase creates a new mock instance
func NewMockFundAccountUseCase(ctrl *gomock.Controller) *MockFundAccountUseCase {
	mock := &MockFundAccountUseCase{ctrl: ctrl}
	mock.recorder = &MockFundAccountUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFundAccountUseCase) EXPECT() *MockFundAccountUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockFundAccountUseCase) Execute(ctx context.Context, identity *entities.Account, chainName string, userInfo *multitenancy.UserInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, identity, chainName, userInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute
func (mr *MockFundAccountUseCaseMockRecorder) Execute(ctx, identity, chainName, userInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockFundAccountUseCase)(nil).Execute), ctx, identity, chainName, userInfo)
}
