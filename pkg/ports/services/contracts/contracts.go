package contracts

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
)

type AnyInt interface {
	~int64 | ~int32 | ~int16
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

type ServiceMapper[C, R, U any, E repoContracts.Entity] interface {
	ToPostModel(dto C, model *E) error
	ToPatchModel(dto U, model *E) error
	ToDTO(model *E, dto *R) error
	GetModelKey(id string) (map[string]any, error)
	ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error)
}

type Service[C, R, U any, E repoContracts.Entity] interface {
	Create(ctx context.Context, dto C) (R, error)
	List(ctx context.Context, params ListParams) ([]R, int64, error)
	GetByID(ctx context.Context, id string) (R, error)
	Update(ctx context.Context, id string, dto U) (R, error)
	Delete(ctx context.Context, id string) error
}

type NestedServiceMapper[C, R, U any, E repoContracts.Entity] interface {
	ServiceMapper[C, R, U, E]
	ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error)
	ApplyParentScopes(parentID string, model *E) error
}

type NestedService[C, R, U any, E repoContracts.Entity] interface {
	Service[C, R, U, E]
	CreateNested(ctx context.Context, parentID string, dto C) (R, error)
	ListNested(ctx context.Context, parentID string, params ListParams) ([]R, int64, error)
}
