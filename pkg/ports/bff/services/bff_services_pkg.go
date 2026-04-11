package services

import "github.com/brunojet/go-infra-backend/pkg/ports/bff/services/contracts"

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
