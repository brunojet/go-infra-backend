package services

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/brunojet/go-infra-backend/internal/database"
	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	"github.com/brunojet/go-infra-backend/internal/ports/repositories"
	repoContracts "github.com/brunojet/go-infra-backend/internal/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/internal/ports/services/contracts"
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

func (m TestMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{"id": id}, nil
}

func (m TestMapper) ToModel(d *TestDTO) (TestModel, error) {
	if d == nil {
		return TestModel{}, nil
	}
	var mdl TestModel
	mdl.Name = ToNullString(d.Name)
	return mdl, nil
}

func (m TestMapper) ToDTO(mdl *TestModel, dto *TestDTO) {
	if dto == nil {
		return
	}
	if mdl == nil {
		*dto = TestDTO{}
		return
	}
	dto.ID = Int64ToString(mdl.ID)
	if mdl.Name.Valid {
		dto.Name = FromNullString(mdl.Name)
	}
	dto.auditDto.ToDTOPtr(&mdl.AuditedEntity)
}

func testCreateDTO(t *testing.T, svc contracts.Service[TestDTO, TestModel], ctx context.Context, in TestDTO) TestDTO {
	out := in
	err := svc.Create(ctx, &out)
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

func openMemoryDB(t *testing.T) (dbcontracts.Database, func()) {
	t.Helper()
	db, err := database.NewSQLiteDatabase("memory")
	assert.NoError(t, err)

	err = db.GormDB().AutoMigrate(&TestModel{})
	assert.NoError(t, err)

	return db, func() { _ = db.Close() }
}

func TestGenericService_Create(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()
	repo := repositories.NewGormRepository[TestModel](db.GormDB())
	svc := NewServiceImpl(repoContracts.Repository[TestModel](repo), TestMapper{})
	ctx := context.Background()

	in := TestDTO{Name: "one"}
	testCreateDTO(t, svc, ctx, in)
}

func TestGenericService_GetByID(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()
	repo := repositories.NewGormRepository[TestModel](db.GormDB())
	svc := NewServiceImpl(repoContracts.Repository[TestModel](repo), TestMapper{})
	ctx := context.Background()

	in := TestDTO{Name: "two"}
	out := testCreateDTO(t, svc, ctx, in)

	got, err := svc.GetByID(ctx, out.ID)
	assert.NoError(t, err)
	assert.Equal(t, out, got)
}

func TestGenericService_List(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()
	repo := repositories.NewGormRepository[TestModel](db.GormDB())
	svc := NewServiceImpl(repoContracts.Repository[TestModel](repo), TestMapper{})
	ctx := context.Background()

	var inTestData = []TestDTO{
		{Name: "List Created 1"},
		{Name: "List Created 2"},
	}

	outTestData := make(map[string]any, len(inTestData))

	for _, in := range inTestData {
		out := testCreateDTO(t, svc, ctx, in)
		outTestData[out.ID] = out
	}

	list, err := svc.List(ctx, len(inTestData))
	assert.NoError(t, err)
	assert.Len(t, list, len(inTestData))
	for _, dto := range list {
		assert.Equal(t, outTestData[dto.ID], dto)
	}
}

func TestGenericService_Update(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()
	repo := repositories.NewGormRepository[TestModel](db.GormDB())
	svc := NewServiceImpl(repoContracts.Repository[TestModel](repo), TestMapper{})
	ctx := context.Background()

	in := TestDTO{Name: "Update_Created"}
	out := testCreateDTO(t, svc, ctx, in)
	time.Sleep(1 * time.Second)

	upd := TestDTO{Name: "Update_Updated"}
	err := svc.Update(ctx, out.ID, &upd)
	assert.NoError(t, err)
	assert.Equal(t, out.ID, upd.ID)
	assert.Greater(t, *upd.UpdatedAt, *upd.CreatedAt)
	assert.Equal(t, "Update_Updated", upd.Name)
}

func TestGenericService_Delete(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()
	repo := repositories.NewGormRepository[TestModel](db.GormDB())
	svc := NewServiceImpl(repoContracts.Repository[TestModel](repo), TestMapper{})
	ctx := context.Background()

	in := TestDTO{Name: "Delete_ToBeDeleted"}
	out := testCreateDTO(t, svc, ctx, in)

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
func (e *errRepo) List(ctx context.Context, listParams repoContracts.ListParams) ([]TestModel, int, error) {
	return nil, 0, errors.New("repo error")
}
func (e *errRepo) Update(ctx context.Context, id map[string]any, inOut *TestModel) error {
	return errors.New("repo error")
}
func (e *errRepo) Delete(ctx context.Context, id map[string]any) error {
	return errors.New("repo error")
}
func (e *errRepo) DB() *gorm.DB { return nil }
func (e *errRepo) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestGenericService_Errors(t *testing.T) {
	svc := NewServiceImpl(&errRepo{}, TestMapper{})
	ctx := context.Background()

	// Create error
	got := TestDTO{Name: "x"}
	err := svc.Create(ctx, &got)
	assert.Error(t, err)

	// GetByID error
	_, err = svc.GetByID(ctx, "nope")
	assert.Error(t, err)

	// List error
	_, err = svc.List(ctx, 0)
	assert.Error(t, err)

	// Update error
	var u TestDTO
	err = svc.Update(ctx, "nope", &u)
	assert.Error(t, err)

	// Delete error
	assert.Error(t, svc.Delete(ctx, "nope"))
}
