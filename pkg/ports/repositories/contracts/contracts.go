package contracts

import (
	"context"

	"gorm.io/gorm"
)

type Entity interface {
	TableName() string
	WhereOnConflict(tx *gorm.DB) *gorm.DB
}

type QueryParams struct {
	Scopes map[string]any
}

type ListParams struct {
	QueryParams
	Page    int
	OrderBy string
	Order   string
}

// ConflictScope representa um filtro de coluna/valor para resolução de conflito.
type ConflictScope struct {
	ColumnName  string
	ColumnValue any
}

// LockValidationSpec defines a reusable contract for transaction+locking
// validations based on SELECT ... FOR UPDATE.
//
//   - SelectColumns: optional columns to fetch from the conflicting row.
//   - WhereSQL/WhereArgs: required filter used to find potential conflicting row.
//   - BlockIfFound: optional callback for custom business-rule checks.
//     If nil and a row is found, caller should consider it a business-rule violation.
type LockValidationSpec[E Entity] struct {
	SelectColumns []string
	WhereSQL      string
	WhereArgs     []any
	BlockIfFound  func(found *E) error
}

type Repository[E Entity] interface {
	Create(ctx context.Context, inOut *E) error
	GetByID(ctx context.Context, id map[string]any, out *E) error
	List(ctx context.Context, params ListParams, out *[]E) (int64, error)
	Update(ctx context.Context, id map[string]any, inOut *E) error
	Delete(ctx context.Context, id map[string]any) error
	GormDB() *gorm.DB
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
