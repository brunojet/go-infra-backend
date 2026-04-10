package handlers

import (
	context "context"
	reflect "reflect"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"

	gomock "github.com/golang/mock/gomock"
)

// ensure the generated mock (for concrete test types) implements the service contract
var _ svcContracts.Service[any, any, any] = (*MockService[any, any, any, repoContracts.Entity])(nil)

var _ svcContracts.ServiceMapper[any, any, any, repoContracts.Entity] = (*MockServiceMapper[any, any, any, repoContracts.Entity])(nil)

// MockServiceMapper is a mock of ServiceMapper interface.
type MockServiceMapper[C, R, U any, E repoContracts.Entity] struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMapperMockRecorder[C, R, U, E]
}

// MockServiceMapperMockRecorder is the mock recorder for MockServiceMapper.
type MockServiceMapperMockRecorder[C, R, U any, E repoContracts.Entity] struct {
	mock *MockServiceMapper[C, R, U, E]
}

// NewMockServiceMapper creates a new mock instance.
func NewMockServiceMapper[C, R, U any, E repoContracts.Entity](ctrl *gomock.Controller) *MockServiceMapper[C, R, U, E] {
	mock := &MockServiceMapper[C, R, U, E]{ctrl: ctrl}
	mock.recorder = &MockServiceMapperMockRecorder[C, R, U, E]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceMapper[C, R, U, E]) EXPECT() *MockServiceMapperMockRecorder[C, R, U, E] {
	return m.recorder
}

// GetModelKey mocks base method.
func (m *MockServiceMapper[C, R, U, E]) GetModelKey(id string) (map[string]any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModelKey", id)
	ret0, _ := ret[0].(map[string]any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplyQueryScopes mocks base method.
func (m *MockServiceMapper[C, R, U, E]) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplyQueryScopes", queryScopes)
	ret0, _ := ret[0].(map[string]any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplyQueryScopes indicates an expected call of ApplyQueryScopes.
func (mr *MockServiceMapperMockRecorder[C, R, U, E]) ApplyQueryScopes(queryScopes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyQueryScopes", reflect.TypeOf((*MockServiceMapper[C, R, U, E])(nil).ApplyQueryScopes), queryScopes)
}

// GetModelKey indicates an expected call of GetModelKey.
func (mr *MockServiceMapperMockRecorder[C, R, U, E]) GetModelKey(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModelKey", reflect.TypeOf((*MockServiceMapper[C, R, U, E])(nil).GetModelKey), id)
}

// ToDTO mocks base method.
func (m *MockServiceMapper[C, R, U, E]) ToDTO(model *E, dto *R) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToDTO", model, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// ToDTO indicates an expected call of ToDTO.
func (mr *MockServiceMapperMockRecorder[C, R, U, E]) ToDTO(model, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToDTO", reflect.TypeOf((*MockServiceMapper[C, R, U, E])(nil).ToDTO), model, dto)
}

// ToPostModel mocks base method.
func (m *MockServiceMapper[C, R, U, E]) ToPostModel(dto C, model *E) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToPostModel", dto, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// ToPostModel indicates an expected call of ToPostModel.
func (mr *MockServiceMapperMockRecorder[C, R, U, E]) ToPostModel(dto, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToPostModel", reflect.TypeOf((*MockServiceMapper[C, R, U, E])(nil).ToPostModel), dto, model)
}

// ToPatchModel mocks base method.
func (m *MockServiceMapper[C, R, U, E]) ToPatchModel(dto U, model *E) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToPatchModel", dto, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// ToPatchModel indicates an expected call of ToPatchModel.
func (mr *MockServiceMapperMockRecorder[C, R, U, E]) ToPatchModel(dto, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToPatchModel", reflect.TypeOf((*MockServiceMapper[C, R, U, E])(nil).ToPatchModel), dto, model)
}

// MockService is a mock of Service interface.
type MockService[C, R, U any, E repoContracts.Entity] struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder[C, R, U, E]
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder[C, R, U any, E repoContracts.Entity] struct {
	mock *MockService[C, R, U, E]
}

// NewMockService creates a new mock instance.
func NewMockService[C, R, U any, E repoContracts.Entity](ctrl *gomock.Controller) *MockService[C, R, U, E] {
	mock := &MockService[C, R, U, E]{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder[C, R, U, E]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService[C, R, U, E]) EXPECT() *MockServiceMockRecorder[C, R, U, E] {
	return m.recorder
}

// Create mocks base method (corrigido para interface atual).
func (m *MockService[C, R, U, E]) Create(ctx context.Context, dto C, response *R) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, dto, response)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockServiceMockRecorder[C, R, U, E]) Create(ctx, dto, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService[C, R, U, E])(nil).Create), ctx, dto, response)
}

// Delete mocks base method.
func (m *MockService[C, R, U, E]) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder[C, R, U, E]) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService[C, R, U, E])(nil).Delete), ctx, id)
}

// GetByID mocks base method (corrigido para interface atual).
func (m *MockService[C, R, U, E]) GetByID(ctx context.Context, id string, response *R) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id, response)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetByID indicates an expected call of GetByID.
func (mr *MockServiceMockRecorder[C, R, U, E]) GetByID(ctx, id, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockService[C, R, U, E])(nil).GetByID), ctx, id, response)
}

// List mocks base method (corrigido para interface atual).
func (m *MockService[C, R, U, E]) List(ctx context.Context, params svcContracts.ListParams, response *[]C) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, params, response)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder[C, R, U, E]) List(ctx, params, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService[C, R, U, E])(nil).List), ctx, params, response)
}

// Update mocks base method (corrigido para interface atual).
func (m *MockService[C, R, U, E]) Update(ctx context.Context, id string, dto U, response *R) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, dto, response)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockServiceMockRecorder[C, R, U, E]) Update(ctx, id, dto, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockService[C, R, U, E])(nil).Update), ctx, id, dto, response)
}
