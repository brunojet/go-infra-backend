package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/brunojet/go-infra-backend/debugassert"
	rpoerrs "github.com/brunojet/go-infra-backend/internal/ports/backend/repositories/errors"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	prterrs "github.com/brunojet/go-infra-backend/pkg/ports/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ctxKeyTx struct{}

type conflictAction uint8

const (
	conflictActionError conflictAction = iota
	conflictActionIgnore
	conflictActionUpdate
)

type LockValidationSpec[E contracts.Entity] = contracts.LockValidationSpec[E]

func hasTransactionInContext(ctx context.Context) bool {
	tx, err := TxFromContext(ctx)
	return err == nil && tx != nil
}

func isTransactionValid(tx *gorm.DB) bool {
	return tx != nil && tx.Statement != nil
}

func isTransactionAndContextValid(tx *gorm.DB) bool {
	return isTransactionValid(tx) && hasTransactionInContext(tx.Statement.Context)
}

// contextWithTx returns a new context that carries the given *gorm.DB transaction.
func contextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

// ContextWithTx returns a new context that carries the given *gorm.DB transaction.
// This is a thin exported wrapper used by higher-level packages and tests to
// annotate contexts with the transaction marker expected by TxFromContext.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return contextWithTx(ctx, tx)
}

func addOnConflict(tx *gorm.DB, action conflictAction, columnNames ...string) error {
	if !isTransactionValid(tx) {
		return rpoerrs.ErrInvalidTx
	}
	if action == conflictActionError {
		return nil
	}
	if len(columnNames) == 0 {
		return rpoerrs.ErrConflictColumnsMissing
	}

	columns := make([]clause.Column, 0, len(columnNames))
	for _, fieldName := range columnNames {
		if strings.TrimSpace(fieldName) == "" {
			return rpoerrs.ErrConflictColumnNameMissing
		}
		columns = append(columns, clause.Column{Name: fieldName})
	}

	onConflict := clause.OnConflict{Columns: columns}
	switch action {
	case conflictActionIgnore:
		onConflict.DoNothing = true
	case conflictActionUpdate:
		onConflict.UpdateAll = true
	}

	tx.Statement.AddClause(onConflict)
	return nil
}

func buildTxWithScopes[E contracts.Entity](db *gorm.DB, scopes map[string]any) (*gorm.DB, error) {
	tx := db.Model(new(E))
	for fieldName, fieldValue := range scopes {
		if fieldName == "" || fieldValue == nil {
			return nil, fmt.Errorf("%w: field '%s' has invalid value", rpoerrs.ErrScopeFieldMissing, fieldName)
		}
		tx = tx.Where(fmt.Sprintf("%s = ?", fieldName), fieldValue)
	}
	return tx, nil
}

func buildTxWithFilledScopes[E contracts.Entity](db *gorm.DB, scopes map[string]any) (*gorm.DB, error) {
	if len(scopes) == 0 {
		return nil, rpoerrs.ErrScopesMissing
	}
	return buildTxWithScopes[E](db, scopes)
}

func getByScope[E contracts.Entity](db *gorm.DB, scopes map[string]any, out *E) error {
	tx, err := buildTxWithFilledScopes[E](db, scopes)
	if err != nil {
		return err
	}
	tx = tx.First(out)
	return rpoerrs.MapTxError(tx)
}

func getExistingWhenConflict[E contracts.Entity](tx *gorm.DB, out *E) error {
	debugassert.Assert(out != nil, "getExistingWhenConflict: out parameter is nil")
	if conflictTx := (*out).WhereOnConflict(tx).First(out); conflictTx.Error != nil || conflictTx.RowsAffected == 0 {
		return rpoerrs.ErrConflictValidationFailed
	}
	return rpoerrs.ErrConflictValidationRequired
}

func setOrderBy(q *gorm.DB, orderBy, order string) error {
	if len(orderBy) == 0 {
		return rpoerrs.ErrOrderByMissing
	}
	orderClause := clause.OrderByColumn{Column: clause.Column{Name: orderBy}, Desc: (strings.ToLower(order) == "desc")}
	q.Order(orderClause)
	return nil
}

func setPagination(q *gorm.DB, page, pageSize int) error {
	if page < 1 {
		return rpoerrs.ErrInvalidPage
	} else if pageSize <= 0 {
		return rpoerrs.ErrInvalidPageSize
	}
	q.Limit(pageSize).Offset((page - 1) * pageSize)
	return nil
}

// TxFromContext extracts a *gorm.DB transaction from the context.
// Returns ErrInvalidTx when no transaction is present or when value has wrong type.
func TxFromContext(ctx context.Context) (*gorm.DB, error) {
	if v := ctx.Value(ctxKeyTx{}); v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx, nil
		}
	}
	return nil, rpoerrs.ErrInvalidTx
}

// ValidateTxWithUpdateLock performs a reusable business-rule validation pattern:
// 1) requires an explicit transaction
// 2) executes SELECT ... FOR UPDATE with provided where clause
// 3) uses tx.Statement.Context when present (for observability/tracing propagation)
// 3) returns nil on not found (no conflict)
// 4) evaluates optional callback for custom blocking rules when a record is found
func ValidateTxWithUpdateLock[E contracts.Entity](tx *gorm.DB, spec contracts.LockValidationSpec[E]) error {
	whereSQL := strings.TrimSpace(spec.WhereSQL)
	if whereSQL == "" {
		return rpoerrs.ErrLockValidationWhere
	}
	if !isTransactionAndContextValid(tx) {
		return rpoerrs.ErrRequiresTransaction
	}
	for _, arg := range spec.WhereArgs {
		if strArg, ok := arg.(string); ok && strings.TrimSpace(strArg) == "" {
			return rpoerrs.ErrLockValidationWhereArgument
		}
		if intArg, ok := arg.(int64); ok && intArg == 0 {
			return rpoerrs.ErrLockValidationWhereArgument
		}
	}
	q := tx.Session(&gorm.Session{NewDB: true}).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	if len(spec.SelectColumns) > 0 {
		q = q.Select(spec.SelectColumns)
	}
	q = q.Where(whereSQL, spec.WhereArgs...)
	var found E
	err := q.First(&found).Error
	if err == nil {
		if spec.BlockIfFound != nil {
			return spec.BlockIfFound(&found)
		}
		return prterrs.ErrBusinessRuleViolation
	}
	if err = rpoerrs.MapDbError(err); err != nil && err != rpoerrs.ErrNotFound {
		return err
	}
	return nil
}

func AddOnConflictDoNothing(tx *gorm.DB, columnNames ...string) error {
	return addOnConflict(tx, conflictActionIgnore, columnNames...)
}

func AddOnConflictUpdateAll(tx *gorm.DB, columnNames ...string) error {
	return addOnConflict(tx, conflictActionUpdate, columnNames...)
}

// WhereOnConflict aplica filtros de conflito ao tx com base nos escopos fornecidos.
func WhereOnConflict(tx *gorm.DB, scopes ...contracts.ConflictScope) *gorm.DB {
	for _, scope := range scopes {
		tx = tx.Where(scope.ColumnName+" = ?", scope.ColumnValue)
	}
	return tx
}
