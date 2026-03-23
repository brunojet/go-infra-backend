package services

import (
	internalservices "github.com/brunojet/go-infra-backend/internal/ports/services"
	repo "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// ---- Contracts ----

type AnyInt = contracts.AnyInt

type ServiceMapper[D any, E repo.Entity] = contracts.ServiceMapper[D, E]

type Service[D any, E repo.Entity] = contracts.Service[D, E]

type NestedServiceMapper[D any, E repo.Entity] = contracts.NestedServiceMapper[D, E]

type NestedService[D any, E repo.Entity] = contracts.NestedService[D, E]

// ---- Constructor (delegating to internal) ----

// NewServiceImpl delegates to the internal implementation.
func NewServiceImpl[D any, E repo.Entity](r repo.Repository[E], m ServiceMapper[D, E]) Service[D, E] {
	return internalservices.NewServiceImpl(r, m)
}

// NewNestedServiceImpl delegates to the internal implementation.
func NewNestedServiceImpl[D any, E repo.Entity](r repo.Repository[E], m NestedServiceMapper[D, E]) NestedService[D, E] {
	return internalservices.NewNestedServiceImpl(r, m)
}

// ---- Utils (delegating to internal) ----
func ParseScopeIntFromString[T AnyInt](value string, minValue T) (T, error) {
	return internalservices.ParseScopeIntFromString(value, minValue)
}

func ParseScopeInt[T AnyInt](scopes map[string]any, key string, minValue T) (T, error) {
	return internalservices.ParseScopeInt(scopes, key, minValue)
}
