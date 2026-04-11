package services

import (
	"context"
	"strings"

	bffrpocts "github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
	bffsvccts "github.com/brunojet/go-infra-backend/pkg/ports/bff/services/contracts"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// ---------------------------------------------------------------------------
// bffServiceImpl
// ---------------------------------------------------------------------------

// bffServiceImpl is the concrete implementation of Service[C,R,U] backed by
// a BffRepository instead of a GORM database. It mirrors the logic of
// serviceImpl in internal/ports/services/service_impl.go but delegates
// domain↔upstream translation to BffServiceMapper.
type bffServiceImpl[C, R, U any, CE, RE, UE bffrpocts.BffEntity] struct {
	rpo    bffrpocts.BffRepository[CE, RE, UE]
	mapper bffsvccts.BffServiceMapper[C, R, U, CE, RE, UE]
}

// NewBffServiceImpl creates a Service[C,R,U] backed by the given BffRepository
// and BffServiceMapper.
//
// C/R/U are the domain DTOs; CE/RE/UE are the upstream (BFF) DTOs.
func NewBffServiceImpl[C, R, U any, CE, RE, UE bffrpocts.BffEntity](
	rpo bffrpocts.BffRepository[CE, RE, UE],
	mapper bffsvccts.BffServiceMapper[C, R, U, CE, RE, UE],
) svccts.Service[C, R, U] {
	return &bffServiceImpl[C, R, U, CE, RE, UE]{rpo: rpo, mapper: mapper}
}

func (s *bffServiceImpl[C, R, U, CE, RE, UE]) Create(ctx context.Context, dto C, response *R) error {
	var upstream CE
	if err := s.mapper.ToUpstreamPost(dto, &upstream); err != nil {
		return err
	}
	var upstreamResponse RE
	if err := s.rpo.Create(ctx, upstream, &upstreamResponse); err != nil {
		return err
	}
	return s.mapper.ToDomainDTO(&upstreamResponse, response)
}

func (s *bffServiceImpl[C, R, U, CE, RE, UE]) GetByID(ctx context.Context, id string, response *R) error {
	upstreamID, err := s.mapper.GetUpstreamID(id)
	if err != nil {
		return err
	}
	var upstreamResponse RE
	if err := s.rpo.GetByID(ctx, upstreamID, &upstreamResponse); err != nil {
		return err
	}
	return s.mapper.ToDomainDTO(&upstreamResponse, response)
}

func (s *bffServiceImpl[C, R, U, CE, RE, UE]) List(ctx context.Context, params svccts.ListParams, responses *[]R) (int64, error) {
	mappedScopes, err := s.mapper.ApplyQueryScopes(params.QueryParams.Scopes)
	if err != nil {
		return 0, err
	}
	repoParams := toBffListParams(params, mappedScopes)
	var upstreamResults []RE
	if err := s.rpo.List(ctx, repoParams, &upstreamResults); err != nil {
		return 0, err
	}
	total := s.mapper.ExtractUpstreamTotal(upstreamResults)
	*responses = make([]R, len(upstreamResults))
	for i := range upstreamResults {
		if err := s.mapper.ToDomainDTO(&upstreamResults[i], &(*responses)[i]); err != nil {
			return 0, err
		}
	}
	return total, nil
}

func (s *bffServiceImpl[C, R, U, CE, RE, UE]) Update(ctx context.Context, id string, dto U, response *R) error {
	upstreamID, err := s.mapper.GetUpstreamID(id)
	if err != nil {
		return err
	}
	var upstream UE
	if err := s.mapper.ToUpstreamPatch(dto, &upstream); err != nil {
		return err
	}
	var upstreamResponse RE
	if err := s.rpo.Update(ctx, upstreamID, upstream, &upstreamResponse); err != nil {
		return err
	}
	return s.mapper.ToDomainDTO(&upstreamResponse, response)
}

func (s *bffServiceImpl[C, R, U, CE, RE, UE]) Delete(ctx context.Context, id string) error {
	upstreamID, err := s.mapper.GetUpstreamID(id)
	if err != nil {
		return err
	}
	return s.rpo.Delete(ctx, upstreamID)
}

// ---------------------------------------------------------------------------
// bffNestedServiceImpl
// ---------------------------------------------------------------------------

// bffNestedServiceImpl extends bffServiceImpl with parent-scoped Create and
// List operations, mirroring nestedServiceImpl for the BFF/HTTP path.
type bffNestedServiceImpl[C, R, U any, CE, RE, UE bffrpocts.BffEntity] struct {
	svccts.Service[C, R, U]
	rpo    bffrpocts.BffNestedRepository[CE, RE, UE]
	mapper bffsvccts.BffNestedServiceMapper[C, R, U, CE, RE, UE]
}

// NewBffNestedServiceImpl creates a NestedService[C,R,U] backed by the given
// BffNestedRepository and BffNestedServiceMapper.
func NewBffNestedServiceImpl[C, R, U any, CE, RE, UE bffrpocts.BffEntity](
	rpo bffrpocts.BffNestedRepository[CE, RE, UE],
	mapper bffsvccts.BffNestedServiceMapper[C, R, U, CE, RE, UE],
) svccts.NestedService[C, R, U] {
	return &bffNestedServiceImpl[C, R, U, CE, RE, UE]{
		Service: NewBffServiceImpl(rpo, mapper),
		rpo:     rpo,
		mapper:  mapper,
	}
}

func (s *bffNestedServiceImpl[C, R, U, CE, RE, UE]) CreateNested(ctx context.Context, parentID string, dto C, response *R) error {
	upstreamParentID, err := s.mapper.GetUpstreamParentID(parentID)
	if err != nil {
		return err
	}
	var upstream CE
	if err := s.mapper.ToUpstreamPost(dto, &upstream); err != nil {
		return err
	}
	if err := s.mapper.ApplyUpstreamParentScopes(parentID, &upstream); err != nil {
		return err
	}
	var upstreamResponse RE
	if err := s.rpo.CreateNested(ctx, upstreamParentID, upstream, &upstreamResponse); err != nil {
		return err
	}
	return s.mapper.ToDomainDTO(&upstreamResponse, response)
}

func (s *bffNestedServiceImpl[C, R, U, CE, RE, UE]) ListNested(ctx context.Context, parentID string, params svccts.ListParams, responses *[]R) (int64, error) {
	upstreamParentID, err := s.mapper.GetUpstreamParentID(parentID)
	if err != nil {
		return 0, err
	}
	mergedScopes, err := s.mapper.ApplyParentQueryScopes(parentID, params.QueryParams.Scopes)
	if err != nil {
		return 0, err
	}
	repoParams := toBffListParams(params, mergedScopes)
	var upstreamResults []RE
	if err := s.rpo.ListNested(ctx, upstreamParentID, repoParams, &upstreamResults); err != nil {
		return 0, err
	}
	total := s.mapper.ExtractUpstreamTotal(upstreamResults)
	*responses = make([]R, len(upstreamResults))
	for i := range upstreamResults {
		if err := s.mapper.ToDomainDTO(&upstreamResults[i], &(*responses)[i]); err != nil {
			return 0, err
		}
	}
	return total, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// toBffListParams converts service-layer ListParams (domain) to the
// BffListParams consumed by the BffRepository, normalising defaults in the
// same way as toRepoListParams in internal/ports/services/service_utils.go.
func toBffListParams(params svccts.ListParams, scopes map[string]any) bffrpocts.BffListParams {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 10
	}
	if strings.TrimSpace(params.OrderBy) == "" {
		params.OrderBy = "created_at"
	}
	order := strings.ToLower(strings.TrimSpace(params.Order))
	if order != "desc" {
		order = "asc"
	}

	effectiveScopes := params.QueryParams.Scopes
	if scopes != nil {
		effectiveScopes = scopes
	}

	return bffrpocts.BffListParams{
		QueryParams: bffrpocts.QueryParams{Scopes: effectiveScopes},
		Page:        params.Page,
		Size:        params.Size,
		OrderBy:     params.OrderBy,
		Order:       order,
	}
}
