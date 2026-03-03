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

type txMarker interface {
	Commit() error
	Rollback() error
}

type LockValidationSpec[E contracts.Entity] = contracts.LockValidationSpec[E]

// ContextWithTx returns a new context that carries the given *gorm.DB transaction.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

// TxFromContext extracts a *gorm.DB transaction from the context.
// Returns ErrInvalidTx when no transaction is present or when value has wrong type.
func TxFromContext(ctx context.Context) (*gorm.DB, error) {
	if v := ctx.Value(ctxKeyTx{}); v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx, nil
		}
		return nil, ErrInvalidTx
	}
	return nil, ErrInvalidTx
}

func GetContextFromTx(tx *gorm.DB) context.Context {
	if tx.Statement != nil && tx.Statement.Context != nil {
		return tx.Statement.Context
	}
	return context.Background()
}

// ValidateTxWithUpdateLock performs a reusable business-rule validation pattern:
// 1) requires an explicit transaction
// 2) executes SELECT ... FOR UPDATE with provided where clause
// 3) uses tx.Statement.Context when present (for observability/tracing propagation)
// 3) returns nil on not found (no conflict)
// 4) evaluates optional callback for custom blocking rules when a record is found
func ValidateTxWithUpdateLock[E contracts.Entity](tx *gorm.DB, spec contracts.LockValidationSpec[E]) error {
	if tx == nil {
		return ErrInvalidTx
	}
	if strings.TrimSpace(spec.WhereSQL) == "" {
		return ErrLockValidationWhere
	}
	if _, ok := tx.Statement.ConnPool.(txMarker); !ok {
		return ErrRequiresTransaction
	}

	ctx := GetContextFromTx(tx)

	q := tx.Session(&gorm.Session{NewDB: true}).WithContext(ctx).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	if len(spec.SelectColumns) > 0 {
		q = q.Select(spec.SelectColumns)
	}
	q = q.Where(spec.WhereSQL, spec.WhereArgs...)

	var found E
	err := q.First(&found).Error
	if err == nil {
		if spec.BlockIfFound != nil {
			return spec.BlockIfFound(&found)
		}
		return ErrBusinessRuleViolation
	}
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return MapDbError(err)
}

func AddOnConflict(tx *gorm.DB, action contracts.ConflictAction, columnNames ...string) error {
	if tx == nil || tx.Statement == nil {
		return ErrInvalidTx
	}
	if action == contracts.ConflictActionError {
		return nil
	}
	if len(columnNames) == 0 {
		return ErrInvalidConflictColumns
	}

	columns := make([]clause.Column, 0, len(columnNames))
	for _, fieldName := range columnNames {
		if strings.TrimSpace(fieldName) == "" {
			return ErrInvalidConflictColumnName
		}
		columns = append(columns, clause.Column{Name: fieldName})
	}

	onConflict := clause.OnConflict{Columns: columns}
	switch action {
	case contracts.ConflictActionIgnore:
		onConflict.DoNothing = true
	case contracts.ConflictActionUpdate:
		onConflict.UpdateAll = true
	}

	tx.Statement.AddClause(onConflict)
	return nil
}

func AddOnConflictDoNothing(tx *gorm.DB, columnNames ...string) error {
	return AddOnConflict(tx, contracts.ConflictActionIgnore, columnNames...)
}

func AddOnConflictUpdateAll(tx *gorm.DB, columnNames ...string) error {
	return AddOnConflict(tx, contracts.ConflictActionUpdate, columnNames...)
}

func buildTxWithScopes[E contracts.Entity](db *gorm.DB, scopes map[string]any) (*gorm.DB, error) {
	tx := db.Model(new(E))
	for fieldName, fieldValue := range scopes {
		if fieldName == "" || fieldValue == nil {
			return nil, fmt.Errorf("%w: field '%s' has invalid value", ErrInvalidScope, fieldName)
		}
		tx = tx.Where(fmt.Sprintf("%s = ?", fieldName), fieldValue)
	}
	return tx, nil
}

func buildTxWithFilledScopes[E contracts.Entity](db *gorm.DB, scopes map[string]any) (*gorm.DB, error) {
	if len(scopes) == 0 {
		return nil, ErrEmptyScopes
	}
	return buildTxWithScopes[E](db, scopes)
}

func getByScope[E contracts.Entity](db *gorm.DB, scopes map[string]any, out *E) error {
	tx, err := buildTxWithFilledScopes[E](db, scopes)
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
