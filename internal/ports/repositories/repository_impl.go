package repositories

import (
	"context"

	"github.com/brunojet/go-infra-backend/debugassert"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
)

var _ contracts.Repository[contracts.Entity] = (*gormRepositoryImpl[contracts.Entity])(nil)

type gormRepositoryImpl[E contracts.Entity] struct {
	db *gorm.DB
}

func NewGormRepository[E contracts.Entity](db dbcontracts.DatabaseAdapter) *gormRepositoryImpl[E] {
	debugassert.Assert(db != nil, "NewGormRepository: db is nil")
	gdb, err := db.GormDB()
	debugassert.Assert(err == nil, "NewGormRepository: failed to obtain gorm DB")
	return &gormRepositoryImpl[E]{db: gdb}
}

func (g *gormRepositoryImpl[E]) GormDB() *gorm.DB {
	return g.db
}

func (g *gormRepositoryImpl[E]) dbFromContext(ctx context.Context) *gorm.DB {
	tx, err := TxFromContext(ctx)
	if err == nil && tx != nil {
		return tx
	}
	return g.db.WithContext(ctx)
}

func (g *gormRepositoryImpl[E]) Create(ctx context.Context, inOut *E) error {
	db := g.dbFromContext(ctx)
	tx := db.Model(new(E)).Create(inOut)
	err := MapTxError(tx)
	if err == ErrNotFound {
		err = getExistingWhenConflict(tx, inOut)
	}
	return err
}

func (g *gormRepositoryImpl[E]) GetByID(ctx context.Context, id map[string]any, out *E) error {
	db := g.dbFromContext(ctx)
	if err := getByScope(db, id, out); err != nil {
		return err
	}
	return nil
}

func (g *gormRepositoryImpl[E]) List(ctx context.Context, listParams contracts.ListParams, out *[]E) (int64, error) {
	db := g.dbFromContext(ctx)
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

func (g *gormRepositoryImpl[E]) Update(ctx context.Context, scopes map[string]any, inOut *E) error {
	db := g.dbFromContext(ctx)
	tx, err := buildTxWithFilledScopes[E](db, scopes)
	if err != nil {
		return err
	}
	if tx = tx.Updates(inOut); tx.Error != nil {
		return MapTxError(tx)
	}
	return getByScope(db, scopes, inOut)
}

func (g *gormRepositoryImpl[E]) Delete(ctx context.Context, scopes map[string]any) error {
	db := g.dbFromContext(ctx)
	tx, err := buildTxWithFilledScopes[E](db, scopes)
	if err != nil {
		return err
	}
	tx = tx.Delete(new(E))
	return MapTxError(tx)
}

func (g *gormRepositoryImpl[E]) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(contextWithTx(ctx, tx))
	})
}
