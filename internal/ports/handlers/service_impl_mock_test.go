package handlers

import (
	context "context"
	reflect "reflect"

	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"

	gomock "github.com/golang/mock/gomock"
)

// ensure the generated mock (for concrete test types) implements the service contract
var _ svcContracts.Service[any, contracts.Entity] = (*MockService[any, contracts.Entity])(nil)

var _ svcContracts.ServiceMapper[any, contracts.Entity] = (*MockServiceMapper[any, contracts.Entity])(nil)

// MockServiceMapper is a mock of ServiceMapper interface.
type MockServiceMapper[D any, E contracts.Entity] struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMapperMockRecorder[D, E]
}

// MockServiceMapperMockRecorder is the mock recorder for MockServiceMapper.
type MockServiceMapperMockRecorder[D any, E contracts.Entity] struct {
	mock *MockServiceMapper[D, E]
}

// NewMockServiceMapper creates a new mock instance.
func NewMockServiceMapper[D any, E contracts.Entity](ctrl *gomock.Controller) *MockServiceMapper[D, E] {
	mock := &MockServiceMapper[D, E]{ctrl: ctrl}
	mock.recorder = &MockServiceMapperMockRecorder[D, E]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceMapper[D, E]) EXPECT() *MockServiceMapperMockRecorder[D, E] {
	return m.recorder
}

// GetModelKey mocks base method.
func (m *MockServiceMapper[D, E]) GetModelKey(id string) (map[string]any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModelKey", id)
	ret0, _ := ret[0].(map[string]any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetModelKey indicates an expected call of GetModelKey.
func (mr *MockServiceMapperMockRecorder[D, E]) GetModelKey(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModelKey", reflect.TypeOf((*MockServiceMapper[D, E])(nil).GetModelKey), id)
}

// ToDTO mocks base method.
func (m *MockServiceMapper[D, E]) ToDTO(model *E, dto *D) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToDTO", model, dto)
}

// ToDTO indicates an expected call of ToDTO.
func (mr *MockServiceMapperMockRecorder[D, E]) ToDTO(model, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToDTO", reflect.TypeOf((*MockServiceMapper[D, E])(nil).ToDTO), model, dto)
}

// ToModel mocks base method.
func (m *MockServiceMapper[D, E]) ToModel(dto *D, model *E) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToModel", dto, model)
}

// ToModel indicates an expected call of ToModel.
func (mr *MockServiceMapperMockRecorder[D, E]) ToModel(dto, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToModel", reflect.TypeOf((*MockServiceMapper[D, E])(nil).ToModel), dto, model)
}

// MockService is a mock of Service interface.
type MockService[D any, E contracts.Entity] struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder[D, E]
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder[D any, E contracts.Entity] struct {
	mock *MockService[D, E]
}

// NewMockService creates a new mock instance.
func NewMockService[D any, E contracts.Entity](ctrl *gomock.Controller) *MockService[D, E] {
	mock := &MockService[D, E]{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder[D, E]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService[D, E]) EXPECT() *MockServiceMockRecorder[D, E] {
	return m.recorder
}

// Create mocks base method.
func (m *MockService[D, E]) Create(ctx context.Context, dto *D) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockServiceMockRecorder[D, E]) Create(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService[D, E])(nil).Create), ctx, dto)
}

// Delete mocks base method.
func (m *MockService[D, E]) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder[D, E]) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService[D, E])(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockService[D, E]) GetByID(ctx context.Context, id string) (D, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	var ret0 D
	if ret[0] != nil {
		ret0 = ret[0].(D)
	}
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockServiceMockRecorder[D, E]) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockService[D, E])(nil).GetByID), ctx, id)
}

// List mocks base method.
func (m *MockService[D, E]) List(ctx context.Context, size int) ([]D, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, size)
	var ret0 []D
	if ret[0] != nil {
		ret0 = ret[0].([]D)
	}
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder[D, E]) List(ctx, size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService[D, E])(nil).List), ctx, size)
}

// Update mocks base method.
func (m *MockService[D, E]) Update(ctx context.Context, id string, dto *D) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockServiceMockRecorder[D, E]) Update(ctx, id, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockService[D, E])(nil).Update), ctx, id, dto)
}
