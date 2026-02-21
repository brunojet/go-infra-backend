package repositories

import (
	"context"
	"errors"
	"testing"

	dbadapters "github.com/brunojet/go-infra-backend/internal/database/adapters"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGormRepository_WithTx(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	err := repo.WithTx(ctx, func(ctx context.Context) error {
		if TxFromContext(ctx) == nil {
			return errors.New("tx not found in context")
		}
		return nil
	})
	assert.NoError(t, err)
}

func TestTxFromContext_AbsentReturnsNil(t *testing.T) {
	ctx := context.Background()
	tx := TxFromContext(ctx)
	assert.Nil(t, tx)
}

func TestTxFromContext_WithGormDB(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()
	gdb := mustGormDB(t, db)
	ctx := ContextWithTx(context.Background(), gdb)
	got := TxFromContext(ctx)
	assert.Same(t, gdb, got)
}

func TestTxFromContext_WrongTypeReturnsNil(t *testing.T) {
	// store a value under the same key but with wrong type
	ctx := context.WithValue(context.Background(), ctxKeyTx{}, "not-a-tx")
	got := TxFromContext(ctx)
	assert.Nil(t, got)
}

// open a minimal in-memory DB for utils tests (reuse logic from other tests)
func openSimpleMemoryDB(t *testing.T) (dbcontracts.DatabaseAdapter, *gorm.DB, func()) {
	t.Helper()
	db, err := dbadapters.NewSQLite("memory")
	require.NoError(t, err)
	gdb, err := db.GormDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	// migrate TestEntity used by utils
	require.NoError(t, gdb.AutoMigrate(&TestEntity{}, &RepoTestModel{}))
	return db, gdb, func() { _ = db.Close() }
}

func TestBuildTxWithConflict_ValidationAndSuccess(t *testing.T) {
	db, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	ctx := context.Background()

	// missing conflict columns when action != Error -> ErrInvalidConflictColumns
	_, err := buildTxWithConflict[TestEntity](ctx, gdb, repoContracts.CreateParams{OnConflictAction: repoContracts.ConflictActionIgnore})
	assert.ErrorIs(t, err, ErrInvalidConflictColumns)

	// empty conflict column name -> ErrInvalidConflictColumnName
	_, err = buildTxWithConflict[TestEntity](ctx, gdb, repoContracts.CreateParams{OnConflictAction: repoContracts.ConflictActionIgnore, ConflictColumns: map[string]any{"": 1}})
	assert.ErrorIs(t, err, ErrInvalidConflictColumnName)

	// valid ignore action with columns -> success
	tx, err := buildTxWithConflict[TestEntity](ctx, gdb, repoContracts.CreateParams{OnConflictAction: repoContracts.ConflictActionIgnore, ConflictColumns: map[string]any{"id": 1}})
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	// action == Error should return a tx even with empty columns
	tx2, err := buildTxWithConflict[TestEntity](ctx, gdb, repoContracts.CreateParams{OnConflictAction: repoContracts.ConflictActionError})
	assert.NoError(t, err)
	assert.NotNil(t, tx2)

	_ = db
}

func TestBuildTxWithScopes_ValidationAndFilledScopes(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	ctx := context.Background()

	// invalid scope: empty field name
	_, err := buildTxWithScopes[TestEntity](ctx, gdb, map[string]any{"": "v"})
	assert.True(t, errors.Is(err, ErrInvalidScope))

	// nil value
	_, err = buildTxWithScopes[TestEntity](ctx, gdb, map[string]any{"id": nil})
	assert.True(t, errors.Is(err, ErrInvalidScope))

	// filled scopes empty -> error
	_, err = buildTxWithFilledScopes[TestEntity](ctx, gdb, map[string]any{})
	assert.ErrorIs(t, err, ErrEmptyScopes)

	// valid filled scopes -> success
	tx, err := buildTxWithFilledScopes[TestEntity](ctx, gdb, map[string]any{"id": "x"})
	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

func TestGetByScope_ErrorsAndSuccess(t *testing.T) {
	db, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	ctx := context.Background()

	// empty scopes -> ErrEmptyScopes
	var out TestEntity
	err := getByScope(ctx, gdb, map[string]any{}, &out)
	assert.ErrorIs(t, err, ErrEmptyScopes)

	// create a record and fetch it
	require.NoError(t, gdb.Create(&TestEntity{ID: "g-1", Name: func() *string { s := "x"; return &s }(), Age: func() *int { i := 1; return &i }()}).Error)
	var got TestEntity
	err = getByScope(ctx, gdb, map[string]any{"id": "g-1"}, &got)
	assert.NoError(t, err)
	assert.Equal(t, "g-1", got.ID)
	_ = db
}

func TestSetOrderByAndPagination_ErrorsAndSuccess(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	q := gdb.Model(&RepoTestModel{})

	// orderBy missing
	assert.ErrorIs(t, setOrderBy(q, "", ""), ErrOrderByMissing)

	// pagination invalid
	assert.ErrorIs(t, setPagination(q, 0, 10), ErrInvalidPage)
	assert.ErrorIs(t, setPagination(q, 1, 0), ErrInvalidPageSize)

	// valid
	assert.NoError(t, setOrderBy(q, "id", "asc"))
	assert.NoError(t, setPagination(q, 2, 5))
}
