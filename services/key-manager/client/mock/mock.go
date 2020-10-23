// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	healthcheck "github.com/heptiolabs/healthcheck"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	reflect "reflect"
)

// MockEthereumAccountClient is a mock of EthereumAccountClient interface
type MockEthereumAccountClient struct {
	ctrl     *gomock.Controller
	recorder *MockEthereumAccountClientMockRecorder
}

// MockEthereumAccountClientMockRecorder is the mock recorder for MockEthereumAccountClient
type MockEthereumAccountClientMockRecorder struct {
	mock *MockEthereumAccountClient
}

// NewMockEthereumAccountClient creates a new mock instance
func NewMockEthereumAccountClient(ctrl *gomock.Controller) *MockEthereumAccountClient {
	mock := &MockEthereumAccountClient{ctrl: ctrl}
	mock.recorder = &MockEthereumAccountClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEthereumAccountClient) EXPECT() *MockEthereumAccountClientMockRecorder {
	return m.recorder
}

// ETHCreateAccount mocks base method
func (m *MockEthereumAccountClient) ETHCreateAccount(ctx context.Context, request *ethereum.CreateETHAccountRequest) (*ethereum.ETHAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHCreateAccount", ctx, request)
	ret0, _ := ret[0].(*ethereum.ETHAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHCreateAccount indicates an expected call of ETHCreateAccount
func (mr *MockEthereumAccountClientMockRecorder) ETHCreateAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHCreateAccount", reflect.TypeOf((*MockEthereumAccountClient)(nil).ETHCreateAccount), ctx, request)
}

// ETHImportAccount mocks base method
func (m *MockEthereumAccountClient) ETHImportAccount(ctx context.Context, request *ethereum.ImportETHAccountRequest) (*ethereum.ETHAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHImportAccount", ctx, request)
	ret0, _ := ret[0].(*ethereum.ETHAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHImportAccount indicates an expected call of ETHImportAccount
func (mr *MockEthereumAccountClientMockRecorder) ETHImportAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHImportAccount", reflect.TypeOf((*MockEthereumAccountClient)(nil).ETHImportAccount), ctx, request)
}

// ETHSign mocks base method
func (m *MockEthereumAccountClient) ETHSign(ctx context.Context, address string, request *keymanager.PayloadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHSign", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHSign indicates an expected call of ETHSign
func (mr *MockEthereumAccountClientMockRecorder) ETHSign(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHSign", reflect.TypeOf((*MockEthereumAccountClient)(nil).ETHSign), ctx, address, request)
}

// ETHSignTransaction mocks base method
func (m *MockEthereumAccountClient) ETHSignTransaction(ctx context.Context, address string, request *ethereum.SignETHTransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHSignTransaction", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHSignTransaction indicates an expected call of ETHSignTransaction
func (mr *MockEthereumAccountClientMockRecorder) ETHSignTransaction(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHSignTransaction", reflect.TypeOf((*MockEthereumAccountClient)(nil).ETHSignTransaction), ctx, address, request)
}

// ETHSignQuorumPrivateTransaction mocks base method
func (m *MockEthereumAccountClient) ETHSignQuorumPrivateTransaction(ctx context.Context, address string, request *ethereum.SignQuorumPrivateTransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHSignQuorumPrivateTransaction", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHSignQuorumPrivateTransaction indicates an expected call of ETHSignQuorumPrivateTransaction
func (mr *MockEthereumAccountClientMockRecorder) ETHSignQuorumPrivateTransaction(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHSignQuorumPrivateTransaction", reflect.TypeOf((*MockEthereumAccountClient)(nil).ETHSignQuorumPrivateTransaction), ctx, address, request)
}

// ETHSignEEATransaction mocks base method
func (m *MockEthereumAccountClient) ETHSignEEATransaction(ctx context.Context, address string, request *ethereum.SignEEATransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHSignEEATransaction", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHSignEEATransaction indicates an expected call of ETHSignEEATransaction
func (mr *MockEthereumAccountClientMockRecorder) ETHSignEEATransaction(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHSignEEATransaction", reflect.TypeOf((*MockEthereumAccountClient)(nil).ETHSignEEATransaction), ctx, address, request)
}

// MockKeyManagerClient is a mock of KeyManagerClient interface
type MockKeyManagerClient struct {
	ctrl     *gomock.Controller
	recorder *MockKeyManagerClientMockRecorder
}

// MockKeyManagerClientMockRecorder is the mock recorder for MockKeyManagerClient
type MockKeyManagerClientMockRecorder struct {
	mock *MockKeyManagerClient
}

// NewMockKeyManagerClient creates a new mock instance
func NewMockKeyManagerClient(ctrl *gomock.Controller) *MockKeyManagerClient {
	mock := &MockKeyManagerClient{ctrl: ctrl}
	mock.recorder = &MockKeyManagerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeyManagerClient) EXPECT() *MockKeyManagerClientMockRecorder {
	return m.recorder
}

// Checker mocks base method
func (m *MockKeyManagerClient) Checker() healthcheck.Check {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checker")
	ret0, _ := ret[0].(healthcheck.Check)
	return ret0
}

// Checker indicates an expected call of Checker
func (mr *MockKeyManagerClientMockRecorder) Checker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checker", reflect.TypeOf((*MockKeyManagerClient)(nil).Checker))
}

// ETHCreateAccount mocks base method
func (m *MockKeyManagerClient) ETHCreateAccount(ctx context.Context, request *ethereum.CreateETHAccountRequest) (*ethereum.ETHAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHCreateAccount", ctx, request)
	ret0, _ := ret[0].(*ethereum.ETHAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHCreateAccount indicates an expected call of ETHCreateAccount
func (mr *MockKeyManagerClientMockRecorder) ETHCreateAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHCreateAccount", reflect.TypeOf((*MockKeyManagerClient)(nil).ETHCreateAccount), ctx, request)
}

// ETHImportAccount mocks base method
func (m *MockKeyManagerClient) ETHImportAccount(ctx context.Context, request *ethereum.ImportETHAccountRequest) (*ethereum.ETHAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHImportAccount", ctx, request)
	ret0, _ := ret[0].(*ethereum.ETHAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHImportAccount indicates an expected call of ETHImportAccount
func (mr *MockKeyManagerClientMockRecorder) ETHImportAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHImportAccount", reflect.TypeOf((*MockKeyManagerClient)(nil).ETHImportAccount), ctx, request)
}

// ETHSign mocks base method
func (m *MockKeyManagerClient) ETHSign(ctx context.Context, address string, request *keymanager.PayloadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHSign", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHSign indicates an expected call of ETHSign
func (mr *MockKeyManagerClientMockRecorder) ETHSign(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHSign", reflect.TypeOf((*MockKeyManagerClient)(nil).ETHSign), ctx, address, request)
}

// ETHSignTransaction mocks base method
func (m *MockKeyManagerClient) ETHSignTransaction(ctx context.Context, address string, request *ethereum.SignETHTransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHSignTransaction", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHSignTransaction indicates an expected call of ETHSignTransaction
func (mr *MockKeyManagerClientMockRecorder) ETHSignTransaction(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHSignTransaction", reflect.TypeOf((*MockKeyManagerClient)(nil).ETHSignTransaction), ctx, address, request)
}

// ETHSignQuorumPrivateTransaction mocks base method
func (m *MockKeyManagerClient) ETHSignQuorumPrivateTransaction(ctx context.Context, address string, request *ethereum.SignQuorumPrivateTransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHSignQuorumPrivateTransaction", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHSignQuorumPrivateTransaction indicates an expected call of ETHSignQuorumPrivateTransaction
func (mr *MockKeyManagerClientMockRecorder) ETHSignQuorumPrivateTransaction(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHSignQuorumPrivateTransaction", reflect.TypeOf((*MockKeyManagerClient)(nil).ETHSignQuorumPrivateTransaction), ctx, address, request)
}

// ETHSignEEATransaction mocks base method
func (m *MockKeyManagerClient) ETHSignEEATransaction(ctx context.Context, address string, request *ethereum.SignEEATransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHSignEEATransaction", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ETHSignEEATransaction indicates an expected call of ETHSignEEATransaction
func (mr *MockKeyManagerClientMockRecorder) ETHSignEEATransaction(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHSignEEATransaction", reflect.TypeOf((*MockKeyManagerClient)(nil).ETHSignEEATransaction), ctx, address, request)
}
