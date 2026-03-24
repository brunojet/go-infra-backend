package services

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type nestedServiceImpl[C, R, U any, E repoContracts.Entity] struct {
	svcContracts.Service[C, R, U, E]
	nestRpo repoContracts.Repository[E]
	nestMap svcContracts.NestedServiceMapper[C, R, U, E]
}

func NewNestedServiceImpl[C, R, U any, E repoContracts.Entity](
	r repoContracts.Repository[E],
	m svcContracts.NestedServiceMapper[C, R, U, E],
) svcContracts.NestedService[C, R, U, E] {
	return &nestedServiceImpl[C, R, U, E]{Service: NewServiceImpl(r, m), nestRpo: r, nestMap: m}
}

func (s *nestedServiceImpl[C, R, U, E]) CreateNested(ctx context.Context, parentID string, dto C) (R, error) {
	var model E
	s.nestMap.ToPostModel(dto, &model)
	if err := s.nestMap.ApplyParentScopes(parentID, &model); err != nil {
		var zero R
		return zero, err
	}
	if err := s.nestRpo.Create(ctx, &model); err != nil {
		var zero R
		return zero, err
	}
	var out R
	s.nestMap.ToDTO(&model, &out)
	return out, nil
}

func (s *nestedServiceImpl[C, R, U, E]) ListNested(ctx context.Context, parentID string, params svcContracts.ListParams) ([]R, int64, error) {
	mergedScopes, err := s.nestMap.ApplyParentQueryScopes(parentID, params.QueryParams.Scopes)
	if err != nil {
		return nil, 0, err
	}
	repoParams := toRepoListParams(params, mergedScopes)
	models, total, err := s.nestRpo.List(ctx, repoParams)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]R, len(models))
	for i := range models {
		s.nestMap.ToDTO(&models[i], &dtos[i])
	}
	return dtos, total, nil
}
