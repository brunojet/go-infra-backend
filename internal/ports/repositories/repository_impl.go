package repositories

import (
	"context"

	"github.com/brunojet/go-infra-backend/debugassert"
	dbcts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
)

var _ contracts.Repository[contracts.Entity] = (*gormRepositoryImpl[contracts.Entity])(nil)

type gormRepositoryImpl[E contracts.Entity] struct {
	db *gorm.DB
}

// NewGormRepository creates a new generic GORM repository for the given entity type.
// It asserts that the provided database adapter is not nil and returns a repository instance.
func NewGormRepository[E contracts.Entity](db dbcts.DatabaseAdapter) *gormRepositoryImpl[E] {
	debugassert.Assert(db != nil, "NewGormRepository: db is nil") // NewGormRepository creates a new generic GORM repository for the given entity type.
	gdb, err := db.GormDB()
	debugassert.Assert(err == nil, "NewGormRepository: failed to obtain gorm DB")
	return &gormRepositoryImpl[E]{db: gdb}
}

// GormDB returns the underlying *gorm.DB instance used by the repository.
func (g *gormRepositoryImpl[E]) GormDB() *gorm.DB {
	return g.db
}

// DbFromContext retrieves a transaction from the context if present, otherwise returns the base DB with context.
func (g *gormRepositoryImpl[E]) DbFromContext(ctx context.Context) *gorm.DB {
	tx, err := TxFromContext(ctx)
	if err == nil && tx != nil {
		return tx // DbFromContext retrieves a transaction from the context if present, otherwise returns the base DB with context.
	}
	return g.db.WithContext(ctx)
}

// Create inserts a new entity into the database. If a conflict occurs, attempts to retrieve the existing entity.
func (g *gormRepositoryImpl[E]) Create(ctx context.Context, inOut *E) error {
	db := g.DbFromContext(ctx)
	tx := db.Model(new(E)).Create(inOut)
	err := MapTxError(tx)
	if err == ErrNotFound {
		err = getExistingWhenConflict(tx, inOut)
	}
	return err
}

// GetByID retrieves an entity by its primary key(s) and stores the result in 'out'.
func (g *gormRepositoryImpl[E]) GetByID(ctx context.Context, id map[string]any, out *E) error {
	db := g.DbFromContext(ctx)
	if err := getByScope(db, id, out); err != nil {
		return err
	}
	return nil
}

// List retrieves a list of entities matching the given parameters and returns the total count.
func (g *gormRepositoryImpl[E]) List(ctx context.Context, listParams contracts.ListParams, out *[]E) (int64, error) {
	db := g.DbFromContext(ctx)
	var total int64
	q, err := buildTxWithScopes[E](db, listParams.QueryParams.Scopes)
	if err != nil {
		return 0, err
	}
	tx := q.Count(&total)
	if err := MapTxError(tx); err != nil {
		return 0, err
	}
	if err := setOrderBy(q, listParams.OrderBy, listParams.Order); err != nil {
		return 0, err
	}
	if err := setPagination(q, listParams.Page, cap(*out)); err != nil {
		return 0, err
	}
	tx = q.Find(out)
	if err := MapTxError(tx); err != nil {
		return 0, err
	}
	return total, nil
}

// Update updates an entity matching the given scopes and refreshes the entity with the latest data from the database.
func (g *gormRepositoryImpl[E]) Update(ctx context.Context, scopes map[string]any, inOut *E) error {
	db := g.DbFromContext(ctx)
	tx, err := buildTxWithFilledScopes[E](db, scopes)
	if err != nil {
		return err
	}
	if tx = tx.Updates(inOut); tx.Error != nil {
		return MapTxError(tx)
	}
	return getByScope(db, scopes, inOut)
}

// Delete removes an entity matching the given scopes from the database.
func (g *gormRepositoryImpl[E]) Delete(ctx context.Context, scopes map[string]any) error {
	db := g.DbFromContext(ctx)
	tx, err := buildTxWithFilledScopes[E](db, scopes)
	if err != nil {
		return err
	}
	tx = tx.Delete(new(E))
	return MapTxError(tx)
}

// WithTx executes the given function within a database transaction, propagating the transaction in the context.
func (g *gormRepositoryImpl[E]) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return g.DbFromContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(contextWithTx(ctx, tx))
	})
}
