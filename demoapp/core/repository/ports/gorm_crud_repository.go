package ports

import (
	"context"

	"gorm.io/gorm"
)

// GormCRUDRepository is a reusable CRUD repository implementation backed by GORM.
//
// Conventions:
//   - Get returns (nil, nil) when the record does not exist.
//   - Delete is idempotent (deleting a missing record returns nil).
type GormCRUDRepository[T any] struct {
	db *gorm.DB
}

func NewGormCRUDRepository[T any](db *gorm.DB) *GormCRUDRepository[T] {
	return &GormCRUDRepository[T]{db: db}
}

func (r *GormCRUDRepository[T]) List(ctx context.Context) ([]T, error) {
	var items []T
	if err := r.db.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GormCRUDRepository[T]) Get(ctx context.Context, id uint) (*T, error) {
	entity := new(T)
	tx := r.db.WithContext(ctx).First(entity, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}
	return entity, nil
}

func (r *GormCRUDRepository[T]) Create(ctx context.Context, v *T) (*T, error) {
	if err := r.db.WithContext(ctx).Create(v).Error; err != nil {
		return nil, err
	}
	return v, nil
}

func (r *GormCRUDRepository[T]) Delete(ctx context.Context, id uint) error {
	// GORM delete is idempotent; deleting a missing record is not an error.
	return r.db.WithContext(ctx).Delete(new(T), id).Error
}

func (r *GormCRUDRepository[T]) Patch(ctx context.Context, id uint, updates map[string]any) (*T, error) {
	if len(updates) == 0 {
		return r.Get(ctx, id)
	}

	// Check existence first to preserve the (nil, nil) contract.
	current, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil
	}

	if err := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}
