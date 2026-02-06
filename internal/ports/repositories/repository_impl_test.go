package repositories

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	dbadapters "github.com/brunojet/go-infra-backend/internal/database/adapters"
	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	repoContracts "github.com/brunojet/go-infra-backend/internal/ports/repositories/contracts"
	"github.com/stretchr/testify/assert"
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

func openMemoryDB(t *testing.T) (dbcontracts.Database, func()) {
	t.Helper()
	db, err := dbadapters.NewSQLite("memory")
	if err != nil {
		t.Fatalf("failed to open sqlite memory: %v", err)
	}
	// enable SQL logging at Info level for tests so queries are logged even on success
	if db != nil && db.GormDB() != nil {
		db.GormDB().Config.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}
	// migrate test model
	if err := db.GormDB().AutoMigrate(&TestEntity{}); err != nil {
		_ = db.Close()
		t.Fatalf("auto migrate failed: %v", err)
	}
	return db, func() { _ = db.Close() }
}
func TestGormrepositories_InvalidIDAndNotFound(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db.GormDB())
	ctx := context.Background()

	// repository now expects callers to provide a primary-key map; validate BuildPKeyMap behaviour
	// repository expects callers to provide a primary-key map; invalid id should be handled upstream
	_, err := repo.GetByID(ctx, map[string]any{})
	// generic repos may map empty id to ErrInvalidID or return a DB error; assert error
	if err == nil {
		t.Fatalf("expected error for empty id map")
	}

	// missing id should return ErrNotFound when using a valid pk map
	pk := map[string]any{"id": "missing-id"}
	_, err = repo.GetByID(ctx, pk)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGormrepositories_CreateAndGet(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db.GormDB())
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

	repo := NewGormRepository[TestEntity](db.GormDB())
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

	repo := NewGormRepository[TestEntity](db.GormDB())
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

	repo := NewGormRepository[TestEntity](db.GormDB())
	ctx := context.Background()

	assert.NoError(t, repo.Create(ctx, &TestEntity{ID: "d-1", Name: ps("to-del"), Age: pi(5)}))
	assert.NoError(t, repo.Delete(ctx, map[string]any{"id": "d-1"}))
	_, err := repo.GetByID(ctx, map[string]any{"id": "d-1"})
	assert.ErrorIs(t, err, ErrNotFound)

	// ensure soft-delete: Unscoped query should find the record and DeletedAt should be set
	var out TestEntity
	err = db.GormDB().Unscoped().First(&out, "id = ?", "d-1").Error
	assert.NoError(t, err)
	if assert.NotNil(t, out.DeletedAt) {
		assert.False(t, out.DeletedAt.Time.IsZero())
	}

	// hard delete the record and ensure it's gone
	assert.NoError(t, db.GormDB().Unscoped().Delete(&out).Error)
	err = db.GormDB().Unscoped().First(&out, "id = ?", "d-1").Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestGormrepositories_Update_NilInputAndDeletedAfterUpdate(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db.GormDB())
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
	repo := NewGormRepository[TestEntity](db.GormDB())
	ctx := context.Background()

	// close underlying DB to force errors
	cleanup()

	// GetByID should return an error (not ErrNotFound)
	_, err := repo.GetByID(ctx, map[string]any{"id": "any"})
	if err == nil || errors.Is(err, ErrNotFound) {
		t.Fatalf("expected DB error (not ErrNotFound), got %v", err)
	}

	// List should return error
	_, _, err = repo.List(ctx, repoContracts.ListParams{Page: 1, Size: 10, OrderBy: "", Order: ""})
	if err == nil {
		t.Fatalf("expected error from List when DB closed")
	}

	// Update should return error when DB closed
	in := &TestEntity{ID: "x"}
	err = repo.Update(ctx, map[string]any{"id": "any"}, in)
	if err == nil {
		t.Fatalf("expected error from Update when DB closed")
	}
}

func TestRepositories_ClosedDB_MapsToErrDBUnavailable(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db.GormDB())

	// close underlying DB
	if sqlDB, err := db.GormDB().DB(); err == nil {
		_ = sqlDB.Close()
	}

	ctx := context.Background()

	_, err := repo.GetByID(ctx, map[string]any{"id": "any"})
	assert.ErrorIs(t, err, ErrDBUnavailable)

	_, _, err = repo.List(ctx, repoContracts.ListParams{Page: 1, Size: 1, OrderBy: "", Order: ""})
	if err == nil {
		t.Fatalf("expected error from List when DB closed")
	}
	// Accept either ErrDBUnavailable (DB closed) or ordering validation error depending on implementation
	if !errors.Is(err, ErrDBUnavailable) && !strings.Contains(err.Error(), "both orderBy and order must be provided together") {
		t.Fatalf("expected ErrDBUnavailable or orderBy validation error, got: %v", err)
	}

	in := &TestEntity{ID: "x"}
	err = repo.Update(ctx, map[string]any{"id": "x"}, in)
	assert.ErrorIs(t, err, ErrDBUnavailable)
}
