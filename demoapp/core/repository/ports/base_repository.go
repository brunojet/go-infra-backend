package ports

import "context"

type CRUDRepository[T any] interface {
	List(ctx context.Context) ([]T, error)
	Get(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, v *T) (*T, error)
	// Patch applies partial updates to an entity identified by its ID.
	//
	// Conventions:
	//   - Patch returns (nil, nil) when the record does not exist.
	//   - updates uses GORM column names (e.g. "nome").
	Patch(ctx context.Context, id uint, updates map[string]any) (*T, error)
	Delete(ctx context.Context, id uint) error
}
