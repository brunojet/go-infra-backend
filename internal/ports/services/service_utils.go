package services

import (
	"strings"

	"github.com/brunojet/go-infra-backend/internal/utils"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type AnyInt interface {
	~int64 | ~int32 | ~int16
}

func normalizeListParams(params svcContracts.ListParams) svcContracts.ListParams {
	if params.Page <= 0 {
		params.Page = DefaultListPage
	}
	if params.Size <= 0 {
		params.Size = DefaultListSize
	}
	if strings.TrimSpace(params.OrderBy) == "" {
		params.OrderBy = DefaultListOrderBy
	}
	order := strings.ToLower(strings.TrimSpace(params.Order))
	if order != ListOrderDesc {
		order = ListOrderAsc
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

// ParseScopeInt extrai e valida um valor inteiro de um map de scopes.
// Retorna o valor se presente, do tipo correto e >= minValue, caso contrário retorna um erro genérico.
func ParseScopeInt[T AnyInt](scopes map[string]any, key string, minValue T) (T, error) {
	rawValue, ok := scopes[key]
	if !ok {
		var zero T
		return zero, ErrScopeKeyNotFound
	}

	value, ok := rawValue.(T)
	if !ok || value < minValue {
		var zero T
		return zero, ErrScopeValueInvalid
	}

	return value, nil
}

// StringToInt converte string para inteiro genérico, validando se > 0
func StringToInt[T AnyInt](s string) (T, error) {
	id, err := utils.StringToInt64(s)
	if err != nil || id <= 0 {
		var zero T
		return zero, ErrScopeValueInvalid
	}
	return T(id), nil
}

// ParseScopeIntFromString converte e valida string para inteiro de escopo, com valor mínimo
func ParseScopeIntFromString[T AnyInt](value string, minValue T) (T, error) {
	id, err := StringToInt[T](value)
	if err != nil || id < minValue {
		var zero T
		return zero, ErrScopeValueInvalid
	}
	return id, nil
}
