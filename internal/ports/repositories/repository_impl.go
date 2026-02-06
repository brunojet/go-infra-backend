package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/brunojet/go-infra-backend/internal/ports/repositories/contracts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GenericGormRepository[E contracts.Entity] struct {
	db *gorm.DB
}

func NewGormRepository[E contracts.Entity](db *gorm.DB) *GenericGormRepository[E] {
	return &GenericGormRepository[E]{db: db}
}

func (g *GenericGormRepository[E]) DB() *gorm.DB {
	return g.db
}

func (g *GenericGormRepository[E]) Create(ctx context.Context, inOut *E) error {
	tx := g.db.WithContext(ctx).Create(inOut)
	return MapTxError(tx)
}

func (g *GenericGormRepository[E]) GetByID(ctx context.Context, id map[string]any) (E, error) {
	var entity E
	tx := g.db.WithContext(ctx).First(&entity, id)
	if err := MapTxError(tx); err != nil {
		var zero E
		return zero, err
	}
	return entity, nil
}

func setPagination(q *gorm.DB, page, pageSize int) error {
	if page < 1 {
		return fmt.Errorf("page must be greater than zero")
	} else if pageSize <= 0 {
		return fmt.Errorf("pageSize must be greater than zero")
	}
	q.Limit(pageSize).Offset((page - 1) * pageSize)
	return nil
}

func setOrderBy(q *gorm.DB, orderBy, order string) error {
	if orderBy == "" || order == "" {
		return fmt.Errorf("both orderBy and order must be provided together")
	}
	orderClause := clause.OrderByColumn{Column: clause.Column{Name: orderBy}, Desc: (strings.ToLower(order) == "desc")}
	q.Order(orderClause)
	return nil
}

func getListSize(total, page, size int) int {
	capacity := size
	offset := (page - 1) * size
	remaining := total - offset
	if remaining < capacity {
		if remaining < 0 {
			capacity = 0
		} else {
			capacity = remaining
		}
	}
	return capacity
}

func (g *GenericGormRepository[E]) List(ctx context.Context, listParams contracts.ListParams) ([]E, int, error) {
	var total int64
	q := g.db.WithContext(ctx).Model(new(E))

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

	return items, int(total), nil
}

func (g *GenericGormRepository[E]) Update(ctx context.Context, id map[string]any, inOut *E) error {
	tx := g.db.WithContext(ctx).Where(id).Updates(inOut)
	if err := MapTxError(tx); err != nil {
		return err
	}

	tx = g.db.WithContext(ctx).First(inOut, id)
	return MapTxError(tx)
}

func (g *GenericGormRepository[E]) Delete(ctx context.Context, id map[string]any) error {
	tx := g.db.WithContext(ctx).Delete(new(E), id)
	return MapTxError(tx)
}

func (g *GenericGormRepository[E]) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ContextWithTx(ctx, tx))
	})
}

var _ contracts.Repository[contracts.Entity] = (*GenericGormRepository[contracts.Entity])(nil)
