package services

import (
	"github.com/brunojet/go-infra-backend/internal/ports/services"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// ---- Contracts ----
type AnyInt = contracts.AnyInt

type ServiceMapper[C, R, U any, E repositories.Entity] = contracts.ServiceMapper[C, R, U, E]

type Service[C, R, U any] = contracts.Service[C, R, U]

type NestedServiceMapper[C, R, U any, E repositories.Entity] = contracts.NestedServiceMapper[C, R, U, E]

type NestedService[C, R, U any] = contracts.NestedService[C, R, U]

// NewServiceImpl delegates to the internal implementation.
func NewServiceImpl[C, R, U any, E repositories.Entity](r repositories.Repository[E], m ServiceMapper[C, R, U, E]) Service[C, R, U] {
	return services.NewServiceImpl(r, m)
}

// NewNestedServiceImpl delegates to the internal implementation.
func NewNestedServiceImpl[C, R, U any, E repositories.Entity](r repositories.Repository[E], m NestedServiceMapper[C, R, U, E]) NestedService[C, R, U] {
	return services.NewNestedServiceImpl(r, m)
}

// ---- Utils (delegating to internal) ----
func ParseScopeIntFromString[T AnyInt](value string, minValue T) (T, error) {
	return services.ParseScopeIntFromString(value, minValue)
}

func ParseScopeInt[T AnyInt](scopes map[string]any, key string, minValue T) (T, error) {
	return services.ParseScopeInt(scopes, key, minValue)
}
