package services

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type serviceImpl[D any, E repoContracts.Entity] struct {
	repo   repoContracts.Repository[E]
	mapper contracts.ServiceMapper[D, E]
}

func NewServiceImpl[D any, E repoContracts.Entity](r repoContracts.Repository[E], m contracts.ServiceMapper[D, E]) contracts.Service[D, E] {
	return &serviceImpl[D, E]{repo: r, mapper: m}
}

func (s *serviceImpl[D, E]) Create(ctx context.Context, dto *D) error {
	var model E
	s.mapper.ToModel(dto, &model)
	if err := s.repo.Create(ctx, &model); err != nil {
		return err
	}
	s.mapper.ToDTO(&model, dto)
	return nil
}

func (s *serviceImpl[D, E]) GetByID(ctx context.Context, id string) (D, error) {
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		var zero D
		return zero, err
	}
	model, err := s.repo.GetByID(ctx, key)
	if err != nil {
		var zero D
		return zero, err
	}
	var dto D
	s.mapper.ToDTO(&model, &dto)
	return dto, nil
}

func (s *serviceImpl[D, E]) List(ctx context.Context, params contracts.ListParams) ([]D, int64, error) {
	mappedScopes, err := s.mapper.ApplyQueryScopes(params.QueryParams.Scopes)
	if err != nil {
		return nil, 0, err
	}
	repoParams := toRepoListParams(params, mappedScopes)
	models, total, err := s.repo.List(ctx, repoParams)
	if err != nil {
		return nil, 0, err
	}
	out := make([]D, len(models))
	for i := range models {
		s.mapper.ToDTO(&models[i], &out[i])
	}
	return out, total, nil
}

func (s *serviceImpl[D, E]) Update(ctx context.Context, id string, dto *D) error {
	var model E
	s.mapper.ToModel(dto, &model)
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	if err := s.repo.Update(ctx, key, &model); err != nil {
		return err
	}
	s.mapper.ToDTO(&model, dto)
	return nil
}

func (s *serviceImpl[D, E]) Delete(ctx context.Context, id string) error {
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, key)
}
