package contracts

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
)

type NestedServiceMapper[D any, E repoContracts.Entity] interface {
	ApplyParentScopes(parentScopes map[string]any, dto *D, model *E) error
	ToModel(dto *D, model *E)
	ToDTO(model *E, dto *D)
}

type NestedService[D any, E repoContracts.Entity] interface {
	CreateNested(ctx context.Context, parentScopes map[string]any, dto *D) error
	ListNested(ctx context.Context, parentScopes map[string]any, params repoContracts.ListParams) ([]D, int64, error)
}
