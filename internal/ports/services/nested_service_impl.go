package services

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type nestedServiceImpl[D any, E repoContracts.Entity] struct {
	nestRpo repoContracts.Repository[E]
	nestMap svcContracts.NestedServiceMapper[D, E]
	baseSvc svcContracts.Service[D, E]
}

func NewNestedServiceImpl[D any, E repoContracts.Entity](
	r repoContracts.Repository[E],
	m svcContracts.NestedServiceMapper[D, E],
) svcContracts.NestedService[D, E] {
	return &nestedServiceImpl[D, E]{nestRpo: r, nestMap: m, baseSvc: NewServiceImpl(r, m)}
}

func (s *nestedServiceImpl[D, E]) Create(ctx context.Context, dto *D) error {
	return s.baseSvc.Create(ctx, dto)
}

func (s *nestedServiceImpl[D, E]) List(ctx context.Context, params svcContracts.ListParams) ([]D, int64, error) {
	return s.baseSvc.List(ctx, params)
}

func (s *nestedServiceImpl[D, E]) GetByID(ctx context.Context, id string) (D, error) {
	return s.baseSvc.GetByID(ctx, id)
}

func (s *nestedServiceImpl[D, E]) Update(ctx context.Context, id string, dto *D) error {
	return s.baseSvc.Update(ctx, id, dto)
}

func (s *nestedServiceImpl[D, E]) Delete(ctx context.Context, id string) error {
	return s.baseSvc.Delete(ctx, id)
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
