package contracts

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
)

type ServiceMapper[D any, E repoContracts.Entity] interface {
	ToModel(dto *D, model *E)
	ToDTO(model *E, dto *D)
	GetModelKey(id string) (map[string]any, error)
	ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error)
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

type Service[D any, E repoContracts.Entity] interface {
	Create(ctx context.Context, dto *D) error
	List(ctx context.Context, params ListParams) ([]D, int64, error)
	GetByID(ctx context.Context, id string) (D, error)
	Update(ctx context.Context, id string, dto *D) error
	Delete(ctx context.Context, id string) error
}

type NestedServiceMapper[D any, E repoContracts.Entity] interface {
	ServiceMapper[D, E]
	ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error)
	ApplyParentScopes(parentID string, model *E) error
}

type NestedService[D any, E repoContracts.Entity] interface {
	Service[D, E]
	CreateNested(ctx context.Context, parentID string, dto *D) error
	ListNested(ctx context.Context, parentID string, params ListParams) ([]D, int64, error)
}
