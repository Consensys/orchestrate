// Code generated by MockGen. DO NOT EDIT.
// Source: use-cases.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entities "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	reflect "reflect"
)

// MockIdentityUseCases is a mock of IdentityUseCases interface
type MockIdentityUseCases struct {
	ctrl     *gomock.Controller
	recorder *MockIdentityUseCasesMockRecorder
}

// MockIdentityUseCasesMockRecorder is the mock recorder for MockIdentityUseCases
type MockIdentityUseCasesMockRecorder struct {
	mock *MockIdentityUseCases
}

// NewMockIdentityUseCases creates a new mock instance
func NewMockIdentityUseCases(ctrl *gomock.Controller) *MockIdentityUseCases {
	mock := &MockIdentityUseCases{ctrl: ctrl}
	mock.recorder = &MockIdentityUseCasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIdentityUseCases) EXPECT() *MockIdentityUseCasesMockRecorder {
	return m.recorder
}

// CreateIdentity mocks base method
func (m *MockIdentityUseCases) CreateIdentity() usecases.CreateIdentityUseCase {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIdentity")
	ret0, _ := ret[0].(usecases.CreateIdentityUseCase)
	return ret0
}

// CreateIdentity indicates an expected call of CreateIdentity
func (mr *MockIdentityUseCasesMockRecorder) CreateIdentity() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIdentity", reflect.TypeOf((*MockIdentityUseCases)(nil).CreateIdentity))
}

// SearchIdentity mocks base method
func (m *MockIdentityUseCases) SearchIdentity() usecases.SearchIdentitiesUseCase {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchIdentity")
	ret0, _ := ret[0].(usecases.SearchIdentitiesUseCase)
	return ret0
}

// SearchIdentity indicates an expected call of SearchIdentity
func (mr *MockIdentityUseCasesMockRecorder) SearchIdentity() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchIdentity", reflect.TypeOf((*MockIdentityUseCases)(nil).SearchIdentity))
}

// MockCreateIdentityUseCase is a mock of CreateIdentityUseCase interface
type MockCreateIdentityUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCreateIdentityUseCaseMockRecorder
}

// MockCreateIdentityUseCaseMockRecorder is the mock recorder for MockCreateIdentityUseCase
type MockCreateIdentityUseCaseMockRecorder struct {
	mock *MockCreateIdentityUseCase
}

// NewMockCreateIdentityUseCase creates a new mock instance
func NewMockCreateIdentityUseCase(ctrl *gomock.Controller) *MockCreateIdentityUseCase {
	mock := &MockCreateIdentityUseCase{ctrl: ctrl}
	mock.recorder = &MockCreateIdentityUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCreateIdentityUseCase) EXPECT() *MockCreateIdentityUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockCreateIdentityUseCase) Execute(ctx context.Context, identity *entities.Identity, privateKey, tenantID string) (*entities.Identity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, identity, privateKey, tenantID)
	ret0, _ := ret[0].(*entities.Identity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockCreateIdentityUseCaseMockRecorder) Execute(ctx, identity, privateKey, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCreateIdentityUseCase)(nil).Execute), ctx, identity, privateKey, tenantID)
}

// MockSearchIdentitiesUseCase is a mock of SearchIdentitiesUseCase interface
type MockSearchIdentitiesUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSearchIdentitiesUseCaseMockRecorder
}

// MockSearchIdentitiesUseCaseMockRecorder is the mock recorder for MockSearchIdentitiesUseCase
type MockSearchIdentitiesUseCaseMockRecorder struct {
	mock *MockSearchIdentitiesUseCase
}

// NewMockSearchIdentitiesUseCase creates a new mock instance
func NewMockSearchIdentitiesUseCase(ctrl *gomock.Controller) *MockSearchIdentitiesUseCase {
	mock := &MockSearchIdentitiesUseCase{ctrl: ctrl}
	mock.recorder = &MockSearchIdentitiesUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSearchIdentitiesUseCase) EXPECT() *MockSearchIdentitiesUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockSearchIdentitiesUseCase) Execute(ctx context.Context, filters *entities.IdentityFilters, tenants []string) ([]*entities.Identity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, filters, tenants)
	ret0, _ := ret[0].([]*entities.Identity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockSearchIdentitiesUseCaseMockRecorder) Execute(ctx, filters, tenants interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockSearchIdentitiesUseCase)(nil).Execute), ctx, filters, tenants)
}
