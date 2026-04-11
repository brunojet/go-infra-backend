package services

import (
	"github.com/brunojet/go-infra-backend/internal/ports/bff/services"
	bffrpocts "github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/bff/services/contracts"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// ---- Contracts ----
// CE/RE/UE are the upstream DTO types for create, read, and update respectively.
// They are unconstrained here; the BffRepository adapter constrains them to
// bffrpo.BffEntity at the implementation layer.

type (
	BffServiceMapper[C, R, U, CE, RE, UE any]       = contracts.BffServiceMapper[C, R, U, CE, RE, UE]
	BffNestedServiceMapper[C, R, U, CE, RE, UE any] = contracts.BffNestedServiceMapper[C, R, U, CE, RE, UE]
	BffService[C, R, U any]                         = contracts.BffService[C, R, U]
	BffNestedService[C, R, U any]                   = contracts.BffNestedService[C, R, U]
)

// ---- Constructors ----

// NewBffServiceImpl creates a Service[C,R,U] backed by the given BffRepository
// and BffServiceMapper. Handlers are identical to database-backed services —
// only the constructor differs.
func NewBffServiceImpl[C, R, U any, CE, RE, UE bffrpocts.BffEntity](
	rpo bffrpocts.BffRepository[CE, RE, UE],
	mapper contracts.BffServiceMapper[C, R, U, CE, RE, UE],
) svccts.Service[C, R, U] {
	return services.NewBffServiceImpl(rpo, mapper)
}

// NewBffNestedServiceImpl creates a NestedService[C,R,U] backed by the given
// BffNestedRepository and BffNestedServiceMapper.
func NewBffNestedServiceImpl[C, R, U any, CE, RE, UE bffrpocts.BffEntity](
	rpo bffrpocts.BffNestedRepository[CE, RE, UE],
	mapper contracts.BffNestedServiceMapper[C, R, U, CE, RE, UE],
) svccts.NestedService[C, R, U] {
	return services.NewBffNestedServiceImpl(rpo, mapper)
}
