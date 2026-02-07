package contracts

import (
	"context"

	"gorm.io/gorm"
)

type ListParams struct {
	Page    int
	Size    int
	OrderBy string
	Order   string
}

type Entity interface {
	TableName() string
}

type Repository[E Entity] interface {
	Create(ctx context.Context, inOut *E) error
	GetByID(ctx context.Context, id map[string]any) (E, error)
	List(ctx context.Context, listParams ListParams) ([]E, int, error)
	Update(ctx context.Context, id map[string]any, inOut *E) error
	Delete(ctx context.Context, id map[string]any) error
	GormDB() *gorm.DB
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
