package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/ports/repositories"
	"github.com/brunojet/go-infra-backend/internal/utils"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type nestedServiceImpl[C, R, U any, E rpocts.Entity] struct {
	contracts.Service[C, R, U, E]
	rpo    rpocts.Repository[E]
	mapper contracts.NestedServiceMapper[C, R, U, E]
}

func NewNestedServiceImpl[C, R, U any, E rpocts.Entity](
	r rpocts.Repository[E],
	m contracts.NestedServiceMapper[C, R, U, E],
) contracts.NestedService[C, R, U, E] {
	return &nestedServiceImpl[C, R, U, E]{Service: NewServiceImpl(r, m), rpo: r, mapper: m}
}

func (s *nestedServiceImpl[C, R, U, E]) CreateNested(ctx context.Context, parentID string, request C, response *R) error {
	var model E
	if err := s.mapper.ToPostModel(request, &model); err != nil {
		return err
	}
	if err := s.mapper.ApplyParentScopes(parentID, &model); err != nil {
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

func (s *nestedServiceImpl[C, R, U, E]) ListNested(ctx context.Context, parentID string, listParams contracts.ListParams, responses *[]R) (int64, error) {
	mergedScopes, err := s.mapper.ApplyParentQueryScopes(parentID, listParams.QueryParams.Scopes)
	if err != nil {
		return 0, err
	}
	repoParams := toRepoListParams(listParams, mergedScopes)
	responseModels := make([]E, cap(*responses))
	total, err := s.rpo.List(ctx, repoParams, &responseModels)
	if err != nil {
		return 0, err
	}
	*responses = (*responses)[:len(responseModels)]
	for i := range responseModels {
		s.mapper.ToDTO(&responseModels[i], &(*responses)[i])
	}
	return total, nil
}
