package contracts

import (
	"context"

	"gorm.io/gorm"
)

type Entity interface {
	TableName() string
}

type ConflictAction int16

const (
	ConflictActionError ConflictAction = iota
	ConflictActionIgnore
	ConflictActionUpdate
)

type CreateParams struct {
	ConflictColumns  map[string]any
	OnConflictAction ConflictAction
}

type QueryParams struct {
	Scopes map[string]any
}

type ListParams struct {
	QueryParams
	Page    int
	Size    int
	OrderBy string
	Order   string
}

type Repository[E Entity] interface {
	Create(ctx context.Context, params CreateParams, inOut *E) error
	GetByID(ctx context.Context, id map[string]any) (E, error)
	List(ctx context.Context, params ListParams) ([]E, int64, error)
	Update(ctx context.Context, id map[string]any, inOut *E) error
	Delete(ctx context.Context, id map[string]any) error
	GormDB() *gorm.DB
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
