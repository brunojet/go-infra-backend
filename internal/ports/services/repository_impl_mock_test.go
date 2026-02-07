package services

import (
	context "context"
	reflect "reflect"

	contracts "github.com/brunojet/go-infra-backend/internal/ports/repositories/contracts"
	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

var _ contracts.Repository[contracts.Entity] = (*MockRepository[contracts.Entity])(nil)

// MockEntity is a mock of Entity interface.
type MockEntity[E contracts.Entity] struct {
	ctrl     *gomock.Controller
	recorder *MockEntityMockRecorder[E]
}

// MockEntityMockRecorder is the mock recorder for MockEntity.
type MockEntityMockRecorder[E contracts.Entity] struct {
	mock *MockEntity[E]
}

// NewMockEntity creates a new mock instance.
func NewMockEntity[E contracts.Entity](ctrl *gomock.Controller) *MockEntity[E] {
	mock := &MockEntity[E]{ctrl: ctrl}
	mock.recorder = &MockEntityMockRecorder[E]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntity[E]) EXPECT() *MockEntityMockRecorder[E] {
	return m.recorder
}

// TableName mocks base method.
func (m *MockEntity[E]) TableName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableName")
	ret0, _ := ret[0].(string)
	return ret0
}

// TableName indicates an expected call of TableName.
func (mr *MockEntityMockRecorder[E]) TableName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableName", reflect.TypeOf((*MockEntity[E])(nil).TableName))
}

// MockRepository is a mock of Repository interface.
type MockRepository[E contracts.Entity] struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder[E]
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder[E contracts.Entity] struct {
	mock *MockRepository[E]
}

// NewMockRepository creates a new mock instance.
func NewMockRepository[E contracts.Entity](ctrl *gomock.Controller) *MockRepository[E] {
	mock := &MockRepository[E]{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder[E]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository[E]) EXPECT() *MockRepositoryMockRecorder[E] {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository[E]) Create(ctx context.Context, inOut *E) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, inOut)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder[E]) Create(ctx, inOut interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository[E])(nil).Create), ctx, inOut)
}

// GormDB mocks base method.
func (m *MockRepository[E]) GormDB() *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DB")
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// GormDB indicates an expected call of GormDB.
func (mr *MockRepositoryMockRecorder[E]) GormDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GormDB", reflect.TypeOf((*MockRepository[E])(nil).GormDB))
}

// Delete mocks base method.
func (m *MockRepository[E]) Delete(ctx context.Context, id map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder[E]) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository[E])(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockRepository[E]) GetByID(ctx context.Context, id map[string]any) (E, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(E)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder[E]) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository[E])(nil).GetByID), ctx, id)
}

// List mocks base method.
func (m *MockRepository[E]) List(ctx context.Context, listParams contracts.ListParams) ([]E, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, listParams)
	ret0, _ := ret[0].([]E)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockRepositoryMockRecorder[E]) List(ctx, listParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository[E])(nil).List), ctx, listParams)
}

// Update mocks base method.
func (m *MockRepository[E]) Update(ctx context.Context, id map[string]any, inOut *E) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, inOut)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder[E]) Update(ctx, id, inOut interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository[E])(nil).Update), ctx, id, inOut)
}

// WithTx mocks base method.
func (m *MockRepository[E]) WithTx(ctx context.Context, fn func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockRepositoryMockRecorder[E]) WithTx(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockRepository[E])(nil).WithTx), ctx, fn)
}
