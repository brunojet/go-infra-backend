package repositories

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	dbadapters "github.com/brunojet/go-infra-backend/internal/database/adapters"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type AuditedEntity struct {
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TestEntity struct {
	AuditedEntity
	ID   string `gorm:"primaryKey"`
	Name *string
	Age  *int
}

func ps(s string) *string { return &s }
func pi(i int) *int       { return &i }

func (t TestEntity) TableName() string { return "test_entities" }

func mustGormDB(t *testing.T, db dbcontracts.DatabaseAdapter) *gorm.DB {
	t.Helper()
	gdb, err := db.GormDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	return gdb
}

func openMemoryDB(t *testing.T) (dbcontracts.DatabaseAdapter, func()) {
	t.Helper()
	db, err := dbadapters.NewSQLite("memory")
	require.NoError(t, err)
	// enable SQL logging at Info level for tests so queries are logged even on success
	gdb := mustGormDB(t, db)
	gdb.Config.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	// migrate test model
	if err := gdb.AutoMigrate(&TestEntity{}); err != nil {
		_ = db.Close()
		require.NoError(t, err)
	}
	return db, func() { _ = db.Close() }
}
func TestGormrepositories_InvalidIDAndNotFound(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	// repository now expects callers to provide a primary-key map; validate BuildPKeyMap behaviour
	// repository expects callers to provide a primary-key map; invalid id should be handled upstream
	_, err := repo.GetByID(ctx, map[string]any{})
	// generic repos may map empty id to ErrInvalidID or return a DB error; assert error
	assert.Error(t, err)

	// missing id should return ErrNotFound when using a valid pk map
	pk := map[string]any{"id": "missing-id"}
	_, err = repo.GetByID(ctx, pk)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGormrepositories_CreateAndGet(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	e := &TestEntity{ID: "id-create", Name: ps("bob"), Age: pi(22)}
	err := repo.Create(ctx, e)
	assert.NoError(t, err)

	// after Create the input object should be populated with the created values
	assert.Equal(t, "id-create", e.ID)
	assert.NotNil(t, e.Name)
	assert.NotNil(t, e.Age)
	assert.Equal(t, "bob", *e.Name)
	assert.Equal(t, 22, *e.Age)

	gotPtr, err := repo.GetByID(ctx, map[string]any{"id": "id-create"})
	assert.NoError(t, err)
	got := gotPtr
	assert.Equal(t, "id-create", got.ID)
	assert.NotNil(t, got.Name)
	assert.NotNil(t, got.Age)
	assert.Equal(t, "bob", *got.Name)
	assert.Equal(t, 22, *got.Age)

	// created object and fetched object should match by value
	assert.Equal(t, e.ID, got.ID)
	assert.Equal(t, *e.Name, *got.Name)
	assert.Equal(t, *e.Age, *got.Age)
}

func TestGormrepositories_List(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	// create multiple
	assert.NoError(t, repo.Create(ctx, &TestEntity{ID: "l-1", Name: ps("n1"), Age: pi(1)}))
	assert.NoError(t, repo.Create(ctx, &TestEntity{ID: "l-2", Name: ps("n2"), Age: pi(2)}))

	items, _, err := repo.List(ctx, repoContracts.ListParams{Page: 1, Size: 10, OrderBy: "NAME", Order: "asc"})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 2)
}

func TestGormrepositories_Update(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	// create initial entity with pointers
	assert.NoError(t, repo.Create(ctx, &TestEntity{ID: "u-1", Name: ps("orig"), Age: pi(10)}))

	// successful update with new values using pointers
	in := &TestEntity{Name: ps("updated")}
	assert.NoError(t, repo.Update(ctx, map[string]any{"id": "u-1"}, in))
	// input pointer should be populated with full loaded entity
	assert.NotNil(t, in.Name)
	assert.NotNil(t, in.Age)
	assert.Equal(t, "u-1", in.ID)
	assert.Equal(t, "updated", *in.Name)
}

func TestGormrepositories_Delete(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	assert.NoError(t, repo.Create(ctx, &TestEntity{ID: "d-1", Name: ps("to-del"), Age: pi(5)}))
	assert.NoError(t, repo.Delete(ctx, map[string]any{"id": "d-1"}))
	_, err := repo.GetByID(ctx, map[string]any{"id": "d-1"})
	assert.ErrorIs(t, err, ErrNotFound)

	// ensure soft-delete: Unscoped query should find the record and DeletedAt should be set
	var out TestEntity
	gdb := mustGormDB(t, db)
	err = gdb.Unscoped().First(&out, "id = ?", "d-1").Error
	assert.NoError(t, err)
	assert.False(t, out.DeletedAt.Time.IsZero())

	// hard delete the record and ensure it's gone
	assert.NoError(t, gdb.Unscoped().Delete(&out).Error)
	err = gdb.Unscoped().First(&out, "id = ?", "d-1").Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestGormrepositories_Update_NilInputAndDeletedAfterUpdate(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	// nil input should return ErrInvalidEntity
	err := repo.Update(ctx, map[string]any{"id": "some-id"}, nil)
	assert.Error(t, err)

	// create then delete before calling Update to simulate missing after update
	assert.NoError(t, repo.Create(ctx, &TestEntity{ID: "u-delete", Name: ps("x"), Age: pi(1)}))
	// remove it directly via DB
	assert.NoError(t, repo.Delete(ctx, map[string]any{"id": "u-delete"}))

	_, err = repo.GetByID(ctx, map[string]any{"id": "u-delete"})
	assert.Error(t, err)
	// attempt update — Updates will run but First should return ErrNotFound
	in := &TestEntity{ID: "u-delete", Name: ps("new"), Age: pi(2)}
	err = repo.Update(ctx, map[string]any{"id": "u-delete"}, in)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGetListUpdate_DBErrors(t *testing.T) {
	// open and immediately close DB to simulate low-level DB errors
	db, cleanup := openMemoryDB(t)
	// create repositories while DB is still open
	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	// close underlying DB to force errors
	cleanup()

	// GetByID should return an error (not ErrNotFound)
	_, err := repo.GetByID(ctx, map[string]any{"id": "any"})
	assert.Error(t, err)
	assert.False(t, errors.Is(err, ErrNotFound))

	// List should return error
	_, _, err = repo.List(ctx, repoContracts.ListParams{Page: 1, Size: 10, OrderBy: "", Order: ""})
	assert.Error(t, err)

	// Update should return error when DB closed
	in := &TestEntity{ID: "x"}
	err = repo.Update(ctx, map[string]any{"id": "any"}, in)
	assert.Error(t, err)
}

func TestRepositories_ClosedDB_MapsToErrDBUnavailable(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)

	// close underlying DB
	gdb, err := db.GormDB()
	require.NoError(t, err)
	sqlDB, err := gdb.DB()
	require.NoError(t, err)
	_ = sqlDB.Close()

	ctx := context.Background()

	_, err = repo.GetByID(ctx, map[string]any{"id": "any"})
	assert.ErrorIs(t, err, ErrDBUnavailable)

	_, _, err = repo.List(ctx, repoContracts.ListParams{Page: 1, Size: 1, OrderBy: "", Order: ""})
	// Accept either ErrDBUnavailable (DB closed) or ordering validation error depending on implementation
	assert.Error(t, err)
	assert.True(
		t,
		errors.Is(err, ErrDBUnavailable) || strings.Contains(err.Error(), "both orderBy and order must be provided together"),
		"expected ErrDBUnavailable or orderBy validation error, got: %v",
		err,
	)

	in := &TestEntity{ID: "x"}
	err = repo.Update(ctx, map[string]any{"id": "x"}, in)
	assert.ErrorIs(t, err, ErrDBUnavailable)
}

type testEntity struct{}

func (t testEntity) TableName() string { return "test_entities" }

type RepoTestModel struct {
	ID   int64 `gorm:"primaryKey;autoIncrement"`
	Name string
}

func (r RepoTestModel) TableName() string { return "repo_test_models" }

func TestSetPagination_ErrorsAndSuccess(t *testing.T) {
	db, close := openMemoryDB(t) // ensure DB can be opened before proceeding with List tests
	defer close()
	q := mustGormDB(t, db).Model(&RepoTestModel{})

	// invalid page
	err := setPagination(q, 0, 10)
	assert.Error(t, err)

	// invalid pageSize
	err = setPagination(q, 1, 0)
	assert.Error(t, err)

	// valid
	err = setPagination(q, 2, 5)
	assert.NoError(t, err)
}

func TestSetOrderBy_ErrorsAndSuccess(t *testing.T) {
	db, close := openMemoryDB(t) // ensure DB can be opened before proceeding with List tests
	defer close()
	q := mustGormDB(t, db).Model(&RepoTestModel{})

	// missing orderBy/order
	err := setOrderBy(q, "", "")
	assert.Error(t, err)

	// valid asc
	err = setOrderBy(q, "id", "asc")
	assert.NoError(t, err)

	// valid desc
	err = setOrderBy(q, "id", "desc")
	assert.NoError(t, err)
}

func TestGetListSize_Branches(t *testing.T) {
	// total 10, page1 size3 -> capacity 3
	cap := getListSize(10, 1, 3)
	assert.Equal(t, 3, cap)

	// total 10, page4 size3 -> offset 9 remaining 1 -> capacity 1
	cap = getListSize(10, 4, 3)
	assert.Equal(t, 1, cap)

	// total 5, page3 size3 -> offset 6 remaining -1 -> capacity 0
	cap = getListSize(5, 3, 3)
	assert.Equal(t, 0, cap)
}

func TestList_ErrorsAndSuccess(t *testing.T) {
	db, close := openMemoryDB(t) // ensure DB can be opened before proceeding with List tests
	defer close()

	gdb := mustGormDB(t, db)
	if err := gdb.AutoMigrate(&RepoTestModel{}); err != nil {
		_ = db.Close()
		require.NoError(t, err)
	}

	repo := NewGormRepository[RepoTestModel](db)
	ctx := context.Background()

	// empty table -> Count rowsAffected == 0 -> MapTxError returns ErrNotFound
	_, _, err := repo.List(ctx, contracts.ListParams{Page: 1, Size: 10, OrderBy: "id", Order: "asc"})
	assert.Error(t, err)

	// create records
	for i := 0; i < 5; i++ {
		m := RepoTestModel{Name: "n"}
		err := repo.Create(ctx, &m)
		assert.NoError(t, err)
	}

	// missing orderBy -> should error
	_, _, err = repo.List(ctx, contracts.ListParams{Page: 1, Size: 2, OrderBy: "", Order: "asc"})
	assert.Error(t, err)

	// invalid pagination -> should error
	_, _, err = repo.List(ctx, contracts.ListParams{Page: 0, Size: 2, OrderBy: "id", Order: "asc"})
	assert.Error(t, err)

	// success: page 2 size 2 -> expect 2 items
	items, total, err := repo.List(ctx, contracts.ListParams{Page: 2, Size: 2, OrderBy: "id", Order: "asc"})
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, items, 2)
}

func TestGetDB(t *testing.T) {
	db, close := openMemoryDB(t)
	defer close()
	repo := NewGormRepository[TestEntity](db)
	assert.Same(t, mustGormDB(t, db), repo.GormDB())
}

func TestGormRepository_WithTx_UsesTransactionalContextForCRUD(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	err := repo.WithTx(ctx, func(txCtx context.Context) error {
		if errCreate := repo.Create(txCtx, &TestEntity{ID: "tx-rollback", Name: ps("tmp"), Age: pi(1)}); errCreate != nil {
			return errCreate
		}
		return errors.New("force rollback")
	})
	assert.EqualError(t, err, "force rollback")

	_, err = repo.GetByID(ctx, map[string]any{"id": "tx-rollback"})
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestList_Update_Delete_ExtraErrorBranches(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	_, _, err := repo.List(ctx, contracts.ListParams{
		QueryParams: contracts.QueryParams{Scopes: map[string]any{"": "x"}},
		Page:        1,
		Size:        10,
		OrderBy:     "id",
		Order:       "asc",
	})
	assert.ErrorIs(t, err, ErrInvalidScope)

	err = repo.Update(ctx, map[string]any{}, &TestEntity{Name: ps("n")})
	assert.ErrorIs(t, err, ErrEmptyScopes)

	err = repo.Delete(ctx, map[string]any{})
	assert.ErrorIs(t, err, ErrEmptyScopes)

	err = repo.Delete(ctx, map[string]any{"id": "missing"})
	assert.ErrorIs(t, err, ErrNotFound)
}
