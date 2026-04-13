package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/infra/database/adapters"
	rpoerrs "github.com/brunojet/go-infra-backend/internal/ports/backend/repositories/errors"
	prterrs "github.com/brunojet/go-infra-backend/internal/ports/errors"
	dbcts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
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
		tx, err := TxFromContext(ctx)
		if err != nil || tx == nil {
			return errors.New("tx not found in context")
		}
		return nil
	})
	assert.NoError(t, err)
}

func TestTxFromContext_AbsentReturnsError(t *testing.T) {
	ctx := context.Background()
	tx, err := TxFromContext(ctx)
	assert.Nil(t, tx)
	assert.ErrorIs(t, err, rpoerrs.ErrInvalidTx)
}

func TestTxFromContext_WithGormDB(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()
	gdb := mustGormDB(t, db)
	ctx := contextWithTx(context.Background(), gdb)
	got, err := TxFromContext(ctx)
	assert.NoError(t, err)
	assert.Same(t, gdb, got)
}

func TestTxFromContext_WrongTypeReturnsError(t *testing.T) {
	// store a value under the same key but with wrong type
	ctx := context.WithValue(context.Background(), ctxKeyTx{}, "not-a-tx")
	got, err := TxFromContext(ctx)
	assert.Nil(t, got)
	assert.ErrorIs(t, err, rpoerrs.ErrInvalidTx)
}

// open a minimal in-memory DB for utils tests (reuse logic from other tests)
func openSimpleMemoryDB(t *testing.T) (dbcts.DatabaseAdapter, *gorm.DB, func()) {
	t.Helper()
	db, err := adapters.NewSQLite("memory")
	require.NoError(t, err)
	gdb, err := db.GormDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	// migrate TestEntity used by utils
	require.NoError(t, gdb.AutoMigrate(&TestEntity{}, &RepoTestModel{}))
	return db, gdb, func() { _ = db.Close() }
}

func TestAddOnConflictDoNothing_ValidationAndSuccess(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	err := AddOnConflictDoNothing(nil, "id")
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = AddOnConflictDoNothing(gdb)
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = AddOnConflictDoNothing(gdb, "")
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = gdb.Transaction(func(tx *gorm.DB) error {
		return AddOnConflictDoNothing(tx, "id")
	})
	assert.NoError(t, err)
}

func TestAddOnConflictUpdateAll_ValidationAndSuccess(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	err := AddOnConflictUpdateAll(nil, "id")
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = AddOnConflictUpdateAll(gdb)
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = AddOnConflictUpdateAll(gdb, "")
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = gdb.Transaction(func(tx *gorm.DB) error {
		return AddOnConflictUpdateAll(tx, "id")
	})
	assert.NoError(t, err)
}

func TestAddOnConflict_ValidationAndModes(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	err := addOnConflict(nil, conflictActionIgnore, "id")
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = addOnConflict(gdb, conflictActionIgnore)
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = addOnConflict(gdb, conflictActionIgnore, "")
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	err = gdb.Transaction(func(tx *gorm.DB) error {
		return addOnConflict(tx, conflictActionIgnore, "id")
	})
	assert.NoError(t, err)

	err = gdb.Transaction(func(tx *gorm.DB) error {
		return addOnConflict(tx, conflictActionUpdate, "id")
	})
	assert.NoError(t, err)

	err = gdb.Transaction(func(tx *gorm.DB) error {
		return addOnConflict(tx, conflictActionError)
	})
	assert.NoError(t, err)
}

func TestBuildTxWithScopes_ValidationAndFilledScopes(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	// invalid scope: empty field name
	_, err := buildTxWithScopes[TestEntity](gdb, map[string]any{"": "v"})
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	// nil value
	_, err = buildTxWithScopes[TestEntity](gdb, map[string]any{"id": nil})
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	// filled scopes empty -> error
	_, err = buildTxWithFilledScopes[TestEntity](gdb, map[string]any{})
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	// valid filled scopes -> success
	tx, err := buildTxWithFilledScopes[TestEntity](gdb, map[string]any{"id": "x"})
	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

func TestGetByScope_ErrorsAndSuccess(t *testing.T) {
	db, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	// empty scopes -> ErrEmptyScopes
	var out TestEntity
	err := getByScope(gdb, map[string]any{}, &out)
	assert.True(t, rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrInvalidParameters))

	// create a record and fetch it
	require.NoError(t, gdb.Create(&TestEntity{ID: "g-1", Name: func() *string { s := "x"; return &s }(), Age: func() *int { i := 1; return &i }()}).Error)
	var got TestEntity
	err = getByScope(gdb, map[string]any{"id": "g-1"}, &got)
	assert.NoError(t, err)
	assert.Equal(t, "g-1", got.ID)
	_ = db
}

func TestSetOrderByAndPagination_ErrorsAndSuccess(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	q := gdb.Model(&RepoTestModel{})

	// orderBy missing
	assert.ErrorIs(t, setOrderBy(q, "", ""), rpoerrs.ErrOrderByMissing)

	// pagination invalid
	assert.ErrorIs(t, setPagination(q, 0, 10), rpoerrs.ErrInvalidPage)
	assert.ErrorIs(t, setPagination(q, 1, 0), rpoerrs.ErrInvalidPageSize)

	// valid
	assert.NoError(t, setOrderBy(q, "id", "asc"))
	assert.NoError(t, setPagination(q, 2, 5))
}

func TestValidateTxWithUpdateLock_RequiresWhereAndTx(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	err := ValidateTxWithUpdateLock(gdb, LockValidationSpec[TestEntity]{
		WhereSQL: "",
	})
	assert.ErrorIs(t, err, rpoerrs.ErrLockValidationWhereClause)

	err = ValidateTxWithUpdateLock(gdb, LockValidationSpec[TestEntity]{
		WhereSQL: "id = ?",
		WhereArgs: []any{
			"x",
		},
	})
	assert.ErrorIs(t, err, rpoerrs.ErrRequiresTransaction)
}

func TestValidateTxWithUpdateLock_NotFoundReturnsNil(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	err := gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(contextWithTx(context.Background(), tx))
		return ValidateTxWithUpdateLock(tx, LockValidationSpec[TestEntity]{
			SelectColumns: []string{"id"},
			WhereSQL:      "id = ?",
			WhereArgs:     []any{"missing"},
		})
	})
	assert.NoError(t, err)
}

func TestValidateTxWithUpdateLock_FoundDefaultAndCallback(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	require.NoError(t, gdb.Create(&TestEntity{ID: "t-1", Name: func() *string { s := "john"; return &s }(), Age: func() *int { i := 30; return &i }()}).Error)

	err := gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(contextWithTx(context.Background(), tx))
		return ValidateTxWithUpdateLock(tx, LockValidationSpec[TestEntity]{
			SelectColumns: []string{"id"},
			WhereSQL:      "id = ?",
			WhereArgs:     []any{"t-1"},
		})
	})
	assert.ErrorIs(t, err, prterrs.ErrBusinessRuleViolation)

	err = gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(contextWithTx(context.Background(), tx))
		return ValidateTxWithUpdateLock(tx, LockValidationSpec[TestEntity]{
			SelectColumns: []string{"id", "name"},
			WhereSQL:      "id = ?",
			WhereArgs:     []any{"t-1"},
			BlockIfFound: func(found *TestEntity) error {
				if found.Name == nil || *found.Name != "john" {
					return prterrs.ErrBusinessRuleViolation
				}
				return nil
			},
		})
	})
	assert.NoError(t, err)

	expected := errors.New("custom-block")
	err = gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(contextWithTx(context.Background(), tx))
		return ValidateTxWithUpdateLock(tx, LockValidationSpec[TestEntity]{
			SelectColumns: []string{"id", "name"},
			WhereSQL:      "id = ?",
			WhereArgs:     []any{"t-1"},
			BlockIfFound: func(_ *TestEntity) error {
				return expected
			},
		})
	})
	assert.ErrorIs(t, err, expected)
}

func TestValidateTxWithUpdateLock_DBErrorPassthrough(t *testing.T) {
	_, gdb, cleanup := openSimpleMemoryDB(t)
	defer cleanup()

	err := gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(contextWithTx(context.Background(), tx))
		return ValidateTxWithUpdateLock(tx, LockValidationSpec[TestEntity]{
			WhereSQL: "id = ? AND (",
			WhereArgs: []any{
				"x",
			},
		})
	})
	assert.Error(t, err)
	assert.NotErrorIs(t, err, rpoerrs.ErrRequiresTransaction)
	assert.NotErrorIs(t, err, prterrs.ErrBusinessRuleViolation)
}
