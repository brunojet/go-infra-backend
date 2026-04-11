package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/ports/backend/repositories"
	"github.com/brunojet/go-infra-backend/internal/utils"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
)

type serviceImpl[C, R, U any, E rpocts.Entity] struct {
	rpo    rpocts.Repository[E]
	mapper contracts.ServiceMapper[C, R, U, E]
}

func NewServiceImpl[C, R, U any, E rpocts.Entity](r rpocts.Repository[E], m contracts.ServiceMapper[C, R, U, E]) contracts.Service[C, R, U] {
	return &serviceImpl[C, R, U, E]{rpo: r, mapper: m}
}

func (s *serviceImpl[C, R, U, E]) Create(ctx context.Context, request C, response *R) error {
	var model E
	if err := s.mapper.ToPostModel(request, &model); err != nil {
		return err
	}
	conflictValidationNeeded := false
	err := s.rpo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.rpo.Create(txCtx, &model); err != nil {
			if err != repositories.ErrConflictValidationRequired {
				return err
			}
			conflictValidationNeeded = true
		}
		if err := s.mapper.ToDTO(&model, response); err != nil {
			return err
		}
		return nil
	})
	if conflictValidationNeeded && !utils.IsSubSetInterface(request, response) {
		err = repositories.ErrConflictValidationFailed
	}
	return err
}

func (s *serviceImpl[C, R, U, E]) GetByID(ctx context.Context, id string, response *R) error {
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	var resposeModel E
	if err := s.rpo.GetByID(ctx, key, &resposeModel); err != nil {
		return err
	}
	s.mapper.ToDTO(&resposeModel, response)
	return nil
}

func (s *serviceImpl[C, R, U, E]) List(ctx context.Context, params contracts.ListParams, responses *[]R) (int64, error) {
	mappedScopes, err := s.mapper.ApplyQueryScopes(params.QueryParams.Scopes)
	if err != nil {
		return 0, err
	}
	repoParams := toRepoListParams(params, mappedScopes)
	modelReponses := make([]E, 0, cap(*responses))
	total, err := s.rpo.List(ctx, repoParams, &modelReponses)
	if err != nil {
		return 0, err
	}
	*responses = (*responses)[:len(modelReponses)] // ensure responses slice has the same length as models
	for i := range modelReponses {
		s.mapper.ToDTO(&modelReponses[i], &(*responses)[i])
	}
	return total, nil
}

func (s *serviceImpl[C, R, U, E]) Update(ctx context.Context, id string, dto U, response *R) error {
	var model E
	if err := s.mapper.ToPatchModel(dto, &model); err != nil {
		return err
	}
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	return s.rpo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.rpo.Update(txCtx, key, &model); err != nil {
			return err
		}
		if err := s.mapper.ToDTO(&model, response); err != nil {
			return err
		}
		return nil
	})
}

func (s *serviceImpl[C, R, U, E]) Delete(ctx context.Context, id string) error {
	key, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	return s.rpo.Delete(ctx, key)
}
