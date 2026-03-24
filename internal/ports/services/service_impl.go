package services

import (
	"context"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type serviceImpl[C, R, U any, E repoContracts.Entity] struct {
	repo   repoContracts.Repository[E]
	mapper contracts.ServiceMapper[C, R, U, E]
}

func NewServiceImpl[C, R, U any, E repoContracts.Entity](r repoContracts.Repository[E], m contracts.ServiceMapper[C, R, U, E]) contracts.Service[C, R, U, E] {
	return &serviceImpl[C, R, U, E]{repo: r, mapper: m}
}

func (s *serviceImpl[C, R, U, E]) Create(ctx context.Context, dto C) (R, error) {
	var model E
	s.mapper.ToPostModel(dto, &model)
	if err := s.repo.Create(ctx, &model); err != nil {
		var zero R
		return zero, err
	}
	var out R
	s.mapper.ToDTO(&model, &out)
	return out, nil
}

func (s *serviceImpl[C, R, U, E]) GetByID(ctx context.Context, id string) (R, error) {
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		var zero R
		return zero, err
	}
	model, err := s.repo.GetByID(ctx, key)
	if err != nil {
		var zero R
		return zero, err
	}
	var dto R
	s.mapper.ToDTO(&model, &dto)
	return dto, nil
}

func (s *serviceImpl[C, R, U, E]) List(ctx context.Context, params contracts.ListParams) ([]R, int64, error) {
	mappedScopes, err := s.mapper.ApplyQueryScopes(params.QueryParams.Scopes)
	if err != nil {
		return nil, 0, err
	}
	repoParams := toRepoListParams(params, mappedScopes)
	models, total, err := s.repo.List(ctx, repoParams)
	if err != nil {
		return nil, 0, err
	}
	out := make([]R, len(models))
	for i := range models {
		s.mapper.ToDTO(&models[i], &out[i])
	}
	return out, total, nil
}

func (s *serviceImpl[C, R, U, E]) Update(ctx context.Context, id string, dto U) (R, error) {
	var model E
	s.mapper.ToPatchModel(dto, &model)
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		var zero R
		return zero, err
	}
	if err := s.repo.Update(ctx, key, &model); err != nil {
		var zero R
		return zero, err
	}
	var out R
	s.mapper.ToDTO(&model, &out)
	return out, nil
}

func (s *serviceImpl[C, R, U, E]) Delete(ctx context.Context, id string) error {
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, key)
}
