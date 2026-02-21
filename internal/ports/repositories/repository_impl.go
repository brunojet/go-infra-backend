package repositories

import (
	"context"
	"fmt"

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

func (g *gormRepositoryImpl[E]) Create(ctx context.Context, params contracts.CreateParams, inOut *E) error {
	tx, err := buildTxWithConflict[E](ctx, g.db, params)
	if err != nil {
		return err
	}
	if tx = tx.Create(inOut); tx.Error != nil {
		return MapTxError(tx)
	}
	if tx.RowsAffected == 0 {
		if err := getByScope(ctx, g.db, params.ConflictColumns, inOut); err != nil {
			return fmt.Errorf("failed to load existing entity after conflict: %w", err)
		}
	}
	return nil
}

func (g *gormRepositoryImpl[E]) GetByID(ctx context.Context, id map[string]any) (E, error) {
	var entity E
	if err := getByScope(ctx, g.db, id, &entity); err != nil {
		var zero E
		return zero, err
	}
	return entity, nil
}

func (g *gormRepositoryImpl[E]) List(ctx context.Context, listParams contracts.ListParams) ([]E, int64, error) {
	var total int64
	q, err := buildTxWithScopes[E](ctx, g.db, listParams.QueryParams.Scopes)
	if err != nil {
		return nil, 0, err
	}

	tx := q.Count(&total)
	if err := MapTxError(tx); err != nil {
		return nil, 0, err
	}

	if err := setOrderBy(q, listParams.OrderBy, listParams.Order); err != nil {
		return nil, 0, err
	}

	if err := setPagination(q, listParams.Page, listParams.Size); err != nil {
		return nil, 0, err
	}

	items := make([]E, 0, getListSize(int(total), listParams.Page, listParams.Size))

	tx = q.Find(&items)
	if err := MapTxError(tx); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (g *gormRepositoryImpl[E]) Update(ctx context.Context, scopes map[string]any, inOut *E) error {
	tx, err := buildTxWithFilledScopes[E](ctx, g.db, scopes)
	if err != nil {
		return err
	}
	if tx = tx.Updates(inOut); tx.Error != nil {
		return MapTxError(tx)
	}
	return getByScope(ctx, g.db, scopes, inOut)
}

func (g *gormRepositoryImpl[E]) Delete(ctx context.Context, scopes map[string]any) error {
	tx, err := buildTxWithFilledScopes[E](ctx, g.db, scopes)
	if err != nil {
		return err
	}
	tx = tx.Delete(new(E))
	return MapTxError(tx)
}

func (g *gormRepositoryImpl[E]) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ContextWithTx(ctx, tx))
	})
}
