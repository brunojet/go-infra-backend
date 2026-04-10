package services

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/brunojet/go-infra-backend/debugassert"

	"github.com/brunojet/go-infra-backend/internal/utils"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type AuditedEntity struct {
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TestModel struct {
	AuditedEntity
	ID   int64 `gorm:"primaryKey,autoIncrement"`
	Name sql.NullString
}

func (m TestModel) TableName() string {
	return "test_models"
}

func (m TestModel) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}

type auditDto struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func ptrTimeUTC(t time.Time) *time.Time {
	tt := t.UTC()
	return &tt
}

func (d *auditDto) ToDTOPtr(e *AuditedEntity) {
	d.CreatedAt = ptrTimeUTC(e.CreatedAt)
	d.UpdatedAt = ptrTimeUTC(e.UpdatedAt)
	if e.DeletedAt.Valid {
		d.DeletedAt = ptrTimeUTC(e.DeletedAt.Time)
	} else {
		d.DeletedAt = nil
	}
}

type TestDTO struct {
	auditDto
	ID   string
	Name string
}

type TestMapper struct{}

// Implementa o método exigido pela interface ServiceMapper
func (m TestMapper) ToPostModel(dto TestDTO, model *TestModel) error {
	return TestMapper{}.ToModel(&dto, model)
}

func (m TestMapper) ToPatchModel(dto TestDTO, model *TestModel) error {
	return TestMapper{}.ToModel(&dto, model)
}

func (m TestMapper) GetModelKey(id string) (map[string]any, error) {
	matched, _ := regexp.MatchString(`^\d+$`, id)
	if !matched {
		return nil, errors.New("invalid id")
	}
	return map[string]any{"id": id}, nil
}

func (m TestMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

type testMapperQueryScopesError struct{}

func (m testMapperQueryScopesError) GetModelKey(id string) (map[string]any, error) {
	return TestMapper{}.GetModelKey(id)
}

func (m testMapperQueryScopesError) ToPostModel(dto TestDTO, model *TestModel) error {
	return TestMapper{}.ToModel(&dto, model)
}
func (m testMapperQueryScopesError) ToModel(dto *TestDTO, model *TestModel) error {
	return TestMapper{}.ToModel(dto, model)
}

func (m testMapperQueryScopesError) ToDTO(model *TestModel, dto *TestDTO) error {
	return TestMapper{}.ToDTO(model, dto)
}

func (m testMapperQueryScopesError) ToPatchModel(dto TestDTO, model *TestModel) error {
	return TestMapper{}.ToModel(&dto, model)
}

func (m testMapperQueryScopesError) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return nil, errors.New("query scope mapping error")
}

func (m TestMapper) ToModel(dto *TestDTO, model *TestModel) error {
	debugassert.Assert(dto != nil, "ToModel: dto is nil")
	debugassert.Assert(model != nil, "ToModel: model is nil")
	model.Name = utils.ToNullString(dto.Name)
	return nil
}

func (m TestMapper) ToDTO(model *TestModel, dto *TestDTO) error {
	debugassert.Assert(model != nil, "ToDTO: model is nil")
	debugassert.Assert(dto != nil, "ToDTO: dto is nil")
	dto.ID = utils.Int64ToString(model.ID)
	dto.Name = utils.FromNullString(model.Name)
	dto.auditDto.ToDTOPtr(&model.AuditedEntity)
	return nil
}

func testCreateDTO(t *testing.T, svc contracts.Service[TestDTO, TestDTO, TestDTO], ctx context.Context, in TestDTO) TestDTO {
	var out TestDTO
	err := svc.Create(ctx, in, &out)
	assert.NoError(t, err)
	v, err := strconv.ParseInt(out.ID, 10, 64)
	assert.NoError(t, err)
	assert.Greater(t, v, int64(0))
	assert.NotNil(t, out.CreatedAt)
	assert.NotNil(t, out.UpdatedAt)
	assert.Equal(t, time.Now().Year(), out.UpdatedAt.Year())
	assert.Equal(t, *out.UpdatedAt, *out.CreatedAt)
	in.ID = out.ID
	in.CreatedAt = out.CreatedAt
	in.UpdatedAt = out.UpdatedAt
	assert.Equal(t, in, out)
	return out
}

// Helper para criar controller, repo, service e context
type testServiceDeps struct {
	ctrl *gomock.Controller
	repo *MockRepository[TestModel]
	svc  contracts.Service[TestDTO, TestDTO, TestDTO]
	ctx  context.Context
}

func newTestServiceDeps(t *testing.T) testServiceDeps {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository[TestModel](ctrl)
	svc := NewServiceImpl(repo, TestMapper{})
	ctx := context.Background()

	return testServiceDeps{ctrl, repo, svc, ctx}
}

func TestGenericService_Create(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	repo.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	ExpectCreateWithTx(repo, func(inOut *TestModel) {
		inOut.ID = 1
		now := time.Now().UTC()
		inOut.CreatedAt = now
		inOut.UpdatedAt = now
	})

	in := TestDTO{Name: "one"}
	testCreateDTO(t, svc, ctx, in)
}

func TestGenericService_GetByID(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	var createdModel TestModel
	ExpectCreateWithTx(repo, func(inOut *TestModel) {
		inOut.ID = 2
		now := time.Now().UTC()
		inOut.CreatedAt = now
		inOut.UpdatedAt = now
		createdModel = *inOut
	})

	in := TestDTO{Name: "two"}
	out := testCreateDTO(t, svc, ctx, in)

	ExpectGetByIDWithTx(repo, func(out *TestModel) {
		*out = createdModel
	})
	var got TestDTO
	err := svc.GetByID(ctx, out.ID, &got)
	assert.NoError(t, err)
	assert.Equal(t, out, got)
}

func TestGenericService_List(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	var inTestData = []TestDTO{
		{Name: "List Created 1"},
		{Name: "List Created 2"},
	}

	models := make([]TestModel, 0, len(inTestData))
	outTestData := make(map[string]any, len(inTestData))
	for i, in := range inTestData {
		id := int64(i + 1)
		now := time.Now().UTC()
		m := TestModel{ID: id, Name: utils.ToNullString(in.Name), AuditedEntity: AuditedEntity{CreatedAt: now, UpdatedAt: now}}
		models = append(models, m)
		var dto TestDTO
		TestMapper{}.ToDTO(&m, &dto)
		outTestData[dto.ID] = dto
	}

	ExpectListWithTx(repo, func(out *[]TestModel) error {
		*out = models
		return nil
	})
	outList := make([]TestDTO, len(inTestData))
	total, err := svc.List(ctx, contracts.ListParams{Page: 1, Size: len(inTestData), OrderBy: "id", Order: "asc"}, &outList)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(inTestData)), total)
	assert.Len(t, outList, len(inTestData))
	for _, dto := range outList {
		assert.Equal(t, outTestData[dto.ID], dto)
	}
}

func TestGenericService_Update(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	var created TestModel
	ExpectCreateWithTx(repo, func(inOut *TestModel) {
		inOut.ID = 10
		now := time.Now().UTC()
		inOut.CreatedAt = now
		inOut.UpdatedAt = now
		created = *inOut
	})

	in := TestDTO{Name: "Update_Created"}
	out := testCreateDTO(t, svc, ctx, in)

	ExpectUpdateWithTx(repo, func(inOut *TestModel) {
		// simulate update: keep ID and CreatedAt, change Name and UpdatedAt
		inOut.ID = created.ID
		inOut.CreatedAt = created.CreatedAt
		now := created.CreatedAt.Add(2 * time.Second)
		inOut.UpdatedAt = now
	})

	upd := TestDTO{Name: "Update_Updated"}
	var updOut TestDTO
	err := svc.Update(ctx, out.ID, upd, &updOut)
	assert.NoError(t, err)
	assert.Equal(t, out.ID, updOut.ID)
	assert.Greater(t, *updOut.UpdatedAt, *updOut.CreatedAt)
	assert.Equal(t, "Update_Updated", updOut.Name)
}

func TestGenericService_Delete(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	ExpectCreateWithTx(repo, func(inOut *TestModel) {
		inOut.ID = 20
		now := time.Now().UTC()
		inOut.CreatedAt = now
		inOut.UpdatedAt = now
	})

	in := TestDTO{Name: "Delete_ToBeDeleted"}
	out := testCreateDTO(t, svc, ctx, in)

	ExpectDeleteWithTx(repo, nil)

	err := svc.Delete(ctx, out.ID)
	assert.NoError(t, err)
}

// errRepo sempre retorna erro para cada operação — usado para testar caminhos de erro
type errRepo struct{}

func (e *errRepo) Create(ctx context.Context, inOut *TestModel) error {
	return errors.New("repo error")
}
func (e *errRepo) GetByID(ctx context.Context, id map[string]any) (TestModel, error) {
	return TestModel{}, errors.New("repo error")
}
func (e *errRepo) List(ctx context.Context, listParams rpocts.ListParams) ([]TestModel, int64, error) {
	return nil, 0, errors.New("repo error")
}
func (e *errRepo) Update(ctx context.Context, id map[string]any, inOut *TestModel) error {
	return errors.New("repo error")
}
func (e *errRepo) Delete(ctx context.Context, id map[string]any) error {
	return errors.New("repo error")
}
func (e *errRepo) GormDB() *gorm.DB { return nil }
func (e *errRepo) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestGenericService_Errors(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	// Sempre mocka WithTx para padronizar
	// Usa o helper para garantir o padrão
	got := TestDTO{Name: "x"}
	ExpectCreateWithTx(repo, nil).Return(errors.New("repo error"))
	ExpectListWithTx(repo, nil).Return(int64(0), errors.New("repo error"))
	var createOut TestDTO
	err := svc.Create(ctx, got, &createOut)
	assert.Error(t, err)

	// GetByID mapper error (non-numeric id)
	var getByIDOut TestDTO
	err = svc.GetByID(ctx, "nope", &getByIDOut)
	assert.Error(t, err)

	// List repo error
	listOut := make([]TestDTO, 10)
	_, err = svc.List(ctx, contracts.ListParams{Page: 1, Size: 10, OrderBy: "id", Order: "asc"}, &listOut)
	assert.Error(t, err)

	// Update: mapper GetModelKey error when id non-numeric
	var u TestDTO
	var updOut TestDTO
	err = svc.Update(ctx, "nope", u, &updOut)
	assert.Error(t, err)

	// Update: repo error when id numeric
	repo.EXPECT().Update(gomock.Any(), gomock.Any(), &TestModel{}).Return(errors.New("repo error"))
	err = svc.Update(ctx, "1", u, &updOut)
	assert.Error(t, err)

	// Delete: mapper error for non-numeric id
	assert.Error(t, svc.Delete(ctx, "nope"))

	// Delete: repo error for numeric id
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("repo error"))
	assert.Error(t, svc.Delete(ctx, "1"))
}

func TestGenericService_List_MapperApplyQueryScopesError(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()

	repo := deps.repo
	// Sempre mocka WithTx para padronizar
	repo.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	).AnyTimes()
	svc := NewServiceImpl(repo, testMapperQueryScopesError{})
	ctx := deps.ctx

	responses := make([]TestDTO, 10)
	_, err := svc.List(ctx, contracts.ListParams{Page: 1, Size: 10, OrderBy: "id", Order: "asc"}, &responses)
	assert.Error(t, err)
}

func TestGenericService_GetByID_RepoError(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	call := ExpectGetByIDWithTx(repo, nil)
	call.Return(errors.New("repo error"))

	response := TestDTO{}
	err := svc.GetByID(ctx, "1", &response)
	assert.Error(t, err)
}

func TestGenericService_Update_RepoError(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	// Usa o helper para garantir o padrão de mock de transação
	call := ExpectUpdateWithTx(repo, nil)
	call.Return(errors.New("repo error"))

	request := TestDTO{Name: "update repo error"}
	response := TestDTO{}
	err := svc.Update(ctx, "1", request, &response)
	assert.Error(t, err)
}

func TestGenericService_Delete_RepoError(t *testing.T) {
	deps := newTestServiceDeps(t)
	defer deps.ctrl.Finish()
	repo := deps.repo
	svc := deps.svc
	ctx := deps.ctx

	call := ExpectDeleteWithTx(repo, nil)
	call.Return(errors.New("repo error"))

	err := svc.Delete(ctx, "1")
	assert.Error(t, err)
}
