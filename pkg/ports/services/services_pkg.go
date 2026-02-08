package services

import (
	internalservices "github.com/brunojet/go-infra-backend/internal/ports/services"
	repo "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// ---- Contracts ----

type ServiceMapper[D any, E repo.Entity] = contracts.ServiceMapper[D, E]

type Service[D any, E repo.Entity] = contracts.Service[D, E]

// ---- Constructor (delegating to internal) ----

// NewServiceImpl delegates to the internal implementation.
func NewServiceImpl[D any, E repo.Entity](r repo.Repository[E], m ServiceMapper[D, E]) Service[D, E] {
	return internalservices.NewServiceImpl[D, E](r, m)
}
