package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ctxKeyTx struct{}

// ContextWithTx returns a new context that carries the given *gorm.DB transaction.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

// TxFromContext extracts a *gorm.DB transaction from the context, or nil if not present.
func TxFromContext(ctx context.Context) *gorm.DB {
	if v := ctx.Value(ctxKeyTx{}); v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return nil
}

func buildTxWithConflict[E contracts.Entity](ctx context.Context, db *gorm.DB, params contracts.CreateParams) (*gorm.DB, error) {
	if params.OnConflictAction != contracts.ConflictActionError && len(params.ConflictColumns) == 0 {
		return nil, ErrInvalidConflictColumns
	}
	tx := db.WithContext(ctx).Model(new(E))
	if params.OnConflictAction == contracts.ConflictActionError {
		return tx, nil
	}
	columns := make([]clause.Column, 0, len(params.ConflictColumns))
	for fieldName := range params.ConflictColumns {
		if fieldName == "" {
			return nil, ErrInvalidConflictColumnName
		}
		columns = append(columns, clause.Column{Name: fieldName})
	}
	onConflict := clause.OnConflict{
		Columns:   columns,
		DoNothing: (params.OnConflictAction == contracts.ConflictActionIgnore),
		UpdateAll: (params.OnConflictAction == contracts.ConflictActionUpdate),
	}
	tx = tx.Clauses(onConflict, clause.Returning{})
	return tx, nil
}

func buildTxWithScopes[E contracts.Entity](ctx context.Context, db *gorm.DB, scopes map[string]any) (*gorm.DB, error) {
	tx := db.WithContext(ctx).Model(new(E))
	for fieldName, fieldValue := range scopes {
		if fieldName == "" || fieldValue == nil {
			return nil, fmt.Errorf("%w: field '%s' has invalid value", ErrInvalidScope, fieldName)
		}
		tx = tx.Where(fmt.Sprintf("%s = ?", fieldName), fieldValue)
	}
	return tx, nil
}

func buildTxWithFilledScopes[E contracts.Entity](ctx context.Context, db *gorm.DB, scopes map[string]any) (*gorm.DB, error) {
	if len(scopes) == 0 {
		return nil, ErrEmptyScopes
	}
	return buildTxWithScopes[E](ctx, db, scopes)
}

func getByScope[E contracts.Entity](ctx context.Context, db *gorm.DB, scopes map[string]any, out *E) error {
	tx, err := buildTxWithFilledScopes[E](ctx, db, scopes)
	if err != nil {
		return err
	}
	tx = tx.First(out)
	return MapTxError(tx)
}

func setOrderBy(q *gorm.DB, orderBy, order string) error {
	if len(orderBy) == 0 {
		return ErrOrderByMissing
	}
	orderClause := clause.OrderByColumn{Column: clause.Column{Name: orderBy}, Desc: (strings.ToLower(order) == "desc")}
	q.Order(orderClause)
	return nil
}

func setPagination(q *gorm.DB, page, pageSize int) error {
	if page < 1 {
		return ErrInvalidPage
	} else if pageSize <= 0 {
		return ErrInvalidPageSize
	}
	q.Limit(pageSize).Offset((page - 1) * pageSize)
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
