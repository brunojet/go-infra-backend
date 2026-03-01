package services

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type nestedServiceImpl[D any, E repoContracts.Entity] struct {
	repo   repoContracts.Repository[E]
	mapper svcContracts.NestedServiceMapper[D, E]
}

func NewNestedServiceImpl[D any, E repoContracts.Entity](
	r repoContracts.Repository[E],
	m svcContracts.NestedServiceMapper[D, E],
) svcContracts.NestedService[D, E] {
	return &nestedServiceImpl[D, E]{repo: r, mapper: m}
}

func (s *nestedServiceImpl[D, E]) CreateNested(ctx context.Context, parentScopes map[string]any, dto *D) error {
	var model E
	s.mapper.ToModel(dto, &model)
	if err := s.mapper.ApplyParentScopes(parentScopes, dto, &model); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, &model); err != nil {
		return err
	}
	s.mapper.ToDTO(&model, dto)
	return nil
}

func (s *nestedServiceImpl[D, E]) ListNested(ctx context.Context, parentScopes map[string]any, params repoContracts.ListParams) ([]D, int64, error) {
	queryScopes := make(map[string]any, len(parentScopes)+len(params.QueryParams.Scopes))
	for key, value := range params.QueryParams.Scopes {
		queryScopes[key] = value
	}
	for key, value := range parentScopes {
		queryScopes[key] = value
	}
	params.QueryParams.Scopes = queryScopes

	models, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]D, len(models))
	for i := range models {
		s.mapper.ToDTO(&models[i], &dtos[i])
	}

	return dtos, total, nil
}
