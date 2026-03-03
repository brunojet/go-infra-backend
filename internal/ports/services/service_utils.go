package services

import (
	"strings"

	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

const (
	defaultListPage    = 1
	defaultListSize    = 10
	defaultListOrderBy = "id"
	listOrderAsc       = "asc"
	listOrderDesc      = "desc"
)

func normalizeListParams(params svcContracts.ListParams) svcContracts.ListParams {
	if params.Page <= 0 {
		params.Page = defaultListPage
	}
	if params.Size <= 0 {
		params.Size = defaultListSize
	}
	if strings.TrimSpace(params.OrderBy) == "" {
		params.OrderBy = defaultListOrderBy
	}
	order := strings.ToLower(strings.TrimSpace(params.Order))
	if order != listOrderDesc {
		order = listOrderAsc
	}
	params.Order = order

	return params
}

func toRepoListParams(params svcContracts.ListParams, scopes map[string]any) repoContracts.ListParams {
	params = normalizeListParams(params)

	effectiveScopes := params.QueryParams.Scopes
	if scopes != nil {
		effectiveScopes = scopes
	}

	return repoContracts.ListParams{
		QueryParams: repoContracts.QueryParams{Scopes: effectiveScopes},
		Page:        params.Page,
		Size:        params.Size,
		OrderBy:     params.OrderBy,
		Order:       params.Order,
	}
}
