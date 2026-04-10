package contracts

import (
	"context"

	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
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

type ServiceMapper[C, R, U any, E contracts.Entity] interface {
	ToPostModel(dto C, model *E) error
	ToPatchModel(dto U, model *E) error
	ToDTO(model *E, dto *R) error
	GetModelKey(id string) (map[string]any, error)
	ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error)
}

type Service[C, R, U any] interface {
	Create(ctx context.Context, dto C, response *R) error
	List(ctx context.Context, params ListParams, response *[]R) (int64, error)
	GetByID(ctx context.Context, id string, response *R) error
	Update(ctx context.Context, id string, dto U, response *R) error
	Delete(ctx context.Context, id string) error
}

type NestedServiceMapper[C, R, U any, E contracts.Entity] interface {
	ServiceMapper[C, R, U, E]
	ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error)
	ApplyParentScopes(parentID string, model *E) error
}

type NestedService[C, R, U any] interface {
	Service[C, R, U]
	CreateNested(ctx context.Context, parentID string, dto C, response *R) error
	ListNested(ctx context.Context, parentID string, params ListParams, response *[]R) (int64, error)
}

type BffServiceMapper[C, R, U, CE, RE, UE any] interface {
	ToExternalPost(dto C, dtoExt *CE) error
	ToExternalPatch(dto U, dtoExt *UE) error
	ToInternalDTO(dtoExt *RE, dto *R) error
}

type BffService[C, R, U any] interface {
	Service[C, R, U]
}
