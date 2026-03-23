package services

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type nestedServiceImpl[D any, E repoContracts.Entity] struct {
	svcContracts.Service[D, E]
	nestRpo repoContracts.Repository[E]
	nestMap svcContracts.NestedServiceMapper[D, E]
}

func NewNestedServiceImpl[D any, E repoContracts.Entity](
	r repoContracts.Repository[E],
	m svcContracts.NestedServiceMapper[D, E],
) svcContracts.NestedService[D, E] {
	return &nestedServiceImpl[D, E]{Service: NewServiceImpl(r, m), nestRpo: r, nestMap: m}
}

func (s *nestedServiceImpl[D, E]) CreateNested(ctx context.Context, parentID string, dto *D) error {
	var model E
	s.nestMap.ToModel(dto, &model)
	if err := s.nestMap.ApplyParentScopes(parentID, &model); err != nil {
		return err
	}
	if err := s.nestRpo.Create(ctx, &model); err != nil {
		return err
	}
	s.nestMap.ToDTO(&model, dto)
	return nil
}

func (s *nestedServiceImpl[D, E]) ListNested(ctx context.Context, parentID string, params svcContracts.ListParams) ([]D, int64, error) {
	mergedScopes, err := s.nestMap.ApplyParentQueryScopes(parentID, params.QueryParams.Scopes)
	if err != nil {
		return nil, 0, err
	}
	repoParams := toRepoListParams(params, mergedScopes)
	models, total, err := s.nestRpo.List(ctx, repoParams)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]D, len(models))
	for i := range models {
		s.nestMap.ToDTO(&models[i], &dtos[i])
	}
	return dtos, total, nil
}
