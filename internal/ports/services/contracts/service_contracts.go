package contracts

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/internal/ports/repositories/contracts"
)

type ServiceMapper[D any, E repoContracts.Entity] interface {
	GetModelKey(id string) (map[string]any, error)
	ToModel(dto *D, model *E)
	ToDTO(model *E, dto *D)
}

type Service[D any, E repoContracts.Entity] interface {
	Create(ctx context.Context, dto *D) error
	GetByID(ctx context.Context, id string) (D, error)
	List(ctx context.Context, size int) ([]D, error)
	Update(ctx context.Context, id string, dto *D) error
	Delete(ctx context.Context, id string) error
}
