// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	healthcheck "github.com/heptiolabs/healthcheck"
	io_prometheus_client "github.com/prometheus/client_model/go"
	api "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	entities "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	reflect "reflect"
)

// MockOrchestrateClient is a mock of OrchestrateClient interface
type MockOrchestrateClient struct {
	ctrl     *gomock.Controller
	recorder *MockOrchestrateClientMockRecorder
}

// MockOrchestrateClientMockRecorder is the mock recorder for MockOrchestrateClient
type MockOrchestrateClientMockRecorder struct {
	mock *MockOrchestrateClient
}

// NewMockOrchestrateClient creates a new mock instance
func NewMockOrchestrateClient(ctrl *gomock.Controller) *MockOrchestrateClient {
	mock := &MockOrchestrateClient{ctrl: ctrl}
	mock.recorder = &MockOrchestrateClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOrchestrateClient) EXPECT() *MockOrchestrateClientMockRecorder {
	return m.recorder
}

// SendContractTransaction mocks base method
func (m *MockOrchestrateClient) SendContractTransaction(ctx context.Context, request *api.SendTransactionRequest) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendContractTransaction", ctx, request)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendContractTransaction indicates an expected call of SendContractTransaction
func (mr *MockOrchestrateClientMockRecorder) SendContractTransaction(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendContractTransaction", reflect.TypeOf((*MockOrchestrateClient)(nil).SendContractTransaction), ctx, request)
}

// SendDeployTransaction mocks base method
func (m *MockOrchestrateClient) SendDeployTransaction(ctx context.Context, request *api.DeployContractRequest) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendDeployTransaction", ctx, request)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendDeployTransaction indicates an expected call of SendDeployTransaction
func (mr *MockOrchestrateClientMockRecorder) SendDeployTransaction(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendDeployTransaction", reflect.TypeOf((*MockOrchestrateClient)(nil).SendDeployTransaction), ctx, request)
}

// SendRawTransaction mocks base method
func (m *MockOrchestrateClient) SendRawTransaction(ctx context.Context, request *api.RawTransactionRequest) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRawTransaction", ctx, request)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendRawTransaction indicates an expected call of SendRawTransaction
func (mr *MockOrchestrateClientMockRecorder) SendRawTransaction(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRawTransaction", reflect.TypeOf((*MockOrchestrateClient)(nil).SendRawTransaction), ctx, request)
}

// SendTransferTransaction mocks base method
func (m *MockOrchestrateClient) SendTransferTransaction(ctx context.Context, request *api.TransferRequest) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTransferTransaction", ctx, request)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTransferTransaction indicates an expected call of SendTransferTransaction
func (mr *MockOrchestrateClientMockRecorder) SendTransferTransaction(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTransferTransaction", reflect.TypeOf((*MockOrchestrateClient)(nil).SendTransferTransaction), ctx, request)
}

// GetTxRequest mocks base method
func (m *MockOrchestrateClient) GetTxRequest(ctx context.Context, txRequestUUID string) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxRequest", ctx, txRequestUUID)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxRequest indicates an expected call of GetTxRequest
func (mr *MockOrchestrateClientMockRecorder) GetTxRequest(ctx, txRequestUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxRequest", reflect.TypeOf((*MockOrchestrateClient)(nil).GetTxRequest), ctx, txRequestUUID)
}

// GetSchedule mocks base method
func (m *MockOrchestrateClient) GetSchedule(ctx context.Context, scheduleUUID string) (*api.ScheduleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchedule", ctx, scheduleUUID)
	ret0, _ := ret[0].(*api.ScheduleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchedule indicates an expected call of GetSchedule
func (mr *MockOrchestrateClientMockRecorder) GetSchedule(ctx, scheduleUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchedule", reflect.TypeOf((*MockOrchestrateClient)(nil).GetSchedule), ctx, scheduleUUID)
}

// GetSchedules mocks base method
func (m *MockOrchestrateClient) GetSchedules(ctx context.Context) ([]*api.ScheduleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchedules", ctx)
	ret0, _ := ret[0].([]*api.ScheduleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchedules indicates an expected call of GetSchedules
func (mr *MockOrchestrateClientMockRecorder) GetSchedules(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchedules", reflect.TypeOf((*MockOrchestrateClient)(nil).GetSchedules), ctx)
}

// CreateSchedule mocks base method
func (m *MockOrchestrateClient) CreateSchedule(ctx context.Context, request *api.CreateScheduleRequest) (*api.ScheduleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSchedule", ctx, request)
	ret0, _ := ret[0].(*api.ScheduleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSchedule indicates an expected call of CreateSchedule
func (mr *MockOrchestrateClientMockRecorder) CreateSchedule(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSchedule", reflect.TypeOf((*MockOrchestrateClient)(nil).CreateSchedule), ctx, request)
}

// GetJob mocks base method
func (m *MockOrchestrateClient) GetJob(ctx context.Context, jobUUID string) (*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJob", ctx, jobUUID)
	ret0, _ := ret[0].(*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJob indicates an expected call of GetJob
func (mr *MockOrchestrateClientMockRecorder) GetJob(ctx, jobUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJob", reflect.TypeOf((*MockOrchestrateClient)(nil).GetJob), ctx, jobUUID)
}

// GetJobs mocks base method
func (m *MockOrchestrateClient) GetJobs(ctx context.Context) ([]*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobs", ctx)
	ret0, _ := ret[0].([]*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobs indicates an expected call of GetJobs
func (mr *MockOrchestrateClientMockRecorder) GetJobs(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobs", reflect.TypeOf((*MockOrchestrateClient)(nil).GetJobs), ctx)
}

// CreateJob mocks base method
func (m *MockOrchestrateClient) CreateJob(ctx context.Context, request *api.CreateJobRequest) (*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJob", ctx, request)
	ret0, _ := ret[0].(*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateJob indicates an expected call of CreateJob
func (mr *MockOrchestrateClientMockRecorder) CreateJob(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJob", reflect.TypeOf((*MockOrchestrateClient)(nil).CreateJob), ctx, request)
}

// UpdateJob mocks base method
func (m *MockOrchestrateClient) UpdateJob(ctx context.Context, jobUUID string, request *api.UpdateJobRequest) (*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateJob", ctx, jobUUID, request)
	ret0, _ := ret[0].(*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateJob indicates an expected call of UpdateJob
func (mr *MockOrchestrateClientMockRecorder) UpdateJob(ctx, jobUUID, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateJob", reflect.TypeOf((*MockOrchestrateClient)(nil).UpdateJob), ctx, jobUUID, request)
}

// StartJob mocks base method
func (m *MockOrchestrateClient) StartJob(ctx context.Context, jobUUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartJob", ctx, jobUUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartJob indicates an expected call of StartJob
func (mr *MockOrchestrateClientMockRecorder) StartJob(ctx, jobUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartJob", reflect.TypeOf((*MockOrchestrateClient)(nil).StartJob), ctx, jobUUID)
}

// ResendJobTx mocks base method
func (m *MockOrchestrateClient) ResendJobTx(ctx context.Context, jobUUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResendJobTx", ctx, jobUUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResendJobTx indicates an expected call of ResendJobTx
func (mr *MockOrchestrateClientMockRecorder) ResendJobTx(ctx, jobUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResendJobTx", reflect.TypeOf((*MockOrchestrateClient)(nil).ResendJobTx), ctx, jobUUID)
}

// SearchJob mocks base method
func (m *MockOrchestrateClient) SearchJob(ctx context.Context, filters *entities.JobFilters) ([]*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchJob", ctx, filters)
	ret0, _ := ret[0].([]*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchJob indicates an expected call of SearchJob
func (mr *MockOrchestrateClientMockRecorder) SearchJob(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchJob", reflect.TypeOf((*MockOrchestrateClient)(nil).SearchJob), ctx, filters)
}

// Checker mocks base method
func (m *MockOrchestrateClient) Checker() healthcheck.Check {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checker")
	ret0, _ := ret[0].(healthcheck.Check)
	return ret0
}

// Checker indicates an expected call of Checker
func (mr *MockOrchestrateClientMockRecorder) Checker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checker", reflect.TypeOf((*MockOrchestrateClient)(nil).Checker))
}

// Prometheus mocks base method
func (m *MockOrchestrateClient) Prometheus(arg0 context.Context) (map[string]*io_prometheus_client.MetricFamily, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prometheus", arg0)
	ret0, _ := ret[0].(map[string]*io_prometheus_client.MetricFamily)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prometheus indicates an expected call of Prometheus
func (mr *MockOrchestrateClientMockRecorder) Prometheus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prometheus", reflect.TypeOf((*MockOrchestrateClient)(nil).Prometheus), arg0)
}

// CreateAccount mocks base method
func (m *MockOrchestrateClient) CreateAccount(ctx context.Context, request *api.CreateAccountRequest) (*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, request)
	ret0, _ := ret[0].(*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount
func (mr *MockOrchestrateClientMockRecorder) CreateAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockOrchestrateClient)(nil).CreateAccount), ctx, request)
}

// SearchAccounts mocks base method
func (m *MockOrchestrateClient) SearchAccounts(ctx context.Context, filters *entities.AccountFilters) ([]*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAccounts", ctx, filters)
	ret0, _ := ret[0].([]*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchAccounts indicates an expected call of SearchAccounts
func (mr *MockOrchestrateClientMockRecorder) SearchAccounts(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAccounts", reflect.TypeOf((*MockOrchestrateClient)(nil).SearchAccounts), ctx, filters)
}

// GetAccount mocks base method
func (m *MockOrchestrateClient) GetAccount(ctx context.Context, address string) (*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, address)
	ret0, _ := ret[0].(*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockOrchestrateClientMockRecorder) GetAccount(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockOrchestrateClient)(nil).GetAccount), ctx, address)
}

// ImportAccount mocks base method
func (m *MockOrchestrateClient) ImportAccount(ctx context.Context, request *api.ImportAccountRequest) (*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportAccount", ctx, request)
	ret0, _ := ret[0].(*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportAccount indicates an expected call of ImportAccount
func (mr *MockOrchestrateClientMockRecorder) ImportAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportAccount", reflect.TypeOf((*MockOrchestrateClient)(nil).ImportAccount), ctx, request)
}

// UpdateAccount mocks base method
func (m *MockOrchestrateClient) UpdateAccount(ctx context.Context, address string, request *api.UpdateAccountRequest) (*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", ctx, address, request)
	ret0, _ := ret[0].(*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccount indicates an expected call of UpdateAccount
func (mr *MockOrchestrateClientMockRecorder) UpdateAccount(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockOrchestrateClient)(nil).UpdateAccount), ctx, address, request)
}

// SignPayload mocks base method
func (m *MockOrchestrateClient) SignPayload(ctx context.Context, address string, request *api.SignPayloadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignPayload", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignPayload indicates an expected call of SignPayload
func (mr *MockOrchestrateClientMockRecorder) SignPayload(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignPayload", reflect.TypeOf((*MockOrchestrateClient)(nil).SignPayload), ctx, address, request)
}

// SignTypedData mocks base method
func (m *MockOrchestrateClient) SignTypedData(ctx context.Context, address string, request *api.SignTypedDataRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTypedData", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTypedData indicates an expected call of SignTypedData
func (mr *MockOrchestrateClientMockRecorder) SignTypedData(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTypedData", reflect.TypeOf((*MockOrchestrateClient)(nil).SignTypedData), ctx, address, request)
}

// VerifySignature mocks base method
func (m *MockOrchestrateClient) VerifySignature(ctx context.Context, request *keymanager.VerifyPayloadRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySignature", ctx, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifySignature indicates an expected call of VerifySignature
func (mr *MockOrchestrateClientMockRecorder) VerifySignature(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySignature", reflect.TypeOf((*MockOrchestrateClient)(nil).VerifySignature), ctx, request)
}

// VerifyTypedDataSignature mocks base method
func (m *MockOrchestrateClient) VerifyTypedDataSignature(ctx context.Context, request *ethereum.VerifyTypedDataRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyTypedDataSignature", ctx, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyTypedDataSignature indicates an expected call of VerifyTypedDataSignature
func (mr *MockOrchestrateClientMockRecorder) VerifyTypedDataSignature(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyTypedDataSignature", reflect.TypeOf((*MockOrchestrateClient)(nil).VerifyTypedDataSignature), ctx, request)
}

// MockTransactionClient is a mock of TransactionClient interface
type MockTransactionClient struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionClientMockRecorder
}

// MockTransactionClientMockRecorder is the mock recorder for MockTransactionClient
type MockTransactionClientMockRecorder struct {
	mock *MockTransactionClient
}

// NewMockTransactionClient creates a new mock instance
func NewMockTransactionClient(ctrl *gomock.Controller) *MockTransactionClient {
	mock := &MockTransactionClient{ctrl: ctrl}
	mock.recorder = &MockTransactionClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTransactionClient) EXPECT() *MockTransactionClientMockRecorder {
	return m.recorder
}

// SendContractTransaction mocks base method
func (m *MockTransactionClient) SendContractTransaction(ctx context.Context, request *api.SendTransactionRequest) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendContractTransaction", ctx, request)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendContractTransaction indicates an expected call of SendContractTransaction
func (mr *MockTransactionClientMockRecorder) SendContractTransaction(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendContractTransaction", reflect.TypeOf((*MockTransactionClient)(nil).SendContractTransaction), ctx, request)
}

// SendDeployTransaction mocks base method
func (m *MockTransactionClient) SendDeployTransaction(ctx context.Context, request *api.DeployContractRequest) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendDeployTransaction", ctx, request)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendDeployTransaction indicates an expected call of SendDeployTransaction
func (mr *MockTransactionClientMockRecorder) SendDeployTransaction(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendDeployTransaction", reflect.TypeOf((*MockTransactionClient)(nil).SendDeployTransaction), ctx, request)
}

// SendRawTransaction mocks base method
func (m *MockTransactionClient) SendRawTransaction(ctx context.Context, request *api.RawTransactionRequest) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRawTransaction", ctx, request)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendRawTransaction indicates an expected call of SendRawTransaction
func (mr *MockTransactionClientMockRecorder) SendRawTransaction(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRawTransaction", reflect.TypeOf((*MockTransactionClient)(nil).SendRawTransaction), ctx, request)
}

// SendTransferTransaction mocks base method
func (m *MockTransactionClient) SendTransferTransaction(ctx context.Context, request *api.TransferRequest) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTransferTransaction", ctx, request)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTransferTransaction indicates an expected call of SendTransferTransaction
func (mr *MockTransactionClientMockRecorder) SendTransferTransaction(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTransferTransaction", reflect.TypeOf((*MockTransactionClient)(nil).SendTransferTransaction), ctx, request)
}

// GetTxRequest mocks base method
func (m *MockTransactionClient) GetTxRequest(ctx context.Context, txRequestUUID string) (*api.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxRequest", ctx, txRequestUUID)
	ret0, _ := ret[0].(*api.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxRequest indicates an expected call of GetTxRequest
func (mr *MockTransactionClientMockRecorder) GetTxRequest(ctx, txRequestUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxRequest", reflect.TypeOf((*MockTransactionClient)(nil).GetTxRequest), ctx, txRequestUUID)
}

// MockScheduleClient is a mock of ScheduleClient interface
type MockScheduleClient struct {
	ctrl     *gomock.Controller
	recorder *MockScheduleClientMockRecorder
}

// MockScheduleClientMockRecorder is the mock recorder for MockScheduleClient
type MockScheduleClientMockRecorder struct {
	mock *MockScheduleClient
}

// NewMockScheduleClient creates a new mock instance
func NewMockScheduleClient(ctrl *gomock.Controller) *MockScheduleClient {
	mock := &MockScheduleClient{ctrl: ctrl}
	mock.recorder = &MockScheduleClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScheduleClient) EXPECT() *MockScheduleClientMockRecorder {
	return m.recorder
}

// GetSchedule mocks base method
func (m *MockScheduleClient) GetSchedule(ctx context.Context, scheduleUUID string) (*api.ScheduleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchedule", ctx, scheduleUUID)
	ret0, _ := ret[0].(*api.ScheduleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchedule indicates an expected call of GetSchedule
func (mr *MockScheduleClientMockRecorder) GetSchedule(ctx, scheduleUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchedule", reflect.TypeOf((*MockScheduleClient)(nil).GetSchedule), ctx, scheduleUUID)
}

// GetSchedules mocks base method
func (m *MockScheduleClient) GetSchedules(ctx context.Context) ([]*api.ScheduleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchedules", ctx)
	ret0, _ := ret[0].([]*api.ScheduleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchedules indicates an expected call of GetSchedules
func (mr *MockScheduleClientMockRecorder) GetSchedules(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchedules", reflect.TypeOf((*MockScheduleClient)(nil).GetSchedules), ctx)
}

// CreateSchedule mocks base method
func (m *MockScheduleClient) CreateSchedule(ctx context.Context, request *api.CreateScheduleRequest) (*api.ScheduleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSchedule", ctx, request)
	ret0, _ := ret[0].(*api.ScheduleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSchedule indicates an expected call of CreateSchedule
func (mr *MockScheduleClientMockRecorder) CreateSchedule(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSchedule", reflect.TypeOf((*MockScheduleClient)(nil).CreateSchedule), ctx, request)
}

// MockJobClient is a mock of JobClient interface
type MockJobClient struct {
	ctrl     *gomock.Controller
	recorder *MockJobClientMockRecorder
}

// MockJobClientMockRecorder is the mock recorder for MockJobClient
type MockJobClientMockRecorder struct {
	mock *MockJobClient
}

// NewMockJobClient creates a new mock instance
func NewMockJobClient(ctrl *gomock.Controller) *MockJobClient {
	mock := &MockJobClient{ctrl: ctrl}
	mock.recorder = &MockJobClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJobClient) EXPECT() *MockJobClientMockRecorder {
	return m.recorder
}

// GetJob mocks base method
func (m *MockJobClient) GetJob(ctx context.Context, jobUUID string) (*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJob", ctx, jobUUID)
	ret0, _ := ret[0].(*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJob indicates an expected call of GetJob
func (mr *MockJobClientMockRecorder) GetJob(ctx, jobUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJob", reflect.TypeOf((*MockJobClient)(nil).GetJob), ctx, jobUUID)
}

// GetJobs mocks base method
func (m *MockJobClient) GetJobs(ctx context.Context) ([]*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobs", ctx)
	ret0, _ := ret[0].([]*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobs indicates an expected call of GetJobs
func (mr *MockJobClientMockRecorder) GetJobs(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobs", reflect.TypeOf((*MockJobClient)(nil).GetJobs), ctx)
}

// CreateJob mocks base method
func (m *MockJobClient) CreateJob(ctx context.Context, request *api.CreateJobRequest) (*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJob", ctx, request)
	ret0, _ := ret[0].(*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateJob indicates an expected call of CreateJob
func (mr *MockJobClientMockRecorder) CreateJob(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJob", reflect.TypeOf((*MockJobClient)(nil).CreateJob), ctx, request)
}

// UpdateJob mocks base method
func (m *MockJobClient) UpdateJob(ctx context.Context, jobUUID string, request *api.UpdateJobRequest) (*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateJob", ctx, jobUUID, request)
	ret0, _ := ret[0].(*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateJob indicates an expected call of UpdateJob
func (mr *MockJobClientMockRecorder) UpdateJob(ctx, jobUUID, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateJob", reflect.TypeOf((*MockJobClient)(nil).UpdateJob), ctx, jobUUID, request)
}

// StartJob mocks base method
func (m *MockJobClient) StartJob(ctx context.Context, jobUUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartJob", ctx, jobUUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartJob indicates an expected call of StartJob
func (mr *MockJobClientMockRecorder) StartJob(ctx, jobUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartJob", reflect.TypeOf((*MockJobClient)(nil).StartJob), ctx, jobUUID)
}

// ResendJobTx mocks base method
func (m *MockJobClient) ResendJobTx(ctx context.Context, jobUUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResendJobTx", ctx, jobUUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResendJobTx indicates an expected call of ResendJobTx
func (mr *MockJobClientMockRecorder) ResendJobTx(ctx, jobUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResendJobTx", reflect.TypeOf((*MockJobClient)(nil).ResendJobTx), ctx, jobUUID)
}

// SearchJob mocks base method
func (m *MockJobClient) SearchJob(ctx context.Context, filters *entities.JobFilters) ([]*api.JobResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchJob", ctx, filters)
	ret0, _ := ret[0].([]*api.JobResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchJob indicates an expected call of SearchJob
func (mr *MockJobClientMockRecorder) SearchJob(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchJob", reflect.TypeOf((*MockJobClient)(nil).SearchJob), ctx, filters)
}

// MockMetricClient is a mock of MetricClient interface
type MockMetricClient struct {
	ctrl     *gomock.Controller
	recorder *MockMetricClientMockRecorder
}

// MockMetricClientMockRecorder is the mock recorder for MockMetricClient
type MockMetricClientMockRecorder struct {
	mock *MockMetricClient
}

// NewMockMetricClient creates a new mock instance
func NewMockMetricClient(ctrl *gomock.Controller) *MockMetricClient {
	mock := &MockMetricClient{ctrl: ctrl}
	mock.recorder = &MockMetricClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMetricClient) EXPECT() *MockMetricClientMockRecorder {
	return m.recorder
}

// Checker mocks base method
func (m *MockMetricClient) Checker() healthcheck.Check {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checker")
	ret0, _ := ret[0].(healthcheck.Check)
	return ret0
}

// Checker indicates an expected call of Checker
func (mr *MockMetricClientMockRecorder) Checker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checker", reflect.TypeOf((*MockMetricClient)(nil).Checker))
}

// Prometheus mocks base method
func (m *MockMetricClient) Prometheus(arg0 context.Context) (map[string]*io_prometheus_client.MetricFamily, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prometheus", arg0)
	ret0, _ := ret[0].(map[string]*io_prometheus_client.MetricFamily)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prometheus indicates an expected call of Prometheus
func (mr *MockMetricClientMockRecorder) Prometheus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prometheus", reflect.TypeOf((*MockMetricClient)(nil).Prometheus), arg0)
}

// MockAccountClient is a mock of AccountClient interface
type MockAccountClient struct {
	ctrl     *gomock.Controller
	recorder *MockAccountClientMockRecorder
}

// MockAccountClientMockRecorder is the mock recorder for MockAccountClient
type MockAccountClientMockRecorder struct {
	mock *MockAccountClient
}

// NewMockAccountClient creates a new mock instance
func NewMockAccountClient(ctrl *gomock.Controller) *MockAccountClient {
	mock := &MockAccountClient{ctrl: ctrl}
	mock.recorder = &MockAccountClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountClient) EXPECT() *MockAccountClientMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method
func (m *MockAccountClient) CreateAccount(ctx context.Context, request *api.CreateAccountRequest) (*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, request)
	ret0, _ := ret[0].(*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount
func (mr *MockAccountClientMockRecorder) CreateAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountClient)(nil).CreateAccount), ctx, request)
}

// SearchAccounts mocks base method
func (m *MockAccountClient) SearchAccounts(ctx context.Context, filters *entities.AccountFilters) ([]*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAccounts", ctx, filters)
	ret0, _ := ret[0].([]*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchAccounts indicates an expected call of SearchAccounts
func (mr *MockAccountClientMockRecorder) SearchAccounts(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAccounts", reflect.TypeOf((*MockAccountClient)(nil).SearchAccounts), ctx, filters)
}

// GetAccount mocks base method
func (m *MockAccountClient) GetAccount(ctx context.Context, address string) (*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, address)
	ret0, _ := ret[0].(*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockAccountClientMockRecorder) GetAccount(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountClient)(nil).GetAccount), ctx, address)
}

// ImportAccount mocks base method
func (m *MockAccountClient) ImportAccount(ctx context.Context, request *api.ImportAccountRequest) (*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportAccount", ctx, request)
	ret0, _ := ret[0].(*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportAccount indicates an expected call of ImportAccount
func (mr *MockAccountClientMockRecorder) ImportAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportAccount", reflect.TypeOf((*MockAccountClient)(nil).ImportAccount), ctx, request)
}

// UpdateAccount mocks base method
func (m *MockAccountClient) UpdateAccount(ctx context.Context, address string, request *api.UpdateAccountRequest) (*api.AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", ctx, address, request)
	ret0, _ := ret[0].(*api.AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccount indicates an expected call of UpdateAccount
func (mr *MockAccountClientMockRecorder) UpdateAccount(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockAccountClient)(nil).UpdateAccount), ctx, address, request)
}

// SignPayload mocks base method
func (m *MockAccountClient) SignPayload(ctx context.Context, address string, request *api.SignPayloadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignPayload", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignPayload indicates an expected call of SignPayload
func (mr *MockAccountClientMockRecorder) SignPayload(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignPayload", reflect.TypeOf((*MockAccountClient)(nil).SignPayload), ctx, address, request)
}

// SignTypedData mocks base method
func (m *MockAccountClient) SignTypedData(ctx context.Context, address string, request *api.SignTypedDataRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTypedData", ctx, address, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTypedData indicates an expected call of SignTypedData
func (mr *MockAccountClientMockRecorder) SignTypedData(ctx, address, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTypedData", reflect.TypeOf((*MockAccountClient)(nil).SignTypedData), ctx, address, request)
}

// VerifySignature mocks base method
func (m *MockAccountClient) VerifySignature(ctx context.Context, request *keymanager.VerifyPayloadRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySignature", ctx, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifySignature indicates an expected call of VerifySignature
func (mr *MockAccountClientMockRecorder) VerifySignature(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySignature", reflect.TypeOf((*MockAccountClient)(nil).VerifySignature), ctx, request)
}

// VerifyTypedDataSignature mocks base method
func (m *MockAccountClient) VerifyTypedDataSignature(ctx context.Context, request *ethereum.VerifyTypedDataRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyTypedDataSignature", ctx, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyTypedDataSignature indicates an expected call of VerifyTypedDataSignature
func (mr *MockAccountClientMockRecorder) VerifyTypedDataSignature(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyTypedDataSignature", reflect.TypeOf((*MockAccountClient)(nil).VerifyTypedDataSignature), ctx, request)
}
