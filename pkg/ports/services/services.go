package services

import (
	internalservices "github.com/brunojet/go-infra-backend/internal/ports/services"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// NewServiceImpl delegates to the internal implementation.
func NewServiceImpl[D any, E contracts.Entity](r contracts.Repository[E], m svccontracts.ServiceMapper[D, E]) svccontracts.Service[D, E] {
	return internalservices.NewServiceImpl[D, E](r, m)
}
