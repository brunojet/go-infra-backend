package repositories

import (
	"bytes"
	"context"
	"encoding/json"
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

// getExistingWhenConflict fetches the existing conflicting record into *out and verifies
// that the original intent is a subset of what is already persisted (idempotency check).
// Returns nil when the conflict is idempotent, ErrConflictValidationFailed otherwise.
func getExistingWhenConflict[E contracts.Entity](tx *gorm.DB, out *E) error {
	debugassert.Assert(out != nil, "getExistingWhenConflict: out parameter is nil")
	// ON CONFLICT DO NOTHING means no rows were inserted and GORM does not write
	// auto-fields back into *out, so *out still holds the caller's original intent.
	// Use a zeroed-out fresh value for the SELECT so GORM does not treat *out's
	// non-zero fields as implicit WHERE conditions.
	selectTx := tx.Session(&gorm.Session{NewDB: true})
	intent := *out
	var fresh E
	if conflictTx := intent.WhereOnConflict(selectTx).First(&fresh); conflictTx.Error != nil || conflictTx.RowsAffected == 0 {
		return rpoerrs.NewDatabaseError(rpoerrs.DBErrConstraint, conflictTx.Error)
	}
	if !isModelSubset(intent, fresh) {
		return rpoerrs.ErrConflictValidationFailed
	}
	*out = fresh
	return nil
}

// isModelSubset reports whether all non-null/non-zero fields in original exist with
// equal values in existing. Uses JSON encoding to avoid direct reflect usage.
func isModelSubset[E any](original, existing E) bool {
	ma, err := toJSONMap(original)
	if err != nil {
		return false
	}
	mb, err := toJSONMap(existing)
	if err != nil {
		return false
	}
	return isSubsetMapNonNull(ma, mb)
}

func toJSONMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	return m, json.Unmarshal(b, &m)
}

// isSubsetMapNonNull checks that every non-null, non-zero field in a exists in b with the same value.
func isSubsetMapNonNull(a, b map[string]any) bool {
	for k, va := range a {
		if isNullOrZero(va) {
			continue
		}
		vb, ok := b[k]
		if !ok {
			return false
		}
		if ma, ok := va.(map[string]any); ok {
			if mb, ok := vb.(map[string]any); ok {
				if !isSubsetMapNonNull(ma, mb) {
					return false
				}
				continue
			}
			return false
		}
		ba, _ := json.Marshal(va)
		bb, _ := json.Marshal(vb)
		if !bytes.Equal(ba, bb) {
			return false
		}
	}
	return true
}

// isNullOrZero returns true when v is considered "not intentionally set" by the caller:
//   - JSON null  (nil)
//   - numeric zero  (float64(0))
//   - sql.Null* pattern: a map with a "Valid" key set to false.
//     database/sql types (NullString, NullTime, NullInt*, etc.) marshal as
//     {"Valid": <bool>, ...} when there is no MarshalJSON implementation,
//     and {"Valid": false} means the value was not provided.
func isNullOrZero(v any) bool {
	switch val := v.(type) {
	case nil:
		return true
	case float64:
		return val == 0
	case map[string]any:
		if valid, ok := val["Valid"]; ok {
			if b, ok := valid.(bool); ok {
				return !b
			}
		}
		return false
	default:
		return false
	}
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
		return rpoerrs.ErrLockValidationWhereClause
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
	if err = rpoerrs.MapDbError(err); err != nil && !rpoerrs.IsDatabaseErrorKind(err, rpoerrs.DBErrNotFound) {
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
